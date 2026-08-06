package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodesMetricsFromInfo(t *testing.T) {
	nodes := []NodeInfo{
		{Name: "n1", State: "allocated"},
		{Name: "n2", State: "completing"},
		{Name: "n3", State: "down"},
		{Name: "n4", State: "draining"},
		{Name: "n5", State: "failing"},
		{Name: "n6", State: "error"},
		{Name: "n7", State: "idle"},
		{Name: "n8", State: "maintenance"},
		{Name: "n9", State: "mixed"},
		{Name: "n10", State: "reserved"},
		{Name: "n11", State: "planned"},
		{Name: "n12", State: "unknown"},
	}

	nm := nodesMetricsFromInfo(nodes)
	assert.Equal(t, &NodesMetrics{
		alloc: 1, comp: 1, down: 1, drain: 1, fail: 1, err: 1,
		idle: 1, maint: 1, mix: 1, resv: 1, plnd: 1,
	}, nm)
}
