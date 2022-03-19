// Package signing provides all functionality within arrebato regarding message signing. This includes both gRPC, raft
// and data store interactions.
package signing

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/sign"
	"google.golang.org/protobuf/proto"
)

// NewKeyPair generates a new public/private key pair for message signing.
func NewKeyPair() (public []byte, private []byte, err error) {
	pu, pb, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return pu[:], pb[:], nil
}

// SignProto encodes the proto.Message and signs it using the private key, returning a signature. The
// proto-encoding is not deterministic so there is no guarantee that the same message produces the exact same signature.
// However, it is only important that the signature is verifiable against the public key, so the consistency of the
// proto-encoding should not matter. Ed25519 is used to sign messages.
func SignProto(m proto.Message, privateKey []byte) ([]byte, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	var arr [64]byte
	copy(arr[:], privateKey[:64])

	return sign.Sign(nil, data, &arr), nil
}

// Verify a signature against the public key.
func Verify(signature []byte, publicKey []byte) bool {
	var arr [32]byte
	copy(arr[:], publicKey[:32])

	_, ok := sign.Open(nil, signature, &arr)
	return ok
}
