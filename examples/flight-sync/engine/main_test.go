package main_test

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"uuid"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

const engines = 4

func TestEnginesInOneFlightExchangeOnce(t *testing.T) {
	t.Parallel()

	bin := buildEngine(t)
	authURL, exchanges := startAuthServer(t)

	token := make([]byte, 32)
	_, err := rand.Read(token)
	require.NoError(t, err)

	req := &tgengine.InitRequest{
		FlightId:    uuid.NewV7().String(),
		FlightToken: token,
		EnvVars: map[string]string{
			"FLIGHT_SYNC_AUTH_URL":  authURL,
			"FLIGHT_SYNC_STATE_DIR": t.TempDir(),
		},
	}

	require.NoError(t, initEngines(t, bin, req))
	assert.Equal(t, int64(1), exchanges.Load())
}

func TestEnginesWithoutFlightExchangeAlone(t *testing.T) {
	t.Parallel()

	bin := buildEngine(t)
	authURL, exchanges := startAuthServer(t)

	req := &tgengine.InitRequest{
		EnvVars: map[string]string{
			"FLIGHT_SYNC_AUTH_URL":  authURL,
			"FLIGHT_SYNC_STATE_DIR": t.TempDir(),
		},
	}

	require.NoError(t, initEngines(t, bin, req))
	assert.Equal(t, int64(engines), exchanges.Load())
}

// initEngines starts the engine processes the way Terragrunt does, and sends
// each of them req at once.
func initEngines(t *testing.T, bin string, req *tgengine.InitRequest) error {
	t.Helper()

	g, ctx := errgroup.WithContext(t.Context())

	for range engines {
		client := startEngine(t, bin)

		g.Go(func() error {
			return initEngine(ctx, client, req)
		})
	}

	return g.Wait()
}

func initEngine(ctx context.Context, client tgengine.EngineClient, req *tgengine.InitRequest) error {
	stream, err := client.Init(ctx, req)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return errors.New("init ended without an exit result")
		}

		if err != nil {
			return err
		}

		if result := resp.GetExitResult(); result != nil {
			if result.GetCode() != 0 {
				return errors.New("init exited non-zero")
			}

			return nil
		}
	}
}

func startEngine(t *testing.T, bin string) tgengine.EngineClient {
	t.Helper()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "engine",
			MagicCookieValue: "terragrunt",
		},
		Plugins: map[string]plugin.Plugin{
			"plugin": &engine.TerragruntGRPCEngine{},
		},
		Cmd:              exec.CommandContext(t.Context(), bin),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
	})
	t.Cleanup(client.Kill)

	rpcClient, err := client.Client()
	require.NoError(t, err)

	raw, err := rpcClient.Dispense("plugin")
	require.NoError(t, err)

	engineClient, ok := raw.(tgengine.EngineClient)
	require.True(t, ok)

	return engineClient
}

func buildEngine(t *testing.T) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "terragrunt-engine-flight-sync")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	out, err := exec.CommandContext(t.Context(), "go", "build", "-o", bin, ".").CombinedOutput()
	require.NoError(t, err, string(out))

	return bin
}

// startAuthServer starts a stand-in auth server and returns its exchange URL
// and the count of exchanges it served.
func startAuthServer(t *testing.T) (string, *atomic.Int64) {
	t.Helper()

	exchanges := new(atomic.Int64)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		exchanges.Add(1)

		_, err := w.Write([]byte(rand.Text()))
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	return server.URL + "/exchange", exchanges
}
