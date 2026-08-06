package collector

import (
	"strconv"
	"strings"
	"time"

	"github.com/astappiev/slurm_exporter/util"
)

type JobInfo struct {
	ID, Account, User, State string
	CPUs, Mem                float64
	Partition, Reason        string
}

func fetchJobInfo() ([]JobInfo, error) {
	out, err := util.RunCommand(*cliTimeout, "squeue", "-a", "-r", "-h", "-o", "%A|%a|%u|%T|%C|%m|%P|%r", "--states=all")
	if err != nil {
		return nil, err
	}
	return parseJobInfo(out), nil
}

func parseJobInfo(out []byte) []JobInfo {
	var result []JobInfo
	for _, line := range util.SplitLines(out) {
		parts := strings.Split(line, "|")
		if len(parts) != 8 {
			continue
		}
		cpus, _ := strconv.ParseFloat(parts[4], 64)
		result = append(result, JobInfo{ID: parts[0], Account: parts[1], User: parts[2], State: parts[3], CPUs: cpus, Mem: ParseMemory(parts[5]), Partition: parts[6], Reason: parts[7]})
	}
	return result
}

var jobCache = newThrottledCache(fetchJobInfo, func() time.Duration { return *defaultPollInterval })
