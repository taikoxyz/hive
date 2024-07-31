package prysm_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"taiko2/common/config/execution"
	"testing"
)

func TestGenesis(t *testing.T) {
	genesisState := &execution_config.GenesisState{}

	ctx := context.Background()

	err := execution_config.CliActionGenerateGenesisState(ctx, genesisState)
	assert.NoError(t, err)
}
