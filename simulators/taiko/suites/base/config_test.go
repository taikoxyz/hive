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

	clientsByRole := clients.ClientsByRole(clientTypes)

	clientCombinations := clientsByRole.Combinations()

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

	env := &testnet.Environment{
		Clients: clientsByRole,
	}

	for _, test := range Tests {
		config := test.GetTestnetConfig(clientCombinations)

		for _, node := range config.NodeDefinitions {
			var (
				anvilDef     = env.Clients.ClientByNameAndRole(node.L1EthClient, "anvil")
				executionDef = env.Clients.ClientByNameAndRole(node.L1EthClient, "eth1")
				beaconDef    = env.Clients.ClientByNameAndRole(node.ConsensusClient, "beacon")
				validatorDef = env.Clients.ClientByNameAndRole(node.ValidatorClientName(), "validator")
				taikoGethDef = env.Clients.ClientByNameAndRole(node.L2EthClient, "taiko-geth")
				driverDef    = env.Clients.ClientByNameAndRole(node.DriverClient, "driver")
				proposerDef  = env.Clients.ClientByNameAndRole(node.ProposerClient, "proposer")
				proverDef    = env.Clients.ClientByNameAndRole(node.ProverClient, "prover")
			)
			t.Log(executionDef, beaconDef, validatorDef, taikoGethDef, driverDef, proposerDef, proverDef, anvilDef)
		}

		genesisState := &execution_config.GenesisState{
			ForkName:         "deneb",
			BeaconConfig:     beaconConfig,
			NumValidators:    64,
			GenesisTimeDelay: 15,
			Genesis:          &genesis,
		}

		prep, err := testnet.PrepareTestnet(config, genesisState)
		assert.NoError(t, err)
		t.Log(prep)
	}
}
