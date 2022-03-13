package acl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/acl"
	"github.com/davidsbond/arrebato/internal/command"
	aclcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/command/v1"
	aclsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/service/v1"
	aclpb "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_Get(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *aclsvc.GetRequest
		ACL          *aclpb.ACL
		ExpectedCode codes.Code
		Error        error
	}{
		{
			Name:    "It should return the server's ACL",
			Request: &aclsvc.GetRequest{},
			ACL: &aclpb.ACL{
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
			Name:         "It should return codes.NotFound if no ACL has been set up",
			Request:      &aclsvc.GetRequest{},
			Error:        acl.ErrNoACL,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			resp, err := acl.NewGRPC(nil, &MockGetter{acl: tc.ACL, err: tc.Error}).Get(ctx, tc.Request)

			require.EqualValues(t, tc.ExpectedCode, status.Code(err))
			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.True(t, proto.Equal(tc.ACL, resp.GetAcl()))
		})
	}
}

func TestGRPC_Set(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *aclsvc.SetRequest
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name: "It should execute a command to set the ACL",
			Request: &aclsvc.SetRequest{
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
			Expected: command.New(&aclcmd.SetACL{
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
			}),
		},
		{
			Name: "It should return a failed precondition if the node is not the leader",
			Request: &aclsvc.SetRequest{
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
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name: "It should return invalid argument if the ACL is invalid",
			Request: &aclsvc.SetRequest{
				Acl: &aclpb.ACL{
					Entries: []*aclpb.Entry{
						{},
					},
				},
			},
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			executor := &MockExecutor{err: tc.Error}

			resp, err := acl.NewGRPC(executor, nil).Set(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*aclcmd.SetACL)
			actual := executor.command.Payload().(*aclcmd.SetACL)

			assert.True(t, proto.Equal(expected, actual))
		})
	}
}
