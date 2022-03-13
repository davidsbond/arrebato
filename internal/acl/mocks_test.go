package acl_test

import (
	"context"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
)

type (
	MockSetter struct {
		set *acl.ACL
		err error
	}

	MockGetter struct {
		acl *acl.ACL
		err error
	}

	MockExecutor struct {
		command command.Command
		err     error
	}
)

func (mm *MockGetter) Get(_ context.Context) (*acl.ACL, error) {
	return mm.acl, mm.err
}

func (mm *MockSetter) Set(_ context.Context, a *acl.ACL) error {
	mm.set = a
	return mm.err
}

func (mm *MockExecutor) Execute(_ context.Context, cmd command.Command) error {
	if mm.err != nil {
		return mm.err
	}

	mm.command = cmd
	return nil
}
