package arrebato

import (
	"context"
	"io"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
)

type (
	backupReader struct {
		buffer []byte
		stream nodesvc.NodeService_BackupClient
	}
)

// Read the contents of the backup, writing it to the provided byte slice.
func (b *backupReader) Read(p []byte) (int, error) {
	// If we already have data in the buffer from a previous read, provide that
	// to the caller.
	if len(b.buffer) > 0 {
		n := copy(p, b.buffer)
		b.buffer = b.buffer[n:]
		return n, nil
	}

	// Otherwise, request new data from the server. gRPC streams will return an io.EOF when
	// completed, we can just propagate the error upwards.
	msg, err := b.stream.Recv()
	if err != nil {
		return 0, err
	}

	n := copy(p, msg.GetData())
	if n < len(msg.GetData()) {
		// If the provided byte slice can't take the entirety of the data we got from the
		// server, store the remainder in memory within the buffer.
		b.buffer = msg.GetData()[n:]
	}

	return n, nil
}

// Close the connection from the server.
func (b *backupReader) Close() error {
	return b.stream.CloseSend()
}

// Backup performs a backup of the server state. It returns an io.ReadCloser implementation that, when read, contains
// the contents of the backup.
func (c *Client) Backup(ctx context.Context) (io.ReadCloser, error) {
	stream, err := nodesvc.NewNodeServiceClient(c.cluster.leader()).Backup(ctx, &nodesvc.BackupRequest{})
	if err != nil {
		return nil, err
	}

	return &backupReader{stream: stream}, nil
}
