package test_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/gruntwork-io/terragrunt-engine-go/proto"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type TestEngineServer struct {
	proto.UnimplementedEngineServer
}

func (m *TestEngineServer) Init(req *proto.InitRequest, stream proto.Engine_InitServer) error {
	log.Printf("Init TestEngineServer")
	return nil
}

func (m *TestEngineServer) Run(req *proto.RunRequest, stream proto.Engine_RunServer) error {
	log.Printf("Run TestEngineServer")
	return nil
}

func (m *TestEngineServer) Shutdown(req *proto.ShutdownRequest, stream proto.Engine_ShutdownServer) error {
	log.Printf("Shutdown TestEngineServer")
	return nil
}

// OneofTestEngineServer is a test server that sends responses using the oneof structure
type OneofTestEngineServer struct {
	proto.UnimplementedEngineServer
}

func (m *OneofTestEngineServer) Init(req *proto.InitRequest, stream proto.Engine_InitServer) error {
	// Send stdout message
	if err := stream.Send(&proto.InitResponse{
		Response: &proto.InitResponse_Stdout{
			Stdout: &proto.StdoutMessage{
				Content: "Initialization started\n",
			},
		},
	}); err != nil {
		return err
	}

	// Send log message
	if err := stream.Send(&proto.InitResponse{
		Response: &proto.InitResponse_Log{
			Log: &proto.LogMessage{
				Level:   proto.LogLevel_LOG_LEVEL_INFO,
				Content: "Initialization in progress",
			},
		},
	}); err != nil {
		return err
	}

	// Send exit result
	if err := stream.Send(&proto.InitResponse{
		Response: &proto.InitResponse_ExitResult{
			ExitResult: &proto.ExitResultMessage{
				Code: 0,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (m *OneofTestEngineServer) Run(req *proto.RunRequest, stream proto.Engine_RunServer) error {
	// Send stdout message
	if err := stream.Send(&proto.RunResponse{
		Response: &proto.RunResponse_Stdout{
			Stdout: &proto.StdoutMessage{
				Content: "Command output\n",
			},
		},
	}); err != nil {
		return err
	}

	// Send stderr message
	if err := stream.Send(&proto.RunResponse{
		Response: &proto.RunResponse_Stderr{
			Stderr: &proto.StderrMessage{
				Content: "Warning message\n",
			},
		},
	}); err != nil {
		return err
	}

	// Send exit result
	if err := stream.Send(&proto.RunResponse{
		Response: &proto.RunResponse_ExitResult{
			ExitResult: &proto.ExitResultMessage{
				Code: 0,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (m *OneofTestEngineServer) Shutdown(req *proto.ShutdownRequest, stream proto.Engine_ShutdownServer) error {
	// Send log message with different log levels
	logLevels := []proto.LogLevel{
		proto.LogLevel_LOG_LEVEL_DEBUG,
		proto.LogLevel_LOG_LEVEL_INFO,
		proto.LogLevel_LOG_LEVEL_WARN,
		proto.LogLevel_LOG_LEVEL_ERROR,
	}

	for _, level := range logLevels {
		if err := stream.Send(&proto.ShutdownResponse{
			Response: &proto.ShutdownResponse_Log{
				Log: &proto.LogMessage{
					Level:   level,
					Content: "Shutdown message with level " + level.String(),
				},
			},
		}); err != nil {
			return err
		}
	}

	// Send exit result
	if err := stream.Send(&proto.ShutdownResponse{
		Response: &proto.ShutdownResponse_ExitResult{
			ExitResult: &proto.ExitResultMessage{
				Code: 0,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func TestGRPCServer(t *testing.T) {
	t.Parallel()

	mockServer := &TestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	s := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	err := grpcEngine.GRPCServer(broker, s)
	require.NoError(t, err, "Expected GRPCServer to not return an error")

	// Check if the service is registered correctly
	serviceInfo := s.GetServiceInfo()
	_, ok := serviceInfo["proto.Engine"]
	assert.True(t, ok, "Expected engine.Engine service to be registered")
}

func TestGRPCClient(t *testing.T) {
	t.Parallel()

	mockServer := &TestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	server := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	var lc net.ListenConfig

	lis, err := lc.Listen(context.Background(), "tcp", ":0")
	require.NoError(t, err, "Expected no error starting listener")

	go func() {
		err := grpcEngine.GRPCServer(broker, server)
		assert.NoError(t, err, "Expected GRPCServer to not return an error")
		err = server.Serve(lis)
		assert.NoError(t, err)
	}()

	defer server.Stop()

	// nolint:staticcheck
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err, "Expected no error dialing GRPC server")

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client, err := grpcEngine.GRPCClient(context.Background(), broker, conn)
	require.NoError(t, err, "Expected GRPCClient to not return an error")
	assert.NotNil(t, client, "Expected client to be non-nil")

	engineClient, ok := client.(proto.EngineClient)
	assert.True(t, ok, "Expected client to be of type engine.EngineClient")

	// Test calling a method on the client
	stream, err := engineClient.Init(context.Background(), &proto.InitRequest{})
	require.NoError(t, err, "Expected no error calling Init")
	assert.NotNil(t, stream, "Expected Init stream to be non-nil")
}

func TestInitResponseOneof(t *testing.T) {
	t.Parallel()

	mockServer := &OneofTestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	server := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	var lc net.ListenConfig

	lis, err := lc.Listen(context.Background(), "tcp", ":0")
	require.NoError(t, err, "Expected no error starting listener")

	go func() {
		err := grpcEngine.GRPCServer(broker, server)
		assert.NoError(t, err, "Expected GRPCServer to not return an error")
		err = server.Serve(lis)
		assert.NoError(t, err)
	}()

	defer server.Stop()

	// nolint:staticcheck
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err, "Expected no error dialing GRPC server")

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client, err := grpcEngine.GRPCClient(context.Background(), broker, conn)
	require.NoError(t, err, "Expected GRPCClient to not return an error")

	engineClient, ok := client.(proto.EngineClient)
	assert.True(t, ok, "Expected client to be of type engine.EngineClient")

	// Test Init stream
	stream, err := engineClient.Init(context.Background(), &proto.InitRequest{})
	require.NoError(t, err, "Expected no error calling Init")

	// Receive and verify stdout message
	resp, err := stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")
	assert.NotNil(t, resp, "Expected response to be non-nil")

	stdout := resp.GetStdout()
	assert.NotNil(t, stdout, "Expected stdout message")
	assert.Equal(t, "Initialization started\n", stdout.GetContent(), "Expected correct stdout content")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetExitResult(), "Expected exit result to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")

	// Receive and verify log message
	resp, err = stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	logMsg := resp.GetLog()
	assert.NotNil(t, logMsg, "Expected log message")
	assert.Equal(t, proto.LogLevel_LOG_LEVEL_INFO, logMsg.GetLevel(), "Expected correct log level")
	assert.Equal(t, "Initialization in progress", logMsg.GetContent(), "Expected correct log content")
	assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetExitResult(), "Expected exit result to be nil")

	// Receive and verify exit result
	resp, err = stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	exitResult := resp.GetExitResult()
	assert.NotNil(t, exitResult, "Expected exit result message")
	assert.Equal(t, int32(0), exitResult.GetCode(), "Expected correct exit code")
	assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")
}

func TestRunResponseOneof(t *testing.T) {
	t.Parallel()

	mockServer := &OneofTestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	server := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	var lc net.ListenConfig

	lis, err := lc.Listen(context.Background(), "tcp", ":0")
	require.NoError(t, err, "Expected no error starting listener")

	go func() {
		err := grpcEngine.GRPCServer(broker, server)
		assert.NoError(t, err, "Expected GRPCServer to not return an error")
		err = server.Serve(lis)
		assert.NoError(t, err)
	}()

	defer server.Stop()

	// nolint:staticcheck
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err, "Expected no error dialing GRPC server")

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client, err := grpcEngine.GRPCClient(context.Background(), broker, conn)
	require.NoError(t, err, "Expected GRPCClient to not return an error")

	engineClient, ok := client.(proto.EngineClient)
	assert.True(t, ok, "Expected client to be of type engine.EngineClient")

	// Test Run stream
	stream, err := engineClient.Run(context.Background(), &proto.RunRequest{})
	require.NoError(t, err, "Expected no error calling Run")

	// Receive and verify stdout message
	resp, err := stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	stdout := resp.GetStdout()
	assert.NotNil(t, stdout, "Expected stdout message")
	assert.Equal(t, "Command output\n", stdout.GetContent(), "Expected correct stdout content")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetExitResult(), "Expected exit result to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")

	// Receive and verify stderr message
	resp, err = stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	stderr := resp.GetStderr()
	assert.NotNil(t, stderr, "Expected stderr message")
	assert.Equal(t, "Warning message\n", stderr.GetContent(), "Expected correct stderr content")
	assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
	assert.Nil(t, resp.GetExitResult(), "Expected exit result to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")

	// Receive and verify exit result
	resp, err = stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	exitResult := resp.GetExitResult()
	assert.NotNil(t, exitResult, "Expected exit result message")
	assert.Equal(t, int32(0), exitResult.GetCode(), "Expected correct exit code")
	assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")
}

func TestShutdownResponseOneof(t *testing.T) {
	t.Parallel()

	mockServer := &OneofTestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	server := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	var lc net.ListenConfig

	lis, err := lc.Listen(context.Background(), "tcp", ":0")
	require.NoError(t, err, "Expected no error starting listener")

	go func() {
		err := grpcEngine.GRPCServer(broker, server)
		assert.NoError(t, err, "Expected GRPCServer to not return an error")
		err = server.Serve(lis)
		assert.NoError(t, err)
	}()

	defer server.Stop()

	// nolint:staticcheck
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err, "Expected no error dialing GRPC server")

	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client, err := grpcEngine.GRPCClient(context.Background(), broker, conn)
	require.NoError(t, err, "Expected GRPCClient to not return an error")

	engineClient, ok := client.(proto.EngineClient)
	assert.True(t, ok, "Expected client to be of type engine.EngineClient")

	// Test Shutdown stream
	stream, err := engineClient.Shutdown(context.Background(), &proto.ShutdownRequest{})
	require.NoError(t, err, "Expected no error calling Shutdown")

	// Receive and verify log messages with different levels
	expectedLevels := []proto.LogLevel{
		proto.LogLevel_LOG_LEVEL_DEBUG,
		proto.LogLevel_LOG_LEVEL_INFO,
		proto.LogLevel_LOG_LEVEL_WARN,
		proto.LogLevel_LOG_LEVEL_ERROR,
	}

	for i, expectedLevel := range expectedLevels {
		resp, err := stream.Recv()
		require.NoError(t, err, "Expected no error receiving response %d", i)

		logMsg := resp.GetLog()
		assert.NotNil(t, logMsg, "Expected log message %d", i)
		assert.Equal(t, expectedLevel, logMsg.GetLevel(), "Expected correct log level for message %d", i)
		assert.Contains(t, logMsg.GetContent(), "Shutdown message with level", "Expected log content to contain expected text")
		assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
		assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
		assert.Nil(t, resp.GetExitResult(), "Expected exit result to be nil")
	}

	// Receive and verify exit result
	resp, err := stream.Recv()
	require.NoError(t, err, "Expected no error receiving response")

	exitResult := resp.GetExitResult()
	assert.NotNil(t, exitResult, "Expected exit result message")
	assert.Equal(t, int32(0), exitResult.GetCode(), "Expected correct exit code")
	assert.Nil(t, resp.GetStdout(), "Expected stdout to be nil")
	assert.Nil(t, resp.GetStderr(), "Expected stderr to be nil")
	assert.Nil(t, resp.GetLog(), "Expected log to be nil")
}
