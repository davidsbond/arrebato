// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: arrebato/node/service/v1/service.proto

package nodesvc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeServiceClient is the client API for NodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeServiceClient interface {
	// Describe the Node.
	Describe(ctx context.Context, in *DescribeRequest, opts ...grpc.CallOption) (*DescribeResponse, error)
	// Watch for changes to the node, including leadership status and new peers. Updates are written to the returned
	// stream.
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (NodeService_WatchClient, error)
	// Backup the node state, streaming the byte contents back to the client.
	Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (NodeService_BackupClient, error)
}

type nodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeServiceClient(cc grpc.ClientConnInterface) NodeServiceClient {
	return &nodeServiceClient{cc}
}

func (c *nodeServiceClient) Describe(ctx context.Context, in *DescribeRequest, opts ...grpc.CallOption) (*DescribeResponse, error) {
	out := new(DescribeResponse)
	err := c.cc.Invoke(ctx, "/arrebato.node.service.v1.NodeService/Describe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (NodeService_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeService_ServiceDesc.Streams[0], "/arrebato.node.service.v1.NodeService/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeServiceWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NodeService_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type nodeServiceWatchClient struct {
	grpc.ClientStream
}

func (x *nodeServiceWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeServiceClient) Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (NodeService_BackupClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeService_ServiceDesc.Streams[1], "/arrebato.node.service.v1.NodeService/Backup", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeServiceBackupClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NodeService_BackupClient interface {
	Recv() (*BackupResponse, error)
	grpc.ClientStream
}

type nodeServiceBackupClient struct {
	grpc.ClientStream
}

func (x *nodeServiceBackupClient) Recv() (*BackupResponse, error) {
	m := new(BackupResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeServiceServer is the server API for NodeService service.
// All implementations should embed UnimplementedNodeServiceServer
// for forward compatibility
type NodeServiceServer interface {
	// Describe the Node.
	Describe(context.Context, *DescribeRequest) (*DescribeResponse, error)
	// Watch for changes to the node, including leadership status and new peers. Updates are written to the returned
	// stream.
	Watch(*WatchRequest, NodeService_WatchServer) error
	// Backup the node state, streaming the byte contents back to the client.
	Backup(*BackupRequest, NodeService_BackupServer) error
}

// UnimplementedNodeServiceServer should be embedded to have forward compatible implementations.
type UnimplementedNodeServiceServer struct {
}

func (UnimplementedNodeServiceServer) Describe(context.Context, *DescribeRequest) (*DescribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Describe not implemented")
}
func (UnimplementedNodeServiceServer) Watch(*WatchRequest, NodeService_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedNodeServiceServer) Backup(*BackupRequest, NodeService_BackupServer) error {
	return status.Errorf(codes.Unimplemented, "method Backup not implemented")
}

// UnsafeNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServiceServer will
// result in compilation errors.
type UnsafeNodeServiceServer interface {
	mustEmbedUnimplementedNodeServiceServer()
}

func RegisterNodeServiceServer(s grpc.ServiceRegistrar, srv NodeServiceServer) {
	s.RegisterService(&NodeService_ServiceDesc, srv)
}

func _NodeService_Describe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Describe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arrebato.node.service.v1.NodeService/Describe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Describe(ctx, req.(*DescribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeServiceServer).Watch(m, &nodeServiceWatchServer{stream})
}

type NodeService_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type nodeServiceWatchServer struct {
	grpc.ServerStream
}

func (x *nodeServiceWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NodeService_Backup_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BackupRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeServiceServer).Backup(m, &nodeServiceBackupServer{stream})
}

type NodeService_BackupServer interface {
	Send(*BackupResponse) error
	grpc.ServerStream
}

type nodeServiceBackupServer struct {
	grpc.ServerStream
}

func (x *nodeServiceBackupServer) Send(m *BackupResponse) error {
	return x.ServerStream.SendMsg(m)
}

// NodeService_ServiceDesc is the grpc.ServiceDesc for NodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arrebato.node.service.v1.NodeService",
	HandlerType: (*NodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Describe",
			Handler:    _NodeService_Describe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _NodeService_Watch_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Backup",
			Handler:       _NodeService_Backup_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arrebato/node/service/v1/service.proto",
}
