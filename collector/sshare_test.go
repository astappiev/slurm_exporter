package collector

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFairShareMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("fixtures/sshare/fairshare.txt")
	data, _ := io.ReadAll(file)
	metrics := ParseFairShareMetrics(data)

	assert.Equal(t, 0.0, metrics["project-a"].fairshare)
	assert.Equal(t, 0.0, metrics["project-d"].fairshare)
}
