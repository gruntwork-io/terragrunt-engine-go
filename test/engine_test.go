package test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/gruntwork-io/terragrunt-engine-go/proto"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
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

func TestGRPCServer(t *testing.T) {
	mockServer := &TestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	s := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	err := grpcEngine.GRPCServer(broker, s)
	assert.Nil(t, err, "Expected GRPCServer to not return an error")

	// Check if the service is registered correctly
	serviceInfo := s.GetServiceInfo()
	_, ok := serviceInfo["proto.Engine"]
	assert.True(t, ok, "Expected engine.Engine service to be registered")
}

func TestGRPCClient(t *testing.T) {
	mockServer := &TestEngineServer{}
	grpcEngine := &engine.TerragruntGRPCEngine{Impl: mockServer}
	server := grpc.NewServer()
	broker := &plugin.GRPCBroker{}

	lis, err := net.Listen("tcp", ":0")
	assert.Nil(t, err, "Expected no error starting listener")

	go func() {
		err := grpcEngine.GRPCServer(broker, server)
		assert.Nil(t, err, "Expected GRPCServer to not return an error")
		err = server.Serve(lis)
		assert.NoError(t, err)
	}()
	defer server.Stop()

	// nolint:staticcheck
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	assert.Nil(t, err, "Expected no error dialing GRPC server")
	defer func() {
		err := conn.Close()
		assert.NoError(t, err)
	}()

	client, err := grpcEngine.GRPCClient(context.Background(), broker, conn)
	assert.Nil(t, err, "Expected GRPCClient to not return an error")
	assert.NotNil(t, client, "Expected client to be non-nil")

	engineClient, ok := client.(proto.EngineClient)
	assert.True(t, ok, "Expected client to be of type engine.EngineClient")

	// Test calling a method on the client
	stream, err := engineClient.Init(context.Background(), &proto.InitRequest{})
	assert.Nil(t, err, "Expected no error calling Init")
	assert.NotNil(t, stream, "Expected Init stream to be non-nil")
}
