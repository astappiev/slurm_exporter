/*
	Copyright 2020 Joeri Hermans, Victor Penso, Matteo Dessalvi
	Copyright 2022 Iztok Lebar Bajec
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

type GPUsMetrics struct {
	alloc       float64
	idle        float64
	total       float64
	utilization float64
}

type GPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
	logger      *slog.Logger
}

func init() {
	registerCollector("gpus", defaultEnabled, NewGPUsCollector)
}

func NewGPUsCollector(logger *slog.Logger) (Collector, error) {
	return &GPUsCollector{
		logger:      logger,
		alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
	}, nil
}

func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) error {
	nodes, err := nodeCache.Get()
	if err != nil {
		return err
	}

	cm := gpusMetricsFromInfo(uniqueNodes(nodes))
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)

	return nil
}

func gpusMetricsFromInfo(nodes []NodeInfo) *GPUsMetrics {
	cm := &GPUsMetrics{}
	for _, node := range nodes {
		for _, gres := range node.GRES {
			if gres.Kind == "gpu" {
				cm.total += gres.Total
			}
		}
		for _, gres := range node.GRESUsed {
			if gres.Kind == "gpu" {
				cm.alloc += gres.Total
			}
		}
	}
	cm.idle = cm.total - cm.alloc
	if cm.total > 0 {
		cm.utilization = cm.alloc / cm.total
	}
	return cm
}
