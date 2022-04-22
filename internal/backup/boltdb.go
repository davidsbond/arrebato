// Package backup provides types responsible for performing backups of various aspects of the server.
package backup

import (
	"io"

	"go.etcd.io/bbolt"
)

type (
	// BoltDB is an io.WriterTo implementation that writes the contents of a boltdb database.
	BoltDB struct {
		db *bbolt.DB
	}
)

// NewBoltDB returns a new instance of the BoltDB type which can be used to take backups of a boltdb
// database.
func NewBoltDB(db *bbolt.DB) *BoltDB {
	return &BoltDB{db: db}
}

// WriteTo writes the contents of the boltdb instance to the io.Writer implementation.
func (b *BoltDB) WriteTo(writer io.Writer) (int64, error) {
	tx, err := b.db.Begin(false)
	if err != nil {
		return 0, err
	}

	return tx.WriteTo(writer)
}
