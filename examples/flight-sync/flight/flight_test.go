package flight_test

import (
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
	"time"
	"uuid"

	"github.com/gruntwork-io/terragrunt-engine-go/examples/flight-sync/flight"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenRejectsInvalidFlights(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7().String()

	testCases := []struct {
		err   error
		name  string
		id    string
		token []byte
	}{
		{name: "empty id", id: "", token: newToken(t), err: flight.ErrInvalidID},
		{name: "not a uuid", id: "not-a-uuid", token: newToken(t), err: flight.ErrInvalidID},
		{name: "path escape", id: "../../outside", token: newToken(t), err: flight.ErrInvalidID},
		{name: "non-canonical uuid", id: strings.ToUpper(id), token: newToken(t), err: flight.ErrInvalidID},
		{name: "empty token", id: id, token: nil, err: flight.ErrInvalidToken},
		{name: "short token", id: id, token: make([]byte, flight.TokenBytes-1), err: flight.ErrInvalidToken},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := flight.Open(t.TempDir(), tc.id, tc.token)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestSecretRoundTripsWithinFlight(t *testing.T) {
	t.Parallel()

	dir, id, token := t.TempDir(), uuid.NewV7().String(), newToken(t)
	now := time.Now()

	writer := openFlight(t, dir, id, token)
	require.NoError(t, writer.WriteSecret("credential", []byte("jwt"), now.Add(time.Hour)))

	reader := openFlight(t, dir, id, token)

	got, found, err := reader.ReadSecret("credential", now)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, []byte("jwt"), got)

	_, found, err = reader.ReadSecret("credential", now.Add(time.Hour))
	require.NoError(t, err)
	assert.False(t, found, "a secret is gone at its expiry")

	_, found, err = reader.ReadSecret("missing", now)
	require.NoError(t, err)
	assert.False(t, found)
}

func TestReadSecretRejectsOtherTokens(t *testing.T) {
	t.Parallel()

	dir, id := t.TempDir(), uuid.NewV7().String()
	now := time.Now()

	writer := openFlight(t, dir, id, newToken(t))
	require.NoError(t, writer.WriteSecret("credential", []byte("jwt"), now.Add(time.Hour)))

	_, found, err := openFlight(t, dir, id, newToken(t)).ReadSecret("credential", now)
	require.NoError(t, err)
	assert.False(t, found)
}

func TestReadSecretRejectsTamperedFiles(t *testing.T) {
	t.Parallel()

	dir, id, token := t.TempDir(), uuid.NewV7().String(), newToken(t)
	now := time.Now()

	f := openFlight(t, dir, id, token)
	require.NoError(t, f.WriteSecret("credential", []byte("jwt"), now.Add(time.Hour)))
	require.NoError(t, f.WriteSecret("other", []byte("other"), now.Add(time.Hour)))

	path := filepath.Join(dir, id, "credential")

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	data[len(data)-1] ^= 1
	require.NoError(t, os.WriteFile(path, data, 0o600))

	_, found, err := f.ReadSecret("credential", now)
	require.NoError(t, err)
	assert.False(t, found, "a flipped bit fails authentication")

	require.NoError(t, os.WriteFile(path, data[:5], 0o600))

	_, found, err = f.ReadSecret("credential", now)
	require.NoError(t, err)
	assert.False(t, found, "a truncated file fails authentication")

	other, err := os.ReadFile(filepath.Join(dir, id, "other"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, other, 0o600))

	_, found, err = f.ReadSecret("credential", now)
	require.NoError(t, err)
	assert.False(t, found, "a secret copied from another name fails authentication")
}

func TestLockExcludesOtherHolders(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		dir, id, token := t.TempDir(), uuid.NewV7().String(), newToken(t)
		first, second := openFlight(t, dir, id, token), openFlight(t, dir, id, token)

		held, err := first.Lock(t.Context(), "credential")
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		_, err = second.Lock(ctx, "credential")
		require.ErrorIs(t, err, context.DeadlineExceeded)

		acquired := make(chan error, 1)

		go func() {
			lock, err := second.Lock(t.Context(), "credential")
			if err == nil {
				err = lock.Unlock()
			}

			acquired <- err
		}()

		synctest.Sleep(time.Second)

		select {
		case err := <-acquired:
			require.FailNow(t, "Lock returned while another holder had the lock", "err: %v", err)
		default:
		}

		require.NoError(t, held.Unlock())

		synctest.Sleep(time.Second)
		require.NoError(t, <-acquired)
	})
}

func openFlight(t *testing.T, dir, id string, token []byte) *flight.Flight {
	t.Helper()

	f, err := flight.Open(dir, id, token)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, f.Close())
	})

	return f
}

func newToken(t *testing.T) []byte {
	t.Helper()

	token := make([]byte, flight.TokenBytes)

	_, err := rand.Read(token)
	require.NoError(t, err)

	return token
}
