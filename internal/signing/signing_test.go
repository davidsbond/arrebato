package signing_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/davidsbond/arrebato/internal/signing"
)

func TestSignProto(t *testing.T) {
	t.Parallel()

	keyPair, err := signing.NewKeyPair()
	require.NoError(t, err)

	message := structpb.NewStringValue("test")

	var signedMessage []byte
	t.Run("It should sign the proto encoded message", func(t *testing.T) {
		signedMessage, err = signing.SignProto(message, keyPair.Private)
		assert.NoError(t, err)
		assert.NotEmpty(t, signedMessage)
	})

	t.Run("The signed message should be verifiable using the public key", func(t *testing.T) {
		verified, err := signing.Verify(signedMessage, keyPair.Public)
		assert.True(t, verified)
		assert.NoError(t, err)

		verified, err = signing.Verify(signedMessage, bytes.Repeat([]byte("a"), 32))
		assert.False(t, verified)
		assert.NoError(t, err)
	})
}
