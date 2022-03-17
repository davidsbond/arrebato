package acl_test

import (
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/arrebato/internal/acl"
	aclcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/command/v1"
	aclpb "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestHandler_Set(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *aclcmd.SetACL
		Expected *aclpb.ACL
		Error    error
	}{
		{
			Name: "It should set the ACL",
			Command: &aclcmd.SetACL{
				Acl: &aclpb.ACL{
					Entries: []*aclpb.Entry{
						{
							Topic:  "test-topic",
							Client: "test-client",
							Permissions: []aclpb.Permission{
								aclpb.Permission_PERMISSION_CONSUME,
							},
						},
					},
				},
			},
			Expected: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic:  "test-topic",
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_CONSUME,
						},
					},
				},
			},
		},
		{
			Name: "It should propagate errors from the setter",
			Command: &aclcmd.SetACL{
				Acl: &aclpb.ACL{
					Entries: []*aclpb.Entry{
						{
							Topic:  "test-topic",
							Client: "test-client",
							Permissions: []aclpb.Permission{
								aclpb.Permission_PERMISSION_CONSUME,
							},
						},
					},
				},
			},
			Error: io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockSetter{err: tc.Error}
			err := acl.NewHandler(mock, hclog.NewNullLogger()).Set(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Expected, mock.set)
		})
	}
}
