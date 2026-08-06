package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountMetricsFromJobs(t *testing.T) {
	jobs := []JobInfo{
		{Account: "project-a", State: "PENDING", CPUs: 8},
		{Account: "project-a", State: "RUNNING", CPUs: 4},
		{Account: "project-a", State: "RUNNING", CPUs: 2},
		{Account: "project-b", State: "SUSPENDED", CPUs: 16},
		{Account: "project-b", State: "COMPLETED", CPUs: 16},
		{Account: "", State: "RUNNING", CPUs: 99},
	}

	am := accountMetricsFromJobs(jobs)
	assert.Len(t, am, 2)
	assert.Equal(t, &JobMetrics{pending: 1, pendingCpus: 8, running: 2, runningCpus: 6}, am["project-a"])
	assert.Equal(t, &JobMetrics{suspended: 1}, am["project-b"])
}
