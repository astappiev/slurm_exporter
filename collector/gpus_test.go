package collector

import (
	"testing"

	"github.com/astappiev/slurm_exporter/util"
	"github.com/stretchr/testify/assert"
)

func TestGPUsMetricsFromInfo(t *testing.T) {
	nodes := []NodeInfo{
		{
			Name: "gpu-node-01",
			GRES: []util.GRES{
				{Kind: "gpu", Type: "a100", Total: 8},
				{Kind: "cpu", Type: "core", Total: 128},
			},
			GRESUsed: []util.GRES{
				{Kind: "gpu", Type: "a100", Total: 3},
			},
		},
		{
			Name: "gpu-node-02",
			GRES: []util.GRES{
				{Kind: "gpu", Type: "h100", Total: 4},
			},
			GRESUsed: []util.GRES{
				{Kind: "gpu", Type: "h100", Total: 4},
			},
		},
	}

	cm := gpusMetricsFromInfo(nodes)
	assert.Equal(t, &GPUsMetrics{alloc: 7, idle: 5, total: 12, utilization: 7.0 / 12.0}, cm)
}

func TestGPUsMetricsFromInfoNoGPUs(t *testing.T) {
	cm := gpusMetricsFromInfo([]NodeInfo{{Name: "cpu-only"}})
	assert.Equal(t, &GPUsMetrics{}, cm)
}
