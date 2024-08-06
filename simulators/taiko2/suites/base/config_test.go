package suite_base

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	taparams "github.com/ethereum/hive/taiko/params"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/stretchr/testify/assert"
	"taiko2/common/clients"
	execution_config "taiko2/common/config/execution"
	"taiko2/common/testnet"
	"testing"
)

func TestCC(t *testing.T) {
	var clientTypes = []*hivesim.ClientDefinition{
		{
			Name:    "go-ethereum",
			Version: "client-1-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"eth1"}},
		},
		{
			Name:    "prysm-bn",
			Version: "client-2-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"beacon"}},
		}, {
			Name:    "prysm-vc",
			Version: "client-3-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"validator"}},
		},
		{
			Name:    "taiko-geth",
			Version: "taiko-geth-4-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"taiko-geth"}},
		},
		{
			Name:    "proposer",
			Version: "proposer-5-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"proposer"}},
		},
		{
			Name:    "driver",
			Version: "driver-6-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"driver"}},
		},
		{
			Name:    "prover",
			Version: "prover-7-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"prover"}},
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
				executionDef = env.Clients.ClientByNameAndRole(node.L1EthClient, "eth1")
				beaconDef    = env.Clients.ClientByNameAndRole(node.ConsensusClient, "beacon")
				validatorDef = env.Clients.ClientByNameAndRole(node.ValidatorClientName(), "validator")
				taikoGethDef = env.Clients.ClientByNameAndRole(node.L2EthClient, "taiko-geth")
				driverDef    = env.Clients.ClientByNameAndRole(node.DriverClient, "driver")
				proposerDef  = env.Clients.ClientByNameAndRole(node.ProposerClient, "proposer")
				proverDef    = env.Clients.ClientByNameAndRole(node.ProverClient, "prover")
			)
			t.Log(executionDef, beaconDef, validatorDef, taikoGethDef, driverDef, proposerDef, proverDef)
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
