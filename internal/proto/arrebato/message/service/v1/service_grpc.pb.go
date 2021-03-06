// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: arrebato/message/service/v1/service.proto

package messagesvc

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	// Produce a new message on a topic. This RPC should return a NOT_FOUND error code if the topic does not exist.
	Produce(ctx context.Context, in *ProduceRequest, opts ...grpc.CallOption) (*ProduceResponse, error)
	// Consume messages from a topic. This RPC should return a NOT_FOUND error code if the topic does not exist.
	Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (MessageService_ConsumeClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) Produce(ctx context.Context, in *ProduceRequest, opts ...grpc.CallOption) (*ProduceResponse, error) {
	out := new(ProduceResponse)
	err := c.cc.Invoke(ctx, "/arrebato.message.service.v1.MessageService/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (MessageService_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/arrebato.message.service.v1.MessageService/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_ConsumeClient interface {
	Recv() (*ConsumeResponse, error)
	grpc.ClientStream
}

type messageServiceConsumeClient struct {
	grpc.ClientStream
}

func (x *messageServiceConsumeClient) Recv() (*ConsumeResponse, error) {
	m := new(ConsumeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations should embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	// Produce a new message on a topic. This RPC should return a NOT_FOUND error code if the topic does not exist.
	Produce(context.Context, *ProduceRequest) (*ProduceResponse, error)
	// Consume messages from a topic. This RPC should return a NOT_FOUND error code if the topic does not exist.
	Consume(*ConsumeRequest, MessageService_ConsumeServer) error
}

// UnimplementedMessageServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) Produce(context.Context, *ProduceRequest) (*ProduceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedMessageServiceServer) Consume(*ConsumeRequest, MessageService_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProduceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arrebato.message.service.v1.MessageService/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Produce(ctx, req.(*ProduceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).Consume(m, &messageServiceConsumeServer{stream})
}

type MessageService_ConsumeServer interface {
	Send(*ConsumeResponse) error
	grpc.ServerStream
}

type messageServiceConsumeServer struct {
	grpc.ServerStream
}

func (x *messageServiceConsumeServer) Send(m *ConsumeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arrebato.message.service.v1.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _MessageService_Produce_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Consume",
			Handler:       _MessageService_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arrebato/message/service/v1/service.proto",
}
