// Package signing provides all functionality within arrebato regarding message signing. This includes both gRPC, raft
// and data store interactions.
package signing

import (
	"crypto/rand"
	"encoding/base64"
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

// SignProto encodes the proto.Message and signs it using the private key, returning a base64-encoded signature. The
// proto-encoding is not deterministic so there is no guarantee that the same message produces the exact same signature.
// However, it is only important that the signature is verifiable against the public key, so the consistency of the
// proto-encoding should not matter. Ed25519 is used to sign messages.
//
// Base64 encoding is used so that the signature can be safely transported via gRPC. The raw signature may contain
// characters that don't play nicely with outbound gRPC requests.
func SignProto(m proto.Message, privateKey []byte) ([]byte, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	var arr [64]byte
	copy(arr[:], privateKey[:64])

	signature := sign.Sign(nil, data, &arr)
	out := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(out, signature)

	return out, nil
}

// Verify a base64-encoded signature against the public key.
func Verify(signature []byte, publicKey []byte) (bool, error) {
	var arr [32]byte
	copy(arr[:], publicKey[:32])

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(signature)))
	n, err := base64.StdEncoding.Decode(dst, signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	_, ok := sign.Open(nil, dst[:n], &arr)
	return ok, nil
}
