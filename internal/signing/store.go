package signing

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/signing/v1"

	"go.etcd.io/bbolt"
)

type (
	// The BoltStore type is responsible for storing client public key data within a bolt database.
	BoltStore struct {
		db *bbolt.DB
	}
)

const (
	signingKey = "signing"
)

// NewBoltStore returns a new instance of the BoltStore type that stores signing key data within the provided
// bbolt.DB instance.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

var (
	// ErrNoPublicKey is the error given when requesting a public key that does not exist.
	ErrNoPublicKey = errors.New("no public key")

	// ErrPublicKeyExists is the error given when performing an operation that would overwrite an existing
	// public key.
	ErrPublicKeyExists = errors.New("public key exists")
)

// Get a client's public key. Returns ErrNoPublic key if it does not exist.
func (bs *BoltStore) Get(ctx context.Context, clientID string) ([]byte, error) {
	var publicKey []byte
	err := bs.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(signingKey))
		if bucket == nil {
			return ErrNoPublicKey
		}

		publicKey = bucket.Get([]byte(clientID))
		if publicKey == nil {
			return ErrNoPublicKey
		}

		return nil
	})

	return publicKey, err
}

// Create a public key record for a client. Returns ErrPublicKeyExists if the client already has a public key.
func (bs *BoltStore) Create(ctx context.Context, clientID string, publicKey []byte) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(signingKey))
		if err != nil {
			return fmt.Errorf("failed to open signing bucket: %w", err)
		}

		key := []byte(clientID)
		value := bucket.Get(key)
		if value != nil || len(value) > 0 {
			return ErrPublicKeyExists
		}

		if err = bucket.Put(key, publicKey); err != nil {
			return fmt.Errorf("failed to store public key: %w", err)
		}

		return nil
	})
}

// List all public keys stored in state.
func (bs *BoltStore) List(ctx context.Context) ([]*signing.PublicKey, error) {
	var keys []*signing.PublicKey
	err := bs.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(signingKey))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			keys = append(keys, &signing.PublicKey{
				ClientId:  string(k),
				PublicKey: v,
			})

			return nil
		})
	})

	return keys, err
}
