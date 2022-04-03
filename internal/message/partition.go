package message

import (
	"fmt"
	"hash/crc32"
	"math"

	"google.golang.org/protobuf/proto"
)

type (
	// The CRC32Partitioner type is used to generate partitions for proto-encoded messages using the standard
	// library's crc32 hashing function.
	CRC32Partitioner struct {
		table *crc32.Table
	}
)

// NewCRC32Partitioner returns a new instance of the CRC32Partitioner type that initialises a crc32 table using
// the IEEE polynomial for all hashes generated.
func NewCRC32Partitioner() *CRC32Partitioner {
	return &CRC32Partitioner{
		// This can never change, otherwise you break partitioning for existing clusters.
		table: crc32.MakeTable(crc32.IEEE),
	}
}

// Partition returns a uint32 value denoting the expected partition of the provided proto.Message implementation. The
// message is deterministically marshalled and passed into the crc32 hashing function.
func (p *CRC32Partitioner) Partition(m proto.Message, max uint32) (uint32, error) {
	// Shortcut for topics that only have a single partition.
	if max == 1 {
		return 0, nil
	}

	// In order to ensure that we always choose the same partition for the same key/message content, we need to do
	// proto encoding deterministically, as a small difference in the binary representation can change the checksum
	// and then the designated partition.
	opts := proto.MarshalOptions{
		Deterministic: true,
	}

	value, err := opts.Marshal(m)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal message: %w", err)
	}

	checksum := crc32.Checksum(value, p.table)

	// We use math.Mod here to bring the generated checksum value into the range of 0 to max-1.
	return uint32(math.Mod(float64(checksum), float64(max-1))), nil
}
