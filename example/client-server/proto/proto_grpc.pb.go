// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: proto/proto.proto

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
	ShellService_RunCommand_FullMethodName = "/proto.ShellService/RunCommand"
)

// ShellServiceClient is the client API for ShellService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShellServiceClient interface {
	RunCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
}

type shellServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShellServiceClient(cc grpc.ClientConnInterface) ShellServiceClient {
	return &shellServiceClient{cc}
}

func (c *shellServiceClient) RunCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, ShellService_RunCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShellServiceServer is the server API for ShellService service.
// All implementations must embed UnimplementedShellServiceServer
// for forward compatibility
type ShellServiceServer interface {
	RunCommand(context.Context, *CommandRequest) (*CommandResponse, error)
	mustEmbedUnimplementedShellServiceServer()
}

// UnimplementedShellServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShellServiceServer struct {
}

func (UnimplementedShellServiceServer) RunCommand(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunCommand not implemented")
}
func (UnimplementedShellServiceServer) mustEmbedUnimplementedShellServiceServer() {}

// UnsafeShellServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShellServiceServer will
// result in compilation errors.
type UnsafeShellServiceServer interface {
	mustEmbedUnimplementedShellServiceServer()
}

func RegisterShellServiceServer(s grpc.ServiceRegistrar, srv ShellServiceServer) {
	s.RegisterService(&ShellService_ServiceDesc, srv)
}

func _ShellService_RunCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShellServiceServer).RunCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShellService_RunCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShellServiceServer).RunCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShellService_ServiceDesc is the grpc.ServiceDesc for ShellService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShellService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ShellService",
	HandlerType: (*ShellServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunCommand",
			Handler:    _ShellService_RunCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proto.proto",
}
