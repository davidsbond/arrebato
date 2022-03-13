package acl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/acl"
	aclpb "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Input        *aclpb.ACL
		Expected     *aclpb.ACL
		ExpectsError bool
	}{
		{
			Name: "It should return an error for an invalid permission",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic:  "test-topic",
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_UNSPECIFIED,
						},
					},
				},
			},
			ExpectsError: true,
		},
		{
			Name: "It should return an error for a missing client name",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic: "test-topic",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_CONSUME,
						},
					},
				},
			},
			ExpectsError: true,
		},
		{
			Name: "It should return an error for a missing topic name",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_CONSUME,
						},
					},
				},
			},
			ExpectsError: true,
		},
		{
			Name: "It should return an error for a missing permissions",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic:  "test-topic",
						Client: "test-client",
					},
				},
			},
			ExpectsError: true,
		},
		{
			Name: "It should remove duplicate permissions from an entry",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic:  "test-topic",
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_PRODUCE,
							aclpb.Permission_PERMISSION_PRODUCE,
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
							aclpb.Permission_PERMISSION_PRODUCE,
						},
					},
				},
			},
		},
		{
			Name: "It should remove duplicate entries",
			Input: &aclpb.ACL{
				Entries: []*aclpb.Entry{
					{
						Topic:  "test-topic",
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_PRODUCE,
						},
					},
					{
						Topic:  "test-topic",
						Client: "test-client",
						Permissions: []aclpb.Permission{
							aclpb.Permission_PERMISSION_PRODUCE,
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
							aclpb.Permission_PERMISSION_PRODUCE,
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			actual, err := acl.Normalize(ctx, tc.Input)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.True(t, proto.Equal(tc.Expected, actual))
		})
	}
}
