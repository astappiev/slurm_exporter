package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeInfoDedupAndPartitionAggregation(t *testing.T) {
	data, err := os.ReadFile("fixtures/sinfo/nodeinfo.txt")
	assert.NoError(t, err)
	nodes := parseNodeInfo(data)
	unique := uniqueNodes(nodes)
	assert.Len(t, unique, 3)
	assert.Equal(t, 768.0, unique[0].CPU.Total+unique[1].CPU.Total+unique[2].CPU.Total)

	partitions := partitionMetricsFromInfo(nodes, nil)
	assert.Equal(t, 256.0, partitions["partition-a"].cpu.Total)
	assert.Equal(t, 768.0, partitions["fallback*"].cpu.Total)
}
