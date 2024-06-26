package test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var listener *bufconn.Listener

func init() {
	listener = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}

type mockCommandExecutor struct {
	engine.UnimplementedCommandExecutorServer
}

func (m *mockCommandExecutor) Init(req *engine.InitRequest, stream engine.CommandExecutor_InitServer) error {
	response := &engine.InitResponse{
		Stdout:     "Mock initialization successful",
		Stderr:     "",
		ResultCode: 0,
	}
	return stream.Send(response)
}

func (m *mockCommandExecutor) Run(req *engine.RunRequest, stream engine.CommandExecutor_RunServer) error {
	response := &engine.RunResponse{
		Stdout:     "Mock command output",
		Stderr:     "",
		ResultCode: 0,
	}
	return stream.Send(response)
}

func (m *mockCommandExecutor) Shutdown(req *engine.ShutdownRequest, stream engine.CommandExecutor_ShutdownServer) error {
	response := &engine.ShutdownResponse{
		Stdout:     "Mock shutdown successful",
		Stderr:     "",
		ResultCode: 0,
	}
	return stream.Send(response)
}

func startTestServer() {
	s := grpc.NewServer()
	engine.RegisterCommandExecutorServer(s, &mockCommandExecutor{})
	go func() {
		if err := s.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()
}

func TestTerragruntGRPCEngine(t *testing.T) {
	startTestServer()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func() {
		assert.NoError(t, conn.Close())
	}()

	client := engine.NewCommandExecutorClient(conn)

	t.Run("Test Init", func(t *testing.T) {
		stream, err := client.Init(context.Background(), &engine.InitRequest{})
		if err != nil {
			t.Fatalf("Failed to call Init: %v", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Failed to receive Init response: %v", err)
			}
			if resp.ResultCode != 0 {
				t.Errorf("Expected result code 0, got %d", resp.ResultCode)
			}
		}
	})

	t.Run("Test Run", func(t *testing.T) {
		stream, err := client.Run(context.Background(), &engine.RunRequest{
			Command: "mock-command",
			Args:    []string{"arg1", "arg2"},
		})
		if err != nil {
			t.Fatalf("Failed to call Run: %v", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Failed to receive Run response: %v", err)
			}
			if resp.Stdout != "Mock command output" {
				t.Errorf("Expected stdout 'Mock command output', got %s", resp.Stdout)
			}
		}
	})

	t.Run("Test Shutdown", func(t *testing.T) {
		stream, err := client.Shutdown(context.Background(), &engine.ShutdownRequest{})
		if err != nil {
			t.Fatalf("Failed to call Shutdown: %v", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Failed to receive Shutdown response: %v", err)
			}
			if resp.ResultCode != 0 {
				t.Errorf("Expected result code 0, got %d", resp.ResultCode)
			}
		}
	})
}
