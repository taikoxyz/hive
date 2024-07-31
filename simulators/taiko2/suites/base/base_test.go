package suite_base

import (
	"context"
	"github.com/stretchr/testify/assert"
	execution_config "taiko2/common/config/execution"
	"testing"
)

func TestGenesis(t *testing.T) {
	ctx := context.Background()

	err := execution_config.CliActionGenerateGenesisState(ctx, &execution_config.GenesisState{
		ChainConfigFile:    "./config.yml",
		GethGenesisJsonIn:  "./genesis.json",
		GethGenesisJsonOut: "./genesis.json",
		OutputSSZ:          "./genesis.ssz",
	})
	assert.NoError(t, err)
}
