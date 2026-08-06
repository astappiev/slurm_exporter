package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMemory(t *testing.T) {
	assert.Equal(t, 0.0, ParseMemory(""))
	assert.Equal(t, 0.0, ParseMemory("0"))
	assert.Equal(t, 4*1024.0, ParseMemory("4K"))
	assert.Equal(t, 4*1024*1024.0, ParseMemory("4M"))
	assert.Equal(t, 4*1024*1024*1024.0, ParseMemory("4G"))
	assert.Equal(t, 4*1024*1024*1024*1024.0, ParseMemory("4T"))
}

func TestUserMetricsFromJobs(t *testing.T) {
	jobs := []JobInfo{
		{User: "user-a", State: "PENDING", CPUs: 8},
		{User: "user-a", State: "RUNNING", CPUs: 4, Mem: 1024},
		{User: "user-a", State: "RUNNING", CPUs: 2, Mem: 512},
		{User: "user-b", State: "SUSPENDED"},
		{User: "user-b", State: "COMPLETED"},
		{User: "", State: "RUNNING", CPUs: 99},
	}

	um := userMetricsFromJobs(jobs)
	assert.Len(t, um, 2)
	assert.Equal(t, &UserJobMetrics{jobsPending: 1, cpusPending: 8, jobsRunning: 2, cpusRunning: 6, memRunning: 1536}, um["user-a"])
	assert.Equal(t, &UserJobMetrics{jobsSuspended: 1}, um["user-b"])
}
