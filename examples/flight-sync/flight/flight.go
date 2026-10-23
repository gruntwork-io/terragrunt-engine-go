// Package flight lets the engine processes of one Terragrunt invocation, a
// flight, coordinate through files. It uses the flight ID that Terragrunt sends
// in InitRequest to name a directory for lock files, and the flight token to
// encrypt the secrets the engines leave there for each other.
package flight

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	"uuid"
)

// TokenBytes is the length of the flight token Terragrunt sends.
const TokenBytes = 32

const (
	dirPerm  = 0o700
	filePerm = 0o600

	lockPollInterval = 50 * time.Millisecond

	// secretsKeyInfo separates the secrets key from any other key an engine
	// derives from the same flight token.
	secretsKeyInfo = "terragrunt-engine-flight-sync secrets"

	// keyBytes selects AES-256.
	keyBytes = 32

	expiryBytes = 8
)

var (
	// ErrInvalidID is returned by [Open] when the flight ID is not a UUID in
	// canonical form.
	ErrInvalidID = errors.New("flight ID is not a canonical UUID")

	// ErrInvalidToken is returned by [Open] when the flight token is not
	// [TokenBytes] long.
	ErrInvalidToken = errors.New("flight token has the wrong length")
)

// Flight is the directory one flight's engines share.
type Flight struct {
	root *os.Root
	aead cipher.AEAD
}

// Lock is a held lock from [Flight.Lock].
type Lock struct {
	file *os.File
}

// Open opens the directory for the flight id under baseDir, creating both as
// needed. Every engine process in the flight opens the same directory.
//
// Returns [ErrInvalidID] when id is not a canonical UUID, and
// [ErrInvalidToken] when token is not [TokenBytes] long.
func Open(baseDir, id string, token []byte) (*Flight, error) {
	parsed, err := uuid.Parse(id)
	if err != nil || parsed.String() != id {
		return nil, ErrInvalidID
	}

	if len(token) != TokenBytes {
		return nil, fmt.Errorf("%w: got %d bytes, want %d", ErrInvalidToken, len(token), TokenBytes)
	}

	aead, err := newAEAD(token)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(baseDir, id)

	//nolint:forbidigo // terragrunt venv rule; this example works on the real filesystem
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return nil, fmt.Errorf("creating flight directory: %w", err)
	}

	//nolint:forbidigo // terragrunt venv rule; this example works on the real filesystem
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("opening flight directory: %w", err)
	}

	return &Flight{root: root, aead: aead}, nil
}

// Close releases the flight directory. It leaves the files in place for the
// flight's other engines.
func (f *Flight) Close() error {
	return f.root.Close()
}

// Lock takes the exclusive lock called name, shared by every process in the
// flight. It waits until the lock is free or ctx ends. The operating system
// releases the lock if the holder exits without calling [Lock.Unlock].
func (f *Flight) Lock(ctx context.Context, name string) (*Lock, error) {
	file, err := f.root.OpenFile(name+".lock", os.O_CREATE|os.O_RDWR, filePerm)
	if err != nil {
		return nil, fmt.Errorf("opening lock %s: %w", name, err)
	}

	ticker := time.NewTicker(lockPollInterval)
	defer ticker.Stop()

	for {
		locked, err := tryLock(file)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("locking %s: %w", name, err), file.Close())
		}

		if locked {
			return &Lock{file: file}, nil
		}

		select {
		case <-ctx.Done():
			return nil, errors.Join(fmt.Errorf("waiting for lock %s: %w", name, ctx.Err()), file.Close())
		case <-ticker.C:
		}
	}
}

// Unlock releases the lock.
func (l *Lock) Unlock() error {
	return errors.Join(unlock(l.file), l.file.Close())
}

// ReadSecret returns the secret stored under name. It reports false when no
// secret is stored, when the stored secret expired at or before now, or when
// the file does not decrypt with this flight's token.
//
// Hold the lock that guards name while calling it, so no other engine
// rewrites the file mid-read.
func (f *Flight) ReadSecret(name string, now time.Time) ([]byte, bool, error) {
	data, err := f.root.ReadFile(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("reading secret %s: %w", name, err)
	}

	payload, ok := f.decrypt(name, data)
	if !ok {
		// A write torn by a crash fails authentication here, and the caller
		// replaces the file.
		return nil, false, nil
	}

	expiresAt := time.Unix(0, int64(binary.BigEndian.Uint64(payload[:expiryBytes])))
	if !now.Before(expiresAt) {
		return nil, false, nil
	}

	return payload[expiryBytes:], true, nil
}

// WriteSecret encrypts secret with the flight token and stores it under name
// until expiresAt. The expiry is encrypted with the secret, so editing the file
// cannot extend it, and the ciphertext is bound to name, so a copy under
// another name does not decrypt.
//
// Hold the lock that guards name while calling it.
func (f *Flight) WriteSecret(name string, secret []byte, expiresAt time.Time) error {
	payload := binary.BigEndian.AppendUint64(nil, uint64(expiresAt.UnixNano()))
	payload = append(payload, secret...)

	nonce := make([]byte, f.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generating nonce: %w", err)
	}

	sealed := f.aead.Seal(nonce, nonce, payload, []byte(name))

	if err := f.root.WriteFile(name, sealed, filePerm); err != nil {
		return fmt.Errorf("writing secret %s: %w", name, err)
	}

	return nil
}

// decrypt returns the payload sealed in data for name, and reports false when
// data does not authenticate.
func (f *Flight) decrypt(name string, data []byte) ([]byte, bool) {
	nonceSize := f.aead.NonceSize()
	if len(data) < nonceSize {
		return nil, false
	}

	payload, err := f.aead.Open(nil, data[:nonceSize], data[nonceSize:], []byte(name))
	if err != nil || len(payload) < expiryBytes {
		return nil, false
	}

	return payload, true
}

func newAEAD(token []byte) (cipher.AEAD, error) {
	key, err := hkdf.Key(sha256.New, token, nil, secretsKeyInfo, keyBytes)
	if err != nil {
		return nil, fmt.Errorf("deriving secrets key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	return aead, nil
}
