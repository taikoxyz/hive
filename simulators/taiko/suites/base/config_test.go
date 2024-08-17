package suite_base

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/stretchr/testify/assert"
	"taiko/common/clients"
	execution_config "taiko/common/config/execution"
	"taiko/common/testnet"
	taparams "taiko/params"
	"testing"
)

func TestCC(t *testing.T) {
	var clientTypes = []*hivesim.ClientDefinition{
		{
			Name:    "go-ethereum",
			Version: "client-1-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"anvil"}},
		},
	}

	clientsByRole := clients.GetClientsByRole(clientTypes)
	_ = clientsByRole

	// Load params.yml
	beaconConfig, err := params.UnmarshalConfig(taparams.ConfigContent, nil)
	if err != nil {
		panic(err)
	}

	var genesis core.Genesis
	// Load genesis.json
	if err = json.Unmarshal(taparams.GenesisContent, &genesis); err != nil {
		panic(err)
	}

	for _, test := range Tests {
		genesisState := &execution_config.GenesisState{
			ForkName:         "deneb",
			BeaconConfig:     beaconConfig,
			NumValidators:    64,
			GenesisTimeDelay: 15,
			Genesis:          &genesis,
		}

		prep, err := testnet.PrepareTestnet(test.GetTestnetConfig(), genesisState)
		_ = prep
		assert.NoError(t, err)
	}
}
