package collector

import (
	"io"
	"os"
	"testing"

	"github.com/astappiev/slurm_exporter/util"
	"github.com/stretchr/testify/assert"
)

func TestCPUsMetrics(t *testing.T) {
	file, _ := os.Open("fixtures/sinfo/cpus.txt")
	data, _ := io.ReadAll(file)
	cpus := util.ParseCPUS(string(data))

	assert.Equal(t, 5725.0, cpus.Allocated, "Miscount of alloc CPUs")
	assert.Equal(t, 877.0, cpus.Idle, "Miscount of idle CPUs")
	assert.Equal(t, 34.0, cpus.Other, "Miscount of other CPUs")
	assert.Equal(t, 6636.0, cpus.Total, "Miscount of total CPUs")
}
