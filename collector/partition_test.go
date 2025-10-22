package collector

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePartitionsMetrics(t *testing.T) {
	// Read the input data from a file
	sinfoFile, _ := os.Open("fixtures/sinfo/partition.txt")
	runningFile, _ := os.Open("fixtures/squeue/partition_running.txt")
	pendingFile, _ := os.Open("fixtures/squeue/partition_pending.txt")
	sinfoData, _ := io.ReadAll(sinfoFile)
	runningData, _ := io.ReadAll(runningFile)
	pendingData, _ := io.ReadAll(pendingFile)
	partitionMetrics := ParsePartitionMetrics(sinfoData, runningData, pendingData)

	assert.Equal(t, 273.0, partitionMetrics["ampere"].cpu.Allocated, "Miscount of allocated CPUs")
	assert.Equal(t, 279.0, partitionMetrics["ampere"].cpu.Idle, "Miscount of idle CPUs")
	assert.Equal(t, 480.0, partitionMetrics["ampere"].cpu.Other, "Miscount of other CPUs")
	assert.Equal(t, 1032.0, partitionMetrics["ampere"].cpu.Total, "Miscount of total CPUs")
	assert.Equal(t, 16.0, partitionMetrics["ampere"].jobsRunning, "Miscount of running jobs")
	assert.Equal(t, 30.0, partitionMetrics["ampere"].jobsPending, "Miscount of pending jobs")
}
