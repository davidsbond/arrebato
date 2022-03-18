package e2e_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type (
	SigningSuite struct {
		suite.Suite

		client *arrebato.Client
	}
)

func (s *SigningSuite) TestKeys() {
	ctx := testutil.Context(s.T())

	var keyPair arrebato.KeyPair
	var err error

	s.T().Run("It should create a key pair for signing messages", func(t *testing.T) {
		keyPair, err = s.client.CreateSigningKeyPair(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, keyPair.PublicKey)
		require.NotEmpty(t, keyPair.PrivateKey)
	})

	s.T().Run("It should obtain the same public key from the server", func(t *testing.T) {
		publicKey, err := s.client.SigningPublicKey(ctx, "test-client")
		require.NoError(t, err)
		require.EqualValues(t, keyPair.PublicKey, publicKey)
	})
}
