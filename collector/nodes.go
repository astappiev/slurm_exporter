/*
	Copyright 2017 Victor Penso, Matteo Dessalvi
	Copyright 2024 Oleh Astappiev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package collector

import (
	"log/slog"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type NodesMetrics struct {
	alloc float64
	comp  float64
	down  float64
	drain float64
	err   float64
	fail  float64
	idle  float64
	maint float64
	mix   float64
	resv  float64
	plnd  float64
}

func nodesMetricsFromInfo(nodes []NodeInfo) *NodesMetrics {
	nm := &NodesMetrics{}
	for _, node := range nodes {
		state := strings.ToLower(node.State)
		switch {
		case strings.HasPrefix(state, "alloc"):
			nm.alloc++
		case strings.HasPrefix(state, "comp"):
			nm.comp++
		case strings.HasPrefix(state, "down"):
			nm.down++
		case strings.HasPrefix(state, "drain"):
			nm.drain++
		case strings.HasPrefix(state, "fail"):
			nm.fail++
		case strings.HasPrefix(state, "err"):
			nm.err++
		case strings.HasPrefix(state, "idle"):
			nm.idle++
		case strings.HasPrefix(state, "maint"):
			nm.maint++
		case strings.HasPrefix(state, "mix"):
			nm.mix++
		case strings.HasPrefix(state, "res"):
			nm.resv++
		case strings.HasPrefix(state, "plan"):
			nm.plnd++
		}
	}
	return nm
}

type NodesCollector struct {
	alloc  *prometheus.Desc
	comp   *prometheus.Desc
	down   *prometheus.Desc
	drain  *prometheus.Desc
	err    *prometheus.Desc
	fail   *prometheus.Desc
	idle   *prometheus.Desc
	maint  *prometheus.Desc
	mix    *prometheus.Desc
	resv   *prometheus.Desc
	plnd   *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("nodes", defaultEnabled, NewNodesCollector)
}

func NewNodesCollector(logger *slog.Logger) (Collector, error) {
	return &NodesCollector{
		logger: logger,
		alloc:  prometheus.NewDesc("slurm_nodes_alloc", "Allocated nodes", nil, nil),
		comp:   prometheus.NewDesc("slurm_nodes_comp", "Completing nodes", nil, nil),
		down:   prometheus.NewDesc("slurm_nodes_down", "Down nodes", nil, nil),
		drain:  prometheus.NewDesc("slurm_nodes_drain", "Drain nodes", nil, nil),
		err:    prometheus.NewDesc("slurm_nodes_err", "Error nodes", nil, nil),
		fail:   prometheus.NewDesc("slurm_nodes_fail", "Fail nodes", nil, nil),
		idle:   prometheus.NewDesc("slurm_nodes_idle", "Idle nodes", nil, nil),
		maint:  prometheus.NewDesc("slurm_nodes_maint", "Maint nodes", nil, nil),
		mix:    prometheus.NewDesc("slurm_nodes_mix", "Mix nodes", nil, nil),
		resv:   prometheus.NewDesc("slurm_nodes_resv", "Reserved nodes", nil, nil),
		plnd:   prometheus.NewDesc("slurm_nodes_plnd", "Planned nodes", nil, nil),
	}, nil
}

func (nc *NodesCollector) Collect(ch chan<- prometheus.Metric) error {
	nodes, err := nodeCache.Get()
	if err != nil {
		return err
	}

	nm := nodesMetricsFromInfo(uniqueNodes(nodes))
	ch <- prometheus.MustNewConstMetric(nc.alloc, prometheus.GaugeValue, nm.alloc)
	ch <- prometheus.MustNewConstMetric(nc.comp, prometheus.GaugeValue, nm.comp)
	ch <- prometheus.MustNewConstMetric(nc.down, prometheus.GaugeValue, nm.down)
	ch <- prometheus.MustNewConstMetric(nc.drain, prometheus.GaugeValue, nm.drain)
	ch <- prometheus.MustNewConstMetric(nc.err, prometheus.GaugeValue, nm.err)
	ch <- prometheus.MustNewConstMetric(nc.fail, prometheus.GaugeValue, nm.fail)
	ch <- prometheus.MustNewConstMetric(nc.idle, prometheus.GaugeValue, nm.idle)
	ch <- prometheus.MustNewConstMetric(nc.maint, prometheus.GaugeValue, nm.maint)
	ch <- prometheus.MustNewConstMetric(nc.mix, prometheus.GaugeValue, nm.mix)
	ch <- prometheus.MustNewConstMetric(nc.resv, prometheus.GaugeValue, nm.resv)
	ch <- prometheus.MustNewConstMetric(nc.plnd, prometheus.GaugeValue, nm.plnd)

	return nil
}
