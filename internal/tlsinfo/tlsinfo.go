// Package tlsinfo provides gRPC interceptors for ensuring the TLS certificates used by clients are well-formed and
// making that information easily available to gRPC handlers.
package tlsinfo

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type (
	// The TLSInfo type describes the contents of a client's TLS certificate as obtained via the inbound gRPC context.
	TLSInfo struct {
		CommonName string
	}

	ctxKey struct{}
)

// FromContext obtains a TLSInfo instance from the provided context.Context. Returns a zero-value TLSInfo if one does
// not exist in the context.
func FromContext(ctx context.Context) TLSInfo {
	val := ctx.Value(ctxKey{})
	if info, ok := val.(TLSInfo); ok {
		return info
	}

	return TLSInfo{}
}

// ToContext adds the provided TLSInfo to the parent context, returning a new context. It can be retrieved from the
// context using FromContext.
func ToContext(ctx context.Context, info TLSInfo) context.Context {
	return context.WithValue(ctx, ctxKey{}, info)
}

// UnaryServerInterceptor is a grpc.UnaryServerInterceptor implementation that ensures the client is using a correctly
// formed TLS certificate for the server to handle. If the server trusts the client, the TLS certificate's common name
// must be present. Otherwise, a codes.Unauthenticated error code is returned to the client.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		clientInfo, err := fromIncomingContext(ctx)
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

// StreamServerInterceptor is a grpc.StreamServerInterceptor implementation that ensures the client is using a correctly
// formed TLS certificate for the server to handle. If the server trusts the client, the TLS certificate's common name
// must be present. Otherwise, a codes.Unauthenticated error code is returned to the client.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientInfo, err := fromIncomingContext(ss.Context())
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

func fromIncomingContext(ctx context.Context) (TLSInfo, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return TLSInfo{}, status.Error(codes.Unauthenticated, "no peer in incoming context")
	}

	raw, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return TLSInfo{}, status.Error(codes.Unauthenticated, "peer is not using TLS, is your client configured correctly?")
	}

	if len(raw.State.VerifiedChains) == 0 || len(raw.State.VerifiedChains[0]) == 0 {
		return TLSInfo{}, status.Error(codes.Unauthenticated, "could not verify peer certificate")
	}

	clientCert := raw.State.VerifiedChains[0][0]
	return TLSInfo{
		CommonName: clientCert.Subject.CommonName,
	}, nil
}
