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

	publicKey, privateKey, err := signing.NewKeyPair()
	require.NoError(t, err)

	message := structpb.NewStringValue("test")

	var signedMessage []byte
	t.Run("It should sign the proto encoded message", func(t *testing.T) {
		signedMessage, err = signing.SignProto(message, privateKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, signedMessage)
	})

	t.Run("The same message should produce the same signature", func(t *testing.T) {
		expected, err := signing.SignProto(message, privateKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, expected)

		for i := 0; i < 100; i++ {
			actual, err := signing.SignProto(message, privateKey)
			assert.NoError(t, err)
			assert.True(t, bytes.Equal(expected, actual))
		}
	})

	t.Run("The signed message should be verifiable using the public key", func(t *testing.T) {
		assert.True(t, signing.Verify(signedMessage, publicKey))
		assert.False(t, signing.Verify(signedMessage, bytes.Repeat([]byte("a"), 32)))
	})
}
