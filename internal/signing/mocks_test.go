package signing_test

import (
	"context"

	"github.com/davidsbond/arrebato/internal/command"
)

type (
	MockPublicKeyCreator struct {
		created []byte
		err     error
	}

	MockExecutor struct {
		command command.Command
		err     error
	}

	MockPublicKeyGetter struct {
		err error
		key []byte
	}
)

func (mm *MockPublicKeyGetter) Get(ctx context.Context, clientID string) ([]byte, error) {
	return mm.key, mm.err
}

func (mm *MockPublicKeyCreator) Create(ctx context.Context, clientID string, publicKey []byte) error {
	mm.created = publicKey
	return mm.err
}

func (mm *MockExecutor) Execute(ctx context.Context, cmd command.Command) error {
	if mm.err != nil {
		return mm.err
	}

	mm.command = cmd
	return nil
}
