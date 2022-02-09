package e2e_test

import "io"

type (
	MockSnapshotSink struct {
		writer io.Writer
	}
)

func (m *MockSnapshotSink) Write(p []byte) (n int, err error) {
	return m.writer.Write(p)
}

func (m *MockSnapshotSink) Close() error {
	return nil
}

func (m *MockSnapshotSink) ID() string {
	return ""
}

func (m *MockSnapshotSink) Cancel() error {
	return nil
}
