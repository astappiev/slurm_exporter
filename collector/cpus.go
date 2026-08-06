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

	"github.com/astappiev/slurm_exporter/util"
	"github.com/prometheus/client_golang/prometheus"
)

type CPUsCollector struct {
	alloc  *prometheus.Desc
	idle   *prometheus.Desc
	other  *prometheus.Desc
	total  *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("cpus", defaultEnabled, NewCPUsCollector)
}

func NewCPUsCollector(logger *slog.Logger) (Collector, error) {
	return &CPUsCollector{
		logger: logger,
		alloc:  prometheus.NewDesc("slurm_cpus_alloc", "Allocated CPUS", nil, nil),
		idle:   prometheus.NewDesc("slurm_cpus_idle", "Idle CPUS", nil, nil),
		other:  prometheus.NewDesc("slurm_cpus_other", "Mix CPUS", nil, nil),
		total:  prometheus.NewDesc("slurm_cpus_total", "Total CPUS", nil, nil),
	}, nil
}

func (cc *CPUsCollector) Collect(ch chan<- prometheus.Metric) error {
	nodes, err := nodeCache.Get()
	if err != nil {
		return err
	}

	cm := cpusMetricsFromInfo(uniqueNodes(nodes))
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.Allocated)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.Idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.Other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.Total)

	return nil
}

func cpusMetricsFromInfo(nodes []NodeInfo) util.CPUS {
	var cm util.CPUS
	for _, node := range nodes {
		cm.Allocated += node.CPU.Allocated
		cm.Idle += node.CPU.Idle
		cm.Other += node.CPU.Other
		cm.Total += node.CPU.Total
	}
	return cm
}
