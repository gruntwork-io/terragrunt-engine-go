// Command terragrunt-engine-flight-sync is an example engine whose processes
// share one credential per Terragrunt invocation. The first engine to take the
// flight's lock exchanges the credential with an auth server and stores it
// encrypted in the flight directory. The other engines wait for the lock and
// read it from there.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"

	"github.com/hashicorp/go-plugin"
)

const (
	pluginVersion = 1

	// fingerprintBytes is how much of the credential's hash the engine logs,
	// enough to tell credentials apart without revealing them.
	fingerprintBytes = 4
)

type flightSyncEngine struct {
	tgengine.UnimplementedEngineServer
}

// Init obtains the credential this engine would use, and logs whether it
// exchanged it or reused the flight's.
func (e *flightSyncEngine) Init(req *tgengine.InitRequest, stream tgengine.Engine_InitServer) error {
	cred, source, err := obtainCredential(stream.Context(), req.GetEnvVars(), req.GetFlightId(), req.GetFlightToken())
	if err != nil {
		return err
	}

	if err := stream.Send(&tgengine.InitResponse{
		Response: &tgengine.InitResponse_Log{
			Log: &tgengine.LogMessage{
				Level:   tgengine.LogLevel_LOG_LEVEL_INFO,
				Content: describeCredential(cred, source),
			},
		},
	}); err != nil {
		return err
	}

	return stream.Send(&tgengine.InitResponse{
		Response: &tgengine.InitResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{Code: 0},
		},
	})
}

// Run runs the command locally, streaming its output back to Terragrunt.
func (e *flightSyncEngine) Run(req *tgengine.RunRequest, stream tgengine.Engine_RunServer) error {
	// exec copies stdout and stderr on separate goroutines, and a gRPC stream
	// does not allow concurrent sends.
	var sendMu sync.Mutex

	cmd := exec.CommandContext(stream.Context(), req.GetCommand(), req.GetArgs()...) //nolint:forbidigo // terragrunt venv rule; this example runs tofu directly
	cmd.Dir = req.GetWorkingDir()
	cmd.Stdout = sendWriter(func(p []byte) error {
		sendMu.Lock()
		defer sendMu.Unlock()

		return stream.Send(&tgengine.RunResponse{
			Response: &tgengine.RunResponse_Stdout{Stdout: &tgengine.StdoutMessage{Content: string(p)}},
		})
	})
	cmd.Stderr = sendWriter(func(p []byte) error {
		sendMu.Lock()
		defer sendMu.Unlock()

		return stream.Send(&tgengine.RunResponse{
			Response: &tgengine.RunResponse_Stderr{Stderr: &tgengine.StderrMessage{Content: string(p)}},
		})
	})

	for key, value := range req.GetEnvVars() {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	code, err := exitCode(cmd.Run())
	if err != nil {
		return err
	}

	return stream.Send(&tgengine.RunResponse{
		Response: &tgengine.RunResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{Code: code},
		},
	})
}

// Shutdown has nothing to clean up. The flight directory stays for the
// flight's other engines.
func (e *flightSyncEngine) Shutdown(_ *tgengine.ShutdownRequest, stream tgengine.Engine_ShutdownServer) error {
	return stream.Send(&tgengine.ShutdownResponse{
		Response: &tgengine.ShutdownResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{Code: 0},
		},
	})
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			MagicCookieKey:   "engine",
			MagicCookieValue: "terragrunt",
		},
		VersionedPlugins: map[int]plugin.PluginSet{
			pluginVersion: {
				"flight-sync-engine": &engine.TerragruntGRPCEngine{Impl: &flightSyncEngine{}},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type sendWriter func(p []byte) error

func (w sendWriter) Write(p []byte) (int, error) {
	if err := w(p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func describeCredential(cred []byte, source credentialSource) string {
	sum := sha256.Sum256(cred)
	fingerprint := hex.EncodeToString(sum[:fingerprintBytes])

	switch source {
	case exchangedForFlight:
		return fmt.Sprintf("Exchanged credential %s and stored it for the flight", fingerprint)
	case reusedFromFlight:
		return "Reused the flight's credential " + fingerprint
	case exchangedAlone:
		return fmt.Sprintf("Exchanged credential %s for this engine alone, because Terragrunt sent no flight", fingerprint)
	}

	panic(fmt.Sprintf("unknown credential source %d", source))
}

// exitCode returns the exit code of a finished command, or err when the
// command did not run to an exit.
func exitCode(err error) (int32, error) {
	if err == nil {
		return 0, nil
	}

	exitErr, ok := errors.AsType[*exec.ExitError](err)
	if !ok {
		return 0, err
	}

	return int32(exitErr.ExitCode()), nil
}
