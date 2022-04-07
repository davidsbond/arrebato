// Package clientinfo provides gRPC interceptors for acquiring client information from inbound gRPC requests and making
// that information available via context keys.
package clientinfo

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type (
	// The ClientInfo type describes the client as obtained via the inbound gRPC context.
	ClientInfo struct {
		ID string
	}

	// The Extractor type is a function that generates a ClientInfo based on the content of the provided context. The
	// context should be that of a gRPC request, so details can be obtained via metadata, tls certificates etc.
	Extractor func(ctx context.Context) (ClientInfo, error)

	ctxKey struct{}
)

// FromContext obtains a ClientInfo instance from the provided context.Context. Returns a zero-value ClientInfo if one does
// not exist in the context.
func FromContext(ctx context.Context) ClientInfo {
	val := ctx.Value(ctxKey{})
	if info, ok := val.(ClientInfo); ok {
		return info
	}

	return ClientInfo{}
}

// ToContext adds the provided ClientInfo to the parent context, returning a new context. It can be retrieved from the
// context using FromContext.
func ToContext(ctx context.Context, info ClientInfo) context.Context {
	return context.WithValue(ctx, ctxKey{}, info)
}

// UnaryServerInterceptor is a grpc.UnaryServerInterceptor implementation that adds a ClientInfo to inbound contexts
// via the provided Extractor implementation.
func UnaryServerInterceptor(fn Extractor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		clientInfo, err := fn(ctx)
		if err != nil {
			return nil, err
		}

		ctx = ToContext(ctx, clientInfo)
		return handler(ctx, req)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream

	ctx context.Context
}

// Context returns the desired context for the grpc.ServerStream instead of the default one.
func (ss wrappedServerStream) Context() context.Context {
	return ss.ctx
}

// StreamServerInterceptor is a grpc.StreamServerInterceptor implementation that adds a ClientInfo to inbound contexts
// via the provided Extractor implementation.
func StreamServerInterceptor(fn Extractor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientInfo, err := fn(ss.Context())
		if err != nil {
			return err
		}

		stream := wrappedServerStream{
			ServerStream: ss,
			ctx:          ToContext(ss.Context(), clientInfo),
		}

		return handler(srv, stream)
	}
}

// TLSExtractor is an Extractor implementation that attempts to generate a ClientInfo using the client's TLS certificate.
// The client identifier will be the peer's SPIFFE identifier.
func TLSExtractor(ctx context.Context) (ClientInfo, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ClientInfo{}, status.Error(codes.Unauthenticated, "no peer in incoming context")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ClientInfo{}, status.Error(codes.Unauthenticated, "peer is not using TLS, is your client configured correctly?")
	}

	if tlsInfo.SPIFFEID == nil {
		return ClientInfo{}, status.Error(codes.InvalidArgument, "peer info is missing SPIFFE ID")
	}

	return ClientInfo{
		ID: tlsInfo.SPIFFEID.String(),
	}, nil
}

// MetadataExtractor is an Extractor implementation that attempts to generate a ClientInfo using the inbound context's
// metadata fields. The client identifier will be taken from the X-Client-ID metadata field.
func MetadataExtractor(ctx context.Context) (ClientInfo, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ClientInfo{}, status.Error(codes.InvalidArgument, "incoming context contains no metadata")
	}

	values := md.Get("X-Client-ID")
	if len(values) == 0 {
		return ClientInfo{}, status.Error(codes.InvalidArgument, "X-Client-ID metadata is not set")
	}

	return ClientInfo{
		ID: strings.Join(values, ""),
	}, nil
}

func UnaryClientInterceptor(clientID string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if clientID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "X-Client-ID", clientID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(clientID string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if clientID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "X-Client-ID", clientID)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
