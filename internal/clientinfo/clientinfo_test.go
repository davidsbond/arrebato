package clientinfo_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Extractor    clientinfo.Extractor
		Handler      func(t *testing.T) grpc.UnaryHandler
		Context      func(ctx context.Context) context.Context
		ExpectedCode codes.Code
	}{
		{
			Name:      "It should return codes.InvalidArgument if there is no inbound metadata",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should return codes.InvalidArgument if the inbound metadata is empty",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return metadata.NewIncomingContext(ctx, metadata.MD{})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should add the client information to the context when the header is present",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					info := clientinfo.FromContext(ctx)
					assert.EqualValues(t, "test", info.ID)
					return nil, nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				return metadata.NewIncomingContext(ctx, metadata.MD{
					"x-client-id": []string{"test"},
				})
			},
		},
		{
			Name:      "It should return codes.Unauthenticated if there is no peer info",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name:      "It should return codes.Unauthenticated if the peer info does not have TLS auth",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name:      "It should return codes.InvalidArgument if there's no SPIFFE ID",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{},
				})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should return codes.InvalidArgument if the SPIFFE ID does not match the metadata header",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						SPIFFEID: &url.URL{
							Scheme: "spiffe",
							Host:   "test",
						},
					},
				})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should add the ClientInfo to the context on success from the TLS certificate",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					info := clientinfo.FromContext(ctx)
					assert.EqualValues(t, "spiffe://test", info.ID)
					return nil, nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				ctx = peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						SPIFFEID: &url.URL{
							Scheme: "spiffe",
							Host:   "test",
						},
					},
				})

				return metadata.NewIncomingContext(ctx, metadata.MD{
					"x-client-id": []string{"spiffe://test"},
				})
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			interceptor := clientinfo.UnaryServerInterceptor(tc.Extractor)
			_, err := interceptor(tc.Context(ctx), nil, nil, tc.Handler(t))
			assert.EqualValues(t, tc.ExpectedCode, status.Code(err))
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Extractor    clientinfo.Extractor
		Handler      func(t *testing.T) grpc.StreamHandler
		Context      func(ctx context.Context) context.Context
		ExpectedCode codes.Code
	}{
		{
			Name:      "It should return codes.InvalidArgument if there is no inbound metadata",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should return codes.InvalidArgument if the inbound metadata is empty",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return metadata.NewIncomingContext(ctx, metadata.MD{})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should add the client information to the context when the header is present",
			Extractor: clientinfo.MetadataExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return func(srv interface{}, stream grpc.ServerStream) error {
					info := clientinfo.FromContext(stream.Context())
					assert.EqualValues(t, "test", info.ID)
					return nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				return metadata.NewIncomingContext(ctx, metadata.MD{
					"x-client-id": []string{"test"},
				})
			},
		},
		{
			Name:      "It should return codes.Unauthenticated if there is no peer info",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return ctx
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name:      "It should return codes.Unauthenticated if the peer info does not have TLS auth",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{})
			},
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name:      "It should return codes.InvalidArgument if there's no SPIFFE ID",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{},
				})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should return codes.InvalidArgument if the SPIFFE ID does not match the metadata header",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return nil
			},
			Context: func(ctx context.Context) context.Context {
				return peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						SPIFFEID: &url.URL{
							Scheme: "spiffe",
							Host:   "test",
						},
					},
				})
			},
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:      "It should add the ClientInfo to the context on success from the TLS certificate",
			Extractor: clientinfo.TLSExtractor,
			Handler: func(t *testing.T) grpc.StreamHandler {
				return func(srv interface{}, stream grpc.ServerStream) error {
					info := clientinfo.FromContext(stream.Context())
					assert.EqualValues(t, "spiffe://test", info.ID)
					return nil
				}
			},
			Context: func(ctx context.Context) context.Context {
				ctx = peer.NewContext(ctx, &peer.Peer{
					AuthInfo: credentials.TLSInfo{
						SPIFFEID: &url.URL{
							Scheme: "spiffe",
							Host:   "test",
						},
					},
				})

				return metadata.NewIncomingContext(ctx, metadata.MD{
					"x-client-id": []string{"spiffe://test"},
				})
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			interceptor := clientinfo.StreamServerInterceptor(tc.Extractor)
			err := interceptor(nil, MockServerStream{ctx: tc.Context(ctx)}, nil, tc.Handler(t))
			assert.EqualValues(t, tc.ExpectedCode, status.Code(err))
		})
	}
}
