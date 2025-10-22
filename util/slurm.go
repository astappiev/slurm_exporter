package util

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// CPUS represents CPU allocation states in SLURM.
type CPUS struct {
	Allocated float64
	Idle      float64
	Other     float64
	Total     float64
}

// GRES represents a SLURM Generic Resource (GRES).
type GRES struct {
	Kind  string
	Type  string
	Total float64
}

// RunCommand executes a SLURM command and returns its output.
func RunCommand(executable string, arguments ...string) ([]byte, error) {
	subprocess := exec.Command(executable, arguments...)
	out, err := subprocess.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("run command error: %w", err)
	}
	return out, nil
}

// SplitLines splits input bytes into lines, handling both Unix and Windows line endings.
func SplitLines(input []byte) []string {
	return strings.Split(strings.ReplaceAll(string(input), "\r", ""), "\n")
}

// RemoveDuplicates removes duplicate strings from a slice while preserving order. Empty strings filtered out.
func RemoveDuplicates(s []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(s))

	for _, v := range s {
		if len(v) > 0 && !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}

	return result
}

// ParseCPUS parses SLURM CPU state string in format "alloc/idle/other/total".
func ParseCPUS(input string) CPUS {
	var cpus CPUS
	parts := strings.Split(strings.TrimSpace(input), "/")
	if len(parts) != 4 {
		return cpus
	}

	cpus.Allocated, _ = strconv.ParseFloat(parts[0], 64)
	cpus.Idle, _ = strconv.ParseFloat(parts[1], 64)
	cpus.Other, _ = strconv.ParseFloat(parts[2], 64)
	cpus.Total, _ = strconv.ParseFloat(parts[3], 64)
	return cpus
}

// ParseGRES parses SLURM GRES string in format "gpu:a100m40:8(IDX:0-7)".
func ParseGRES(input string) []GRES {
	if input == "" || !strings.Contains(input, ":") {
		return []GRES{}
	}

	// Remove parenthesis and all content inside
	for strings.Contains(input, "(") {
		openIndex := strings.Index(input, "(")
		closeIndex := strings.Index(input, ")")
		if closeIndex < openIndex {
			break
		}
		input = input[:openIndex] + input[closeIndex+1:]
	}

	var results []GRES
	for _, single := range strings.Split(input, ",") {
		parts := strings.Split(single, ":")
		if len(parts) < 2 {
			continue
		}

		gres := GRES{Kind: parts[0]}
		if len(parts) == 2 {
			gres.Total, _ = strconv.ParseFloat(parts[1], 64)
		} else if len(parts) >= 3 {
			gres.Type = parts[1]
			gres.Total, _ = strconv.ParseFloat(parts[2], 64)
		}
		results = append(results, gres)
	}

	return results
}
