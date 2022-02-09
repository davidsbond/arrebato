package e2e_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"

	"github.com/davidsbond/arrebato/internal/server"
)

type (
	SnapshotSuite struct {
		suite.Suite

		server *server.Server
		config server.Config
	}
)

func (s *SnapshotSuite) TestSnapshotting() {
	buf := bytes.NewBuffer([]byte{})

	s.Run("It should produce a snapshot of the FSM", func() {
		snapshot, err := s.server.Snapshot()
		require.NoError(s.T(), err)

		require.NoError(s.T(), snapshot.Persist(&MockSnapshotSink{writer: buf}))
		assert.True(s.T(), buf.Len() > 0)
	})

	snapshotPath := filepath.Join(s.config.DataPath, "arrebato_restore.db")

	s.Run("It should restore from the snapshot as a local file", func() {
		require.NoError(s.T(), s.server.Restore(
			ioutil.NopCloser(buf),
		))

		_, err := os.Stat(snapshotPath)
		assert.NoError(s.T(), err)
	})

	s.Run("The snapshot file should be a valid bbolt database", func() {
		db, err := bbolt.Open(snapshotPath, 0o755, bbolt.DefaultOptions)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), db)
		require.NoError(s.T(), db.Close())
	})
}
