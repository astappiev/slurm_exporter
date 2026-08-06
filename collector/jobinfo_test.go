package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJobInfo(t *testing.T) {
	data, err := os.ReadFile("fixtures/squeue/jobinfo.txt")
	assert.NoError(t, err)

	jobs := parseJobInfo(data)
	assert.Len(t, jobs, 6)
	assert.Equal(t, JobInfo{ID: "10001", Account: "project-a", User: "user-a", State: "PENDING", CPUs: 8, Mem: 4 * 1024 * 1024 * 1024, Partition: "partition-a", Reason: "QOSMaxGRESPerUser"}, jobs[0])
	assert.Equal(t, "Dependency", jobs[1].Reason)
}
