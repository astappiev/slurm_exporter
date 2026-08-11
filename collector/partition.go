/*
	Copyright 2020 Victor Penso
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

type PartitionMetrics struct {
	cpu         util.CPUS
	jobsPending float64
	jobsRunning float64
}

type PartitionCollector struct {
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	running   *prometheus.Desc
	total     *prometheus.Desc
	logger    *slog.Logger
}

func init() {
	registerCollector("partition", defaultEnabled, NewPartitionCollector)
}

func NewPartitionCollector(logger *slog.Logger) (Collector, error) {
	return &PartitionCollector{
		logger:    logger,
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUS for partition", []string{"partition"}, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUS for partition", []string{"partition"}, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUS for partition", []string{"partition"}, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", []string{"partition"}, nil),
		running:   prometheus.NewDesc("slurm_partition_jobs_running", "Running jobs for partition", []string{"partition"}, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUS for partition", []string{"partition"}, nil),
	}, nil
}

func (pc *PartitionCollector) Collect(ch chan<- prometheus.Metric) error {
	nodes, err := nodeCache.Get()
	if err != nil {
		return err
	}
	jobs, err := jobCache.Get()
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return ErrNoData
	}

	pm := partitionMetricsFromInfo(nodes, jobs)
	for p := range pm {
		if pm[p].cpu.Allocated > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].cpu.Allocated, p)
		}
		if pm[p].cpu.Idle > 0 {
			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].cpu.Idle, p)
		}
		if pm[p].cpu.Other > 0 {
			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].cpu.Other, p)
		}
		if pm[p].cpu.Total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].cpu.Total, p)
		}
		if pm[p].jobsPending > 0 {
			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].jobsPending, p)
		}
		if pm[p].jobsRunning > 0 {
			ch <- prometheus.MustNewConstMetric(pc.running, prometheus.GaugeValue, pm[p].jobsRunning, p)
		}
	}

	return nil
}

func partitionMetricsFromInfo(nodes []NodeInfo, jobs []JobInfo) map[string]*PartitionMetrics {
	partitions := make(map[string]*PartitionMetrics)
	for _, node := range nodes {
		if partitions[node.Partition] == nil {
			partitions[node.Partition] = &PartitionMetrics{}
		}
		p := partitions[node.Partition]
		p.cpu.Allocated += node.CPU.Allocated
		p.cpu.Idle += node.CPU.Idle
		p.cpu.Other += node.CPU.Other
		p.cpu.Total += node.CPU.Total
	}
	for _, job := range jobs {
		if partitions[job.Partition] == nil {
			continue
		}
		if job.State == "PENDING" {
			partitions[job.Partition].jobsPending++
		}
		if job.State == "RUNNING" {
			partitions[job.Partition].jobsRunning++
		}
	}
	return partitions
}
