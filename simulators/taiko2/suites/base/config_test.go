package suite_base

import (
	"github.com/ethereum/hive/hivesim"
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
	}

	clientsByRole := clients.ClientsByRole(clientTypes)

	clientCombinations := clientsByRole.Combinations()

	mnemonic := "couple kiwi radio river setup fortune hunt grief buddy forward perfect empty slim wear bounce drift execute nation tobacco dutch chapter festival ice fog"

	for _, test := range Tests {
		keys := test.GetValidatorKeys(mnemonic)
		env := &testnet.Environment{
			Clients:    clientsByRole,
			Validators: keys,
		}
		config := test.GetTestnetConfig(clientCombinations)

		genesisState := &execution_config.GenesisState{
			ForkName:           "deneb",
			NumValidators:      64,
			GenesisTimeDelay:   15,
			ChainConfigFile:    "/Users/huan/projects/taiko/hive/simulators/taiko2/config.yml",
			OutputSSZ:          "/Users/huan/projects/taiko/hive/simulators/taiko2/genesis.ssz",
			GethGenesisJsonIn:  "/Users/huan/projects/taiko/hive/simulators/taiko2/genesis.json",
			GethGenesisJsonOut: "/Users/huan/projects/taiko/hive/simulators/taiko2/genesis.json",
		}

		prep, err := testnet.PrepareTestnet(env, config, genesisState)
		assert.NoError(t, err)
		t.Log(prep)
	}
}
