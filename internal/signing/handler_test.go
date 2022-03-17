package signing_test

import (
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	signingcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/command/v1"
	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestHandler_Set(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *signingcmd.CreatePublicKey
		Expected []byte
		Error    error
	}{
		{
			Name: "It should set the ACL",
			Command: &signingcmd.CreatePublicKey{
				ClientId:  "test",
				PublicKey: []byte("test"),
			},
			Expected: []byte("test"),
		},
		{
			Name:    "It should propagate errors from the setter",
			Command: &signingcmd.CreatePublicKey{},
			Error:   io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockPublicKeyCreator{err: tc.Error}
			err := signing.NewHandler(mock, hclog.NewNullLogger()).Create(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Expected, mock.created)
		})
	}
}
