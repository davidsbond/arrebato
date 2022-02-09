// Package clientinfo provides functions and gRPC interceptors to obtain and validate metadata sent via the client
// to the server.
package clientinfo

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	// The ClientInfo type describes client metadata as sent to the server.
	ClientInfo struct {
		// The Client's identifier.
		ID string
	}

	ctxKey struct{}
)

// FromContext attempts to return a ClientInfo instance stored within the provided context. Returns a zero-value
// ClientInfo if one does not exist in the context.
func FromContext(ctx context.Context) ClientInfo {
	val := ctx.Value(ctxKey{})
	if info, ok := val.(ClientInfo); ok {
		return info
	}

	return ClientInfo{}
}

// ToContext adds the provided ClientInfo instance to the context.Context. Returning a new context.Context.
func ToContext(ctx context.Context, info ClientInfo) context.Context {
	return context.WithValue(ctx, ctxKey{}, info)
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor implementation that obtains client metadata from the
// incoming gRPC context and adds it to the context provided to the handling function.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		clientInfo := fromIncomingContext(ctx)
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

// StreamServerInterceptor returns a grpc.StreamServerInterceptor implementation that obtains client metadata from the
// incoming gRPC context and adds it to the context provided to the handling function.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientInfo := fromIncomingContext(ss.Context())
		stream := wrappedServerStream{
			ServerStream: ss,
			ctx:          ToContext(ss.Context(), clientInfo),
		}

		return handler(srv, stream)
	}
}

func fromIncomingContext(ctx context.Context) ClientInfo {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ClientInfo{}
	}

	values := md.Get("client_id")
	if len(values) == 0 {
		return ClientInfo{}
	}

	return ClientInfo{
		ID: strings.Join(values, ";"),
	}
}
