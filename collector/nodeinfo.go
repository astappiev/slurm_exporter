package collector

import (
	"strconv"
	"strings"
	"time"

	"github.com/astappiev/slurm_exporter/util"
)

// NodeInfo is the shared per-node, per-partition representation returned by sinfo.
type NodeInfo struct {
	Name, Partition, State string
	CPU                    util.CPUS
	MemAlloc, MemTotal     float64
	GRES, GRESUsed         []util.GRES
}

func parseNodeInfo(input []byte) []NodeInfo {
	var result []NodeInfo
	for _, line := range util.SplitLines(input) {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		alloc, _ := strconv.ParseFloat(fields[4], 64)
		total, _ := strconv.ParseFloat(fields[5], 64)
		info := NodeInfo{Name: fields[0], Partition: fields[1], State: fields[2], CPU: util.ParseCPUS(fields[3]), MemAlloc: alloc, MemTotal: total}
		if len(fields) > 6 && fields[6] != "(null)" {
			info.GRES = util.ParseGRES(fields[6])
		}
		if len(fields) > 7 && fields[7] != "(null)" {
			info.GRESUsed = util.ParseGRES(fields[7])
		}
		result = append(result, info)
	}
	return result
}

func fetchNodeInfo() ([]NodeInfo, error) {
	out, err := util.RunCommand(*cliTimeout, "sinfo", "-h", "-a", "-N", "-O", "NodeList: ,Partition: ,StateLong: ,CPUsState: ,AllocMem: ,Memory: ,GRES: ,Gresused: ")
	if err != nil {
		return nil, err
	}
	return parseNodeInfo(out), nil
}

var nodeCache = newThrottledCache(fetchNodeInfo, func() time.Duration { return *defaultPollInterval })

func uniqueNodes(nodes []NodeInfo) []NodeInfo {
	seen := make(map[string]struct{}, len(nodes))
	result := make([]NodeInfo, 0, len(nodes))
	for _, node := range nodes {
		if _, ok := seen[node.Name]; ok {
			continue
		}
		seen[node.Name] = struct{}{}
		result = append(result, node)
	}
	return result
}
