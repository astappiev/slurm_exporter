package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueMetricsFromJobs(t *testing.T) {
	jobs := []JobInfo{
		{State: "PENDING", Reason: "Dependency"},
		{State: "PENDING", Reason: "Resources"},
		{State: "RUNNING"},
		{State: "SUSPENDED"},
		{State: "CANCELLED"},
		{State: "COMPLETING"},
		{State: "COMPLETED"},
		{State: "CONFIGURING"},
		{State: "FAILED"},
		{State: "TIMEOUT"},
		{State: "PREEMPTED"},
		{State: "NODE_FAIL"},
		{State: "OUT_OF_MEMORY"},
	}

	qm := queueMetricsFromJobs(jobs)
	assert.Equal(t, &QueueMetrics{
		pending: 2, pendingDep: 1, running: 1, suspended: 1, cancelled: 1,
		completing: 1, completed: 1, configuring: 1, failed: 1, timeout: 1,
		preempted: 1, nodeFail: 1, outOfMemory: 1,
	}, qm)
}
