package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCPUS(t *testing.T) {
	assert.Equal(t, CPUS{}, ParseCPUS(""))
	assert.Equal(t, CPUS{Allocated: 64, Idle: 0, Other: 64, Total: 128}, ParseCPUS("64/0/64/128"))
	assert.Equal(t, CPUS{Allocated: 16, Idle: 16, Other: 224, Total: 256}, ParseCPUS(" 16/16/224/256"))
}

func TestParseGRES(t *testing.T) {
	assert.Equal(t, []GRES{}, ParseGRES(""))
	assert.Equal(t, []GRES{{Kind: "gpu", Type: "", Total: 8}}, ParseGRES("gpu:8"))
	assert.Equal(t, []GRES{{Kind: "gpu", Type: "tesla", Total: 2}, {Kind: "gpu", Type: "kepler", Total: 2}, {Kind: "mps", Total: 400}}, ParseGRES("gpu:tesla:2,gpu:kepler:2,mps:400"))
	assert.Equal(t, []GRES{{Kind: "gpu", Type: "a100m40", Total: 8}}, ParseGRES("gpu:a100m40:8"))
	assert.Equal(t, []GRES{{Kind: "gpu", Type: "a100m40", Total: 4}, {Kind: "gpu", Type: "a100m80", Total: 2}}, ParseGRES("gpu:a100m40:4(IDX:0-3),gpu:a100m80:2(IDX:4-5)"))
	assert.Equal(t, []GRES{{Kind: "gpu", Total: 8}}, ParseGRES("gpu:8(S:0-1)"))
	assert.Equal(t, []GRES{{Kind: "gpu", Total: 3}}, ParseGRES("gpu:(null):3(IDX:0-7)"))
	assert.Equal(t, []GRES{{Kind: "gpu", Type: "A30", Total: 4}, {Kind: "gpu", Type: "Q6K", Total: 4}}, ParseGRES("gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)"))
}
