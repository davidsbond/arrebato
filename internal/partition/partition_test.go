package partition_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/davidsbond/arrebato/internal/partition"
)

func TestPartitioner_Partition(t *testing.T) {
	p := partition.NewCRC32Partitioner()

	t.Run("It should produce the same value for the same message within the range of partitions", func(t *testing.T) {
		for i := uint32(1); i < 1000; i++ {
			m := structpb.NewStringValue(fmt.Sprint(i))
			first, err := p.Partition(m, i)
			require.NoError(t, err)
			second, err := p.Partition(m, i)
			require.NoError(t, err)

			assert.True(t, first < i)
			assert.True(t, second < i)
			assert.EqualValues(t, first, second)
		}
	})
}
