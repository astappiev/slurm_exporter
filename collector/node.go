/*
	Copyright 2021 Chris Read
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

	"github.com/prometheus/client_golang/prometheus"
)

type NodeCollector struct {
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
	gpuAlloc *prometheus.Desc
	gpuTotal *prometheus.Desc
	logger   *slog.Logger
}

func init() {
	registerCollector("node", defaultEnabled, NewNodeCollector)
}

// NewNodeCollector creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector(logger *slog.Logger) (Collector, error) {
	return &NodeCollector{
		logger:   logger,
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUS per node", []string{"node", "status"}, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUS per node", []string{"node", "status"}, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUS per node", []string{"node", "status"}, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUS per node", []string{"node", "status"}, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", []string{"node", "status"}, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", []string{"node", "status"}, nil),
		gpuAlloc: prometheus.NewDesc("slurm_node_gpu_alloc", "Allocated GPUs per node", []string{"node", "status", "gputype"}, nil),
		gpuTotal: prometheus.NewDesc("slurm_node_gpu_total", "Total GPUs per node", []string{"node", "status", "gputype"}, nil),
	}, nil
}

func (c *NodeCollector) Collect(ch chan<- prometheus.Metric) error {
	infos, err := nodeCache.Get()
	if err != nil {
		return err
	}

	for _, node := range uniqueNodes(infos) {
		ch <- prometheus.MustNewConstMetric(c.cpuAlloc, prometheus.GaugeValue, node.CPU.Allocated, node.Name, node.State)
		ch <- prometheus.MustNewConstMetric(c.cpuIdle, prometheus.GaugeValue, node.CPU.Idle, node.Name, node.State)
		ch <- prometheus.MustNewConstMetric(c.cpuOther, prometheus.GaugeValue, node.CPU.Other, node.Name, node.State)
		ch <- prometheus.MustNewConstMetric(c.cpuTotal, prometheus.GaugeValue, node.CPU.Total, node.Name, node.State)
		ch <- prometheus.MustNewConstMetric(c.memAlloc, prometheus.GaugeValue, node.MemAlloc, node.Name, node.State)
		ch <- prometheus.MustNewConstMetric(c.memTotal, prometheus.GaugeValue, node.MemTotal, node.Name, node.State)
		for _, tres := range node.GRESUsed {
			if tres.Kind == "gpu" {
				ch <- prometheus.MustNewConstMetric(c.gpuAlloc, prometheus.GaugeValue, tres.Total, node.Name, node.State, tres.Type)
			}
		}
		for _, tres := range node.GRES {
			if tres.Kind == "gpu" {
				ch <- prometheus.MustNewConstMetric(c.gpuTotal, prometheus.GaugeValue, tres.Total, node.Name, node.State, tres.Type)
			}
		}
	}

	return nil
}
