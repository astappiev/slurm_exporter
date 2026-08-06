package collector

import (
	"testing"

	"github.com/astappiev/slurm_exporter/util"
	"github.com/stretchr/testify/assert"
)

func TestCPUsMetricsFromInfo(t *testing.T) {
	nodes := []NodeInfo{
		{Name: "n1", CPU: util.CPUS{Allocated: 56, Idle: 72, Other: 0, Total: 128}},
		{Name: "n2", CPU: util.CPUS{Allocated: 32, Idle: 480, Other: 0, Total: 512}},
	}

	cm := cpusMetricsFromInfo(nodes)
	assert.Equal(t, util.CPUS{Allocated: 88, Idle: 552, Other: 0, Total: 640}, cm)
}
