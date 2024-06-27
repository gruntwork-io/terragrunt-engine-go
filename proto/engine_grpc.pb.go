// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: proto/engine.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Engine_Init_FullMethodName     = "/proto.Engine/Init"
	Engine_Run_FullMethodName      = "/proto.Engine/Run"
	Engine_Shutdown_FullMethodName = "/proto.Engine/Shutdown"
)

// EngineClient is the client API for Engine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EngineClient interface {
	// Initializes the engine with the provided request parameters.
	// Returns a stream of InitResponse messages.
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (Engine_InitClient, error)
	// Runs a command with the provided request parameters.
	// Returns a stream of RunResponse messages.
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (Engine_RunClient, error)
	// Shuts down the engine with the provided request parameters.
	// Returns a stream of ShutdownResponse messages.
	Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (Engine_ShutdownClient, error)
}

type engineClient struct {
	cc grpc.ClientConnInterface
}

func NewEngineClient(cc grpc.ClientConnInterface) EngineClient {
	return &engineClient{cc}
}

func (c *engineClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (Engine_InitClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Engine_ServiceDesc.Streams[0], Engine_Init_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &engineInitClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Engine_InitClient interface {
	Recv() (*InitResponse, error)
	grpc.ClientStream
}

type engineInitClient struct {
	grpc.ClientStream
}

func (x *engineInitClient) Recv() (*InitResponse, error) {
	m := new(InitResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *engineClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (Engine_RunClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Engine_ServiceDesc.Streams[1], Engine_Run_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &engineRunClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Engine_RunClient interface {
	Recv() (*RunResponse, error)
	grpc.ClientStream
}

type engineRunClient struct {
	grpc.ClientStream
}

func (x *engineRunClient) Recv() (*RunResponse, error) {
	m := new(RunResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *engineClient) Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (Engine_ShutdownClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Engine_ServiceDesc.Streams[2], Engine_Shutdown_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &engineShutdownClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Engine_ShutdownClient interface {
	Recv() (*ShutdownResponse, error)
	grpc.ClientStream
}

type engineShutdownClient struct {
	grpc.ClientStream
}

func (x *engineShutdownClient) Recv() (*ShutdownResponse, error) {
	m := new(ShutdownResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EngineServer is the server API for Engine service.
// All implementations must embed UnimplementedEngineServer
// for forward compatibility
type EngineServer interface {
	// Initializes the engine with the provided request parameters.
	// Returns a stream of InitResponse messages.
	Init(*InitRequest, Engine_InitServer) error
	// Runs a command with the provided request parameters.
	// Returns a stream of RunResponse messages.
	Run(*RunRequest, Engine_RunServer) error
	// Shuts down the engine with the provided request parameters.
	// Returns a stream of ShutdownResponse messages.
	Shutdown(*ShutdownRequest, Engine_ShutdownServer) error
	mustEmbedUnimplementedEngineServer()
}

// UnimplementedEngineServer must be embedded to have forward compatible implementations.
type UnimplementedEngineServer struct {
}

func (UnimplementedEngineServer) Init(*InitRequest, Engine_InitServer) error {
	return status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedEngineServer) Run(*RunRequest, Engine_RunServer) error {
	return status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (UnimplementedEngineServer) Shutdown(*ShutdownRequest, Engine_ShutdownServer) error {
	return status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedEngineServer) mustEmbedUnimplementedEngineServer() {}

// UnsafeEngineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EngineServer will
// result in compilation errors.
type UnsafeEngineServer interface {
	mustEmbedUnimplementedEngineServer()
}

func RegisterEngineServer(s grpc.ServiceRegistrar, srv EngineServer) {
	s.RegisterService(&Engine_ServiceDesc, srv)
}

func _Engine_Init_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InitRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EngineServer).Init(m, &engineInitServer{ServerStream: stream})
}

type Engine_InitServer interface {
	Send(*InitResponse) error
	grpc.ServerStream
}

type engineInitServer struct {
	grpc.ServerStream
}

func (x *engineInitServer) Send(m *InitResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Engine_Run_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RunRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EngineServer).Run(m, &engineRunServer{ServerStream: stream})
}

type Engine_RunServer interface {
	Send(*RunResponse) error
	grpc.ServerStream
}

type engineRunServer struct {
	grpc.ServerStream
}

func (x *engineRunServer) Send(m *RunResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Engine_Shutdown_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ShutdownRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EngineServer).Shutdown(m, &engineShutdownServer{ServerStream: stream})
}

type Engine_ShutdownServer interface {
	Send(*ShutdownResponse) error
	grpc.ServerStream
}

type engineShutdownServer struct {
	grpc.ServerStream
}

func (x *engineShutdownServer) Send(m *ShutdownResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Engine_ServiceDesc is the grpc.ServiceDesc for Engine service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Engine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Engine",
	HandlerType: (*EngineServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Init",
			Handler:       _Engine_Init_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Run",
			Handler:       _Engine_Run_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Shutdown",
			Handler:       _Engine_Shutdown_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/engine.proto",
}