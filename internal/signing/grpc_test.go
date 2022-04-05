package signing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/command"
	signingcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/command/v1"
	signingsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/service/v1"
	signingpb "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/v1"
	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_CreateKeyPair(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *signingsvc.CreateKeyPairRequest
		ClientInfo   clientinfo.ClientInfo
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name:    "It should generate a new signing key pair",
			Request: &signingsvc.CreateKeyPairRequest{},
			ClientInfo: clientinfo.ClientInfo{
				ID: "test-client",
			},
			Expected: command.New(&signingcmd.CreatePublicKey{
				ClientId: "test-client",
			}),
		},
		{
			Name:    "It should return codes.FailedPrecondition if the node is not the leader",
			Request: &signingsvc.CreateKeyPairRequest{},
			ClientInfo: clientinfo.ClientInfo{
				ID: "test-client",
			},
			Error:        command.ErrNotLeader,
			ExpectedCode: codes.FailedPrecondition,
		},
		{
			Name:    "It should return codes.AlreadyExists if the client already has a key pair",
			Request: &signingsvc.CreateKeyPairRequest{},
			ClientInfo: clientinfo.ClientInfo{
				ID: "test-client",
			},
			Error:        signing.ErrPublicKeyExists,
			ExpectedCode: codes.AlreadyExists,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			executor := &MockExecutor{err: tc.Error}
			ctx := clientinfo.ToContext(testutil.Context(t), tc.ClientInfo)

			resp, err := signing.NewGRPC(executor, nil).CreateKeyPair(ctx, tc.Request)
			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*signingcmd.CreatePublicKey)
			actual := executor.command.Payload().(*signingcmd.CreatePublicKey)

			assert.EqualValues(t, expected.GetClientId(), actual.GetClientId())
			assert.EqualValues(t, resp.GetKeyPair().GetPublicKey(), actual.GetPublicKey())
			assert.NotEmpty(t, actual.GetPublicKey())
		})
	}
}

func TestGRPC_GetPublicKey(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *signingsvc.GetPublicKeyRequest
		Expected     *signingsvc.GetPublicKeyResponse
		ExpectedCode codes.Code
		Error        error
	}{
		{
			Name: "It should return the public key",
			Request: &signingsvc.GetPublicKeyRequest{
				ClientId: "test-client",
			},
			Expected: &signingsvc.GetPublicKeyResponse{
				PublicKey: &signingpb.PublicKey{
					ClientId:  "test-client",
					PublicKey: []byte("test"),
				},
			},
		},
		{
			Name: "It should return codes.NotFound if the public key does not exit",
			Request: &signingsvc.GetPublicKeyRequest{
				ClientId: "test-client",
			},
			ExpectedCode: codes.NotFound,
			Error:        signing.ErrNoPublicKey,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getter := &MockPublicKeyGetter{err: tc.Error, key: []byte("test")}
			ctx := testutil.Context(t)

			resp, err := signing.NewGRPC(nil, getter).GetPublicKey(ctx, tc.Request)
			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.True(t, proto.Equal(tc.Expected, resp))
		})
	}
}

func TestGRPC_ListPublicKeys(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Expected     *signingsvc.ListPublicKeysResponse
		ExpectedCode codes.Code
		Error        error
	}{
		{
			Name: "It should return all public keys",
			Expected: &signingsvc.ListPublicKeysResponse{
				PublicKeys: []*signingpb.PublicKey{
					{
						ClientId:  "test-client",
						PublicKey: []byte("test"),
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getter := &MockPublicKeyGetter{err: tc.Error, key: []byte("test")}
			ctx := testutil.Context(t)

			resp, err := signing.NewGRPC(nil, getter).ListPublicKeys(ctx, nil)
			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.True(t, proto.Equal(tc.Expected, resp))
		})
	}
}
