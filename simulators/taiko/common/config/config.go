package config

import (
	"bytes"
	"io"
	"math/big"
)

type ForkConfig struct {
	TerminalTotalDifficulty *big.Int `json:"terminal_total_difficulty,omitempty"`
	AltairForkEpoch         *big.Int `json:"altair_fork_epoch,omitempty"`
	BellatrixForkEpoch      *big.Int `json:"bellatrix_fork_epoch,omitempty"`
	CapellaForkEpoch        *big.Int `json:"capella_fork_epoch,omitempty"`
	DenebForkEpoch          *big.Int `json:"deneb_fork_epoch,omitempty"`
}

func BytesSource(data []byte) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
}
