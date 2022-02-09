// Package command provides a wrapper around the raft log that used strongly-typed commands for communicating state
// changes across cluster members.
package command

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type (
	// The Command type represents a proto-encoded command that can be applied to a raft cluster.
	Command struct {
		payload proto.Message
	}

	// The Executor type is responsible for executing commands across the raft cluster.
	Executor struct {
		applier LogApplier
		timeout time.Duration
	}

	// The LogApplier interface describes types that can apply logs to the raft node.
	LogApplier interface {
		ApplyLog(log raft.Log, timeout time.Duration) raft.ApplyFuture
		State() raft.RaftState
	}
)

// NewExecutor returns a new instance of the Executor type that will execute Command instances via the given
// LogApplier. The timeout is applied across all Command executions.
func NewExecutor(applier LogApplier, timeout time.Duration) *Executor {
	return &Executor{
		applier: applier,
		timeout: timeout,
	}
}

// ErrNotLeader is the error given when trying to execute a Command on a node that is not the leader.
var ErrNotLeader = errors.New("not leader")

// Execute the provided command.
func (ex *Executor) Execute(ctx context.Context, cmd Command) error {
	if ex.applier.State() != raft.Leader {
		return ErrNotLeader
	}

	data, err := Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	log := raft.Log{
		Data: data,
	}

	future := ex.applier.ApplyLog(log, ex.timeout)
	if future.Error() != nil {
		return future.Error()
	}

	resp := future.Response()
	switch val := resp.(type) {
	case error:
		return val
	default:
		return nil
	}
}

// New returns a new Command wrapping a proto.Message implementation.
func New(message proto.Message) Command {
	return Command{payload: message}
}

// Payload returns the underlying proto.Message implementation describing the Command.
func (cmd Command) Payload() proto.Message {
	return cmd.payload
}

// Unmarshal converts a byte slice into a Command.
func Unmarshal(b []byte) (*Command, error) {
	input := bytes.NewBuffer(b)
	writer := bytes.NewBuffer([]byte{})

	reader, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}

	for {
		_, err = io.CopyN(writer, reader, 1024)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	if err = reader.Close(); err != nil {
		return nil, err
	}

	var a anypb.Any
	if err := proto.Unmarshal(writer.Bytes(), &a); err != nil {
		return nil, err
	}

	payload, err := a.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	return &Command{payload: payload}, nil
}

// Marshal converts a command into its byte representation.
func Marshal(cmd Command) ([]byte, error) {
	a, err := anypb.New(cmd.payload)
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(a)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer([]byte{})
	writer := gzip.NewWriter(reader)

	if _, err = io.Copy(writer, bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	return reader.Bytes(), nil
}
