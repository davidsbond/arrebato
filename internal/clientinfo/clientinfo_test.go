package clientinfo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(testutil.Context(t), metadata.New(map[string]string{
		"client_id": "test",
	}))

	t.Run("It should pass the client info from the metadata to the context", func(t *testing.T) {
		interceptor := clientinfo.UnaryServerInterceptor()

		_, err := interceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			info := clientinfo.FromContext(ctx)
			assert.EqualValues(t, "test", info.ID)
			return nil, nil
		})

		assert.NoError(t, err)
	})
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(testutil.Context(t), metadata.New(map[string]string{
		"client_id": "test",
	}))

	t.Run("It should pass the client info from the metadata to the context", func(t *testing.T) {
		interceptor := clientinfo.StreamServerInterceptor()

		err := interceptor(ctx, &MockServerStream{ctx: ctx}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			info := clientinfo.FromContext(stream.Context())
			assert.EqualValues(t, "test", info.ID)
			return nil
		})

		assert.NoError(t, err)
	})
}
