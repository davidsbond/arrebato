package tlsinfo_test

import (
	"context"

	"google.golang.org/grpc"
)

type (
	MockServerStream struct {
		grpc.ServerStream

		ctx context.Context
	}
)

func (m MockServerStream) Context() context.Context {
	return m.ctx
}
