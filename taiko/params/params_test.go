package params

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestENV(t *testing.T) {
	assert.Equal(t, true, len(EnvParams) > 0)
	assert.Equal(t, true, len(ContractTxs) > 0)

	var genesis core.Genesis
	// Load genesis.json
	if err := json.Unmarshal(GenesisContent, &genesis); err != nil {
		panic(err)
	}
	t.Log(genesis.Config.ChainID.Uint64())

	t.Logf(ContractTxs[0].Hash().String())
}
