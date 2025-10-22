package collector

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGPUsMetrics(t *testing.T) {
	testDataPaths, _ := filepath.Glob("fixtures/sinfo/slurm-*")
	for _, testDataPath := range testDataPaths {
		slurmVersion := strings.TrimPrefix(testDataPath, "fixtures/sinfo/slurm-")
		t.Logf("slurm-%s", slurmVersion)

		file, _ := os.Open(testDataPath + "/gpus.txt")
		data, _ := io.ReadAll(file)

		metrics := ParseGPUsMetrics(data)
		assert.Equal(t, 41.0, metrics.idle)
		assert.Equal(t, 7.0, metrics.alloc)
		assert.Equal(t, 48.0, metrics.total)
	}
}
