package collector

import (
	"testing"

	"github.com/astappiev/slurm_exporter/util"
	"github.com/stretchr/testify/assert"
)

func TestPartitionMetricsFromInfo(t *testing.T) {
	nodes := []NodeInfo{
		{Name: "n1", Partition: "partition-a", CPU: util.CPUS{Allocated: 4, Idle: 4, Total: 8}},
		{Name: "n2", Partition: "partition-a", CPU: util.CPUS{Allocated: 2, Idle: 6, Total: 8}},
		{Name: "n3", Partition: "partition-b", CPU: util.CPUS{Allocated: 8, Idle: 0, Total: 8}},
	}
	jobs := []JobInfo{
		{Partition: "partition-a", State: "PENDING"},
		{Partition: "partition-a", State: "PENDING"},
		{Partition: "partition-a", State: "RUNNING"},
		{Partition: "partition-b", State: "RUNNING"},
		{Partition: "unknown-partition", State: "RUNNING"},
	}

	pm := partitionMetricsFromInfo(nodes, jobs)
	assert.Len(t, pm, 2)
	assert.Equal(t, &PartitionMetrics{
		cpu:         util.CPUS{Allocated: 6, Idle: 10, Total: 16},
		jobsPending: 2,
		jobsRunning: 1,
	}, pm["partition-a"])
	assert.Equal(t, &PartitionMetrics{
		cpu:         util.CPUS{Allocated: 8, Idle: 0, Total: 8},
		jobsRunning: 1,
	}, pm["partition-b"])
}
