package tlsinfo_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/tlsinfo"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Handler      func(t *testing.T) grpc.UnaryHandler
		Context      func(ctx context.Context) context.Context
		ExpectedCode codes.Code
	}{
		{
			Name: "It should return codes.Unauthenticated if there is no peer info",
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should return codes.Unauthenticated if the peer info does not have TLS auth",
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should return codes.Unauthenticated if there is no client certificate data",
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: &credentials.TLSInfo{
						State: tls.ConnectionState{
							VerifiedChains: [][]*x509.Certificate{},
						},
					},
				})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should add the TLSInfo to the context on success",
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					info := tlsinfo.FromContext(ctx)
					assert.EqualValues(t, "test", info.CommonName)
					return nil, nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						State: tls.ConnectionState{
							VerifiedChains: [][]*x509.Certificate{
								{
									{
										Subject: pkix.Name{
											CommonName: "test",
										},
									},
								},
							},
						},
					},
				})
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			interceptor := tlsinfo.UnaryServerInterceptor()
			_, err := interceptor(tc.Context(ctx), nil, nil, tc.Handler(t))
			assert.EqualValues(t, tc.ExpectedCode, status.Code(err))
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Handler      func(t *testing.T) grpc.StreamHandler
		Context      func(ctx context.Context) context.Context
		ExpectedCode codes.Code
	}{
		{
			Name: "It should return codes.Unauthenticated if there is no peer info",
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should return codes.Unauthenticated if the peer info does not have TLS auth",
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should return codes.Unauthenticated if there is no client certificate data",
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: &credentials.TLSInfo{
						State: tls.ConnectionState{
							VerifiedChains: [][]*x509.Certificate{},
						},
					},
				})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "It should add the TLSInfo to the context on success",
			Handler: func(t *testing.T) grpc.StreamHandler {
				return func(srv interface{}, stream grpc.ServerStream) error {
					info := tlsinfo.FromContext(stream.Context())
					assert.EqualValues(t, "test", info.CommonName)
					return nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						State: tls.ConnectionState{
							VerifiedChains: [][]*x509.Certificate{
								{
									{
										Subject: pkix.Name{
											CommonName: "test",
										},
									},
								},
							},
						},
					},
				})
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			interceptor := tlsinfo.StreamServerInterceptor()
			err := interceptor(nil, MockServerStream{ctx: tc.Context(ctx)}, nil, tc.Handler(t))
			assert.EqualValues(t, tc.ExpectedCode, status.Code(err))
		})
	}
}
