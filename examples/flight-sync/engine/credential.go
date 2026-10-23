package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gruntwork-io/terragrunt-engine-go/examples/flight-sync/flight"
)

const (
	envAuthURL  = "FLIGHT_SYNC_AUTH_URL"
	envStateDir = "FLIGHT_SYNC_STATE_DIR"

	defaultAuthURL     = "http://127.0.0.1:8080/exchange"
	stateDirName       = "terragrunt-engine-flight-sync"
	credentialName     = "credential"
	credentialTTL      = time.Hour
	lockTimeout        = time.Minute
	exchangeTimeout    = 30 * time.Second
	maxCredentialBytes = 64 << 10
)

// credentialSource says where obtainCredential got the credential.
type credentialSource int

const (
	exchangedForFlight credentialSource = iota
	reusedFromFlight

	// exchangedAlone means Terragrunt sent no flight, so this engine
	// exchanged a credential for itself.
	exchangedAlone
)

// obtainCredential returns a credential for this engine, exchanging one with
// the auth server only when no other engine in the flight already has.
func obtainCredential(ctx context.Context, env map[string]string, flightID string, flightToken []byte) ([]byte, credentialSource, error) {
	authURL := env[envAuthURL]
	if authURL == "" {
		authURL = defaultAuthURL
	}

	if flightID == "" && len(flightToken) == 0 {
		cred, err := exchange(ctx, authURL)

		return cred, exchangedAlone, err
	}

	stateDir, err := resolveStateDir(env)
	if err != nil {
		return nil, 0, err
	}

	f, err := flight.Open(stateDir, flightID, flightToken)
	if err != nil {
		return nil, 0, err
	}

	cred, source, err := sharedCredential(ctx, f, authURL)

	return cred, source, errors.Join(err, f.Close())
}

// sharedCredential returns the flight's stored credential, exchanging and
// storing one first when none is stored or it expired.
func sharedCredential(ctx context.Context, f *flight.Flight, authURL string) (cred []byte, source credentialSource, err error) {
	lockCtx, cancel := context.WithTimeout(ctx, lockTimeout)
	defer cancel()

	lock, err := f.Lock(lockCtx, credentialName)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		err = errors.Join(err, lock.Unlock())
	}()

	now := time.Now()

	cred, found, err := f.ReadSecret(credentialName, now)
	if err != nil {
		return nil, 0, err
	}

	if found {
		return cred, reusedFromFlight, nil
	}

	cred, err = exchange(ctx, authURL)
	if err != nil {
		return nil, 0, err
	}

	if err := f.WriteSecret(credentialName, cred, now.Add(credentialTTL)); err != nil {
		return nil, 0, err
	}

	return cred, exchangedForFlight, nil
}

// exchange asks the auth server at authURL for a new credential.
//
// Errors leave out authURL, which may carry credentials of its own.
func exchange(ctx context.Context, authURL string) (cred []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, exchangeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, http.NoBody)
	if err != nil {
		return nil, errors.New("exchanging credential: invalid auth URL")
	}

	client := &http.Client{} //nolint:forbidigo // terragrunt venv rule; this example uses the real network

	resp, err := client.Do(req)
	if err != nil {
		if urlErr, ok := errors.AsType[*url.Error](err); ok {
			err = urlErr.Err
		}

		return nil, fmt.Errorf("exchanging credential: %w", err)
	}

	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchanging credential: auth server returned %s", resp.Status)
	}

	cred, err = io.ReadAll(io.LimitReader(resp.Body, maxCredentialBytes+1))
	if err != nil {
		return nil, fmt.Errorf("exchanging credential: %w", err)
	}

	if len(cred) > maxCredentialBytes {
		return nil, fmt.Errorf("exchanging credential: response is over %d bytes", maxCredentialBytes)
	}

	if len(cred) == 0 {
		return nil, errors.New("exchanging credential: response is empty")
	}

	return cred, nil
}

func resolveStateDir(env map[string]string) (string, error) {
	if dir := env[envStateDir]; dir != "" {
		return dir, nil
	}

	cacheDir, err := os.UserCacheDir() //nolint:forbidigo // terragrunt venv rule; this example uses the real user directories
	if err != nil {
		return "", fmt.Errorf("finding a state directory: %w", err)
	}

	return filepath.Join(cacheDir, stateDirName), nil
}
