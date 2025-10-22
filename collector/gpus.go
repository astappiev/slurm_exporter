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
	"strconv"
	"strings"

	"github.com/astappiev/slurm_exporter/util"
	"github.com/prometheus/client_golang/prometheus"
)

type GPUsMetrics struct {
	alloc       float64
	idle        float64
	total       float64
	utilization float64
}

func ParseGPUsMetrics(data []byte) *GPUsMetrics {
	var usedGpus = 0.0
	var totalGpus = 0.0
	for _, line := range util.SplitLines(data) {
		if len(line) > 0 && strings.Contains(line, "gpu:") {
			parts := strings.Fields(line)
			numNodes, _ := strconv.ParseFloat(parts[0], 64)
			if numNodes == 0 { // for old slurm the results were not grouped and `NodeName` was shown in first column
				numNodes = 1
			}
			gres := util.ParseGRES(parts[1])
			gresUsed := util.ParseGRES(parts[2])

			numNodeGpus := 0.0
			for _, gres := range gres {
				if gres.Kind == "gpu" {
					numNodeGpus += gres.Total
				}
			}
			totalGpus += numNodes * numNodeGpus

			numNodeUsedGpus := 0.0
			for _, gres := range gresUsed {
				if gres.Kind == "gpu" {
					numNodeUsedGpus += gres.Total
				}
			}
			usedGpus += numNodes * numNodeUsedGpus
		}
	}

	return &GPUsMetrics{
		alloc:       usedGpus,
		idle:        totalGpus - usedGpus,
		total:       totalGpus,
		utilization: usedGpus / totalGpus,
	}
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
	out, err := util.RunCommand("sinfo", "-a", "-h", "--Format=Nodes: ,GRES: ,GresUsed:")
	if err != nil {
		return err
	}

	cm := ParseGPUsMetrics(out)
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)

	return nil
}
