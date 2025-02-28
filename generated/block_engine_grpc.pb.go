// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: block_engine.proto

package types

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BlockEngineValidator_StreamMempool_FullMethodName    = "/block_engine.BlockEngineValidator/StreamMempool"
	BlockEngineValidator_SubscribeBundles_FullMethodName = "/block_engine.BlockEngineValidator/SubscribeBundles"
)

// BlockEngineValidatorClient is the client API for BlockEngineValidator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// / Validators can connect to Block Engines to receive packets and bundles.
type BlockEngineValidatorClient interface {
	// / Validators can submit message to the block engine
	StreamMempool(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[MempoolPacket, StreamMempoolResponse], error)
	// / Validators can subscribe to the block engine to receive a stream of simulated and profitable bundles
	SubscribeBundles(ctx context.Context, in *SubscribeBundlesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ValidatorBundle], error)
}

type blockEngineValidatorClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockEngineValidatorClient(cc grpc.ClientConnInterface) BlockEngineValidatorClient {
	return &blockEngineValidatorClient{cc}
}

func (c *blockEngineValidatorClient) StreamMempool(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[MempoolPacket, StreamMempoolResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BlockEngineValidator_ServiceDesc.Streams[0], BlockEngineValidator_StreamMempool_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[MempoolPacket, StreamMempoolResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockEngineValidator_StreamMempoolClient = grpc.ClientStreamingClient[MempoolPacket, StreamMempoolResponse]

func (c *blockEngineValidatorClient) SubscribeBundles(ctx context.Context, in *SubscribeBundlesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ValidatorBundle], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BlockEngineValidator_ServiceDesc.Streams[1], BlockEngineValidator_SubscribeBundles_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeBundlesRequest, ValidatorBundle]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockEngineValidator_SubscribeBundlesClient = grpc.ServerStreamingClient[ValidatorBundle]

// BlockEngineValidatorServer is the server API for BlockEngineValidator service.
// All implementations must embed UnimplementedBlockEngineValidatorServer
// for forward compatibility.
//
// / Validators can connect to Block Engines to receive packets and bundles.
type BlockEngineValidatorServer interface {
	// / Validators can submit message to the block engine
	StreamMempool(grpc.ClientStreamingServer[MempoolPacket, StreamMempoolResponse]) error
	// / Validators can subscribe to the block engine to receive a stream of simulated and profitable bundles
	SubscribeBundles(*SubscribeBundlesRequest, grpc.ServerStreamingServer[ValidatorBundle]) error
	mustEmbedUnimplementedBlockEngineValidatorServer()
}

// UnimplementedBlockEngineValidatorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBlockEngineValidatorServer struct{}

func (UnimplementedBlockEngineValidatorServer) StreamMempool(grpc.ClientStreamingServer[MempoolPacket, StreamMempoolResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamMempool not implemented")
}
func (UnimplementedBlockEngineValidatorServer) SubscribeBundles(*SubscribeBundlesRequest, grpc.ServerStreamingServer[ValidatorBundle]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeBundles not implemented")
}
func (UnimplementedBlockEngineValidatorServer) mustEmbedUnimplementedBlockEngineValidatorServer() {}
func (UnimplementedBlockEngineValidatorServer) testEmbeddedByValue()                              {}

// UnsafeBlockEngineValidatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockEngineValidatorServer will
// result in compilation errors.
type UnsafeBlockEngineValidatorServer interface {
	mustEmbedUnimplementedBlockEngineValidatorServer()
}

func RegisterBlockEngineValidatorServer(s grpc.ServiceRegistrar, srv BlockEngineValidatorServer) {
	// If the following call pancis, it indicates UnimplementedBlockEngineValidatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BlockEngineValidator_ServiceDesc, srv)
}

func _BlockEngineValidator_StreamMempool_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BlockEngineValidatorServer).StreamMempool(&grpc.GenericServerStream[MempoolPacket, StreamMempoolResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockEngineValidator_StreamMempoolServer = grpc.ClientStreamingServer[MempoolPacket, StreamMempoolResponse]

func _BlockEngineValidator_SubscribeBundles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeBundlesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BlockEngineValidatorServer).SubscribeBundles(m, &grpc.GenericServerStream[SubscribeBundlesRequest, ValidatorBundle]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockEngineValidator_SubscribeBundlesServer = grpc.ServerStreamingServer[ValidatorBundle]

// BlockEngineValidator_ServiceDesc is the grpc.ServiceDesc for BlockEngineValidator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockEngineValidator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "block_engine.BlockEngineValidator",
	HandlerType: (*BlockEngineValidatorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamMempool",
			Handler:       _BlockEngineValidator_StreamMempool_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SubscribeBundles",
			Handler:       _BlockEngineValidator_SubscribeBundles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "block_engine.proto",
}
