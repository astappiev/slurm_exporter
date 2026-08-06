package collector

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sacct/jobhistory.txt")
	data, _ := io.ReadAll(file)
	metrics := ParseJobMetrics(nil, data)

	assert.Equal(t, 119, len(metrics))
	assert.Equal(t, "91193_761", metrics[0].JobID)
	assert.Equal(t, "extract", metrics[0].JobName)
	assert.Equal(t, "user1", metrics[0].User)
	assert.Equal(t, 70839.0, metrics[0].Elapsed)
	assert.Equal(t, 120243.0, metrics[68].Elapsed)
}
