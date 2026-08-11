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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type JobMetrics struct {
	pending     float64
	pendingCpus float64
	running     float64
	runningCpus float64
	suspended   float64
}

type AccountCollector struct {
	pending     *prometheus.Desc
	pendingCpus *prometheus.Desc
	running     *prometheus.Desc
	runningCpus *prometheus.Desc
	suspended   *prometheus.Desc
	logger      *slog.Logger
}

func init() {
	registerCollector("account", defaultEnabled, NewAccountCollector)
}

func NewAccountCollector(logger *slog.Logger) (Collector, error) {
	return &AccountCollector{
		logger:      logger,
		pending:     prometheus.NewDesc("slurm_account_jobs_pending", "Pending jobs for account", []string{"account"}, nil),
		pendingCpus: prometheus.NewDesc("slurm_account_cpus_pending", "Pending jobs for account", []string{"account"}, nil),
		running:     prometheus.NewDesc("slurm_account_jobs_running", "Running jobs for account", []string{"account"}, nil),
		runningCpus: prometheus.NewDesc("slurm_account_cpus_running", "Running cpus for account", []string{"account"}, nil),
		suspended:   prometheus.NewDesc("slurm_account_jobs_suspended", "Suspended jobs for account", []string{"account"}, nil),
	}, nil
}

func (ac *AccountCollector) Collect(ch chan<- prometheus.Metric) error {
	jobs, err := jobCache.Get()
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return ErrNoData
	}

	am := accountMetricsFromJobs(jobs)
	for a := range am {
		if am[a].pending > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pending, prometheus.GaugeValue, am[a].pending, a)
		}
		if am[a].pendingCpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pendingCpus, prometheus.GaugeValue, am[a].pendingCpus, a)
		}
		if am[a].running > 0 {
			ch <- prometheus.MustNewConstMetric(ac.running, prometheus.GaugeValue, am[a].running, a)
		}
		if am[a].runningCpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.runningCpus, prometheus.GaugeValue, am[a].runningCpus, a)
		}
		if am[a].suspended > 0 {
			ch <- prometheus.MustNewConstMetric(ac.suspended, prometheus.GaugeValue, am[a].suspended, a)
		}
	}

	return nil
}

func accountMetricsFromJobs(jobs []JobInfo) map[string]*JobMetrics {
	accounts := make(map[string]*JobMetrics)
	for _, job := range jobs {
		if job.Account == "" {
			continue
		}
		if accounts[job.Account] == nil {
			accounts[job.Account] = &JobMetrics{}
		}
		switch {
		case strings.HasPrefix(strings.ToLower(job.State), "pending"):
			accounts[job.Account].pending++
			accounts[job.Account].pendingCpus += job.CPUs
		case strings.HasPrefix(strings.ToLower(job.State), "running"):
			accounts[job.Account].running++
			accounts[job.Account].runningCpus += job.CPUs
		case strings.HasPrefix(strings.ToLower(job.State), "suspended"):
			accounts[job.Account].suspended++
		}
	}
	return accounts
}
