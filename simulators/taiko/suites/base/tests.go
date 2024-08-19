package suite_base

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"taiko/common/clients"
	"taiko/common/config/execution"
	taparams "taiko/params"
	"taiko/suites"
)

var testSuite = hivesim.Suite{
	Name:        "taiko-deneb-testnet",
	DisplayName: "Deneb Testnet",
	Description: `Collection of test vectors that use a L1EthClient+BeaconNode+ValidatorClient testnet for Cancun+Deneb.`,
	Location:    "suites/base",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		BaseTestSpec{
			Name:        "test-deneb-genesis",
			DisplayName: "Deneb Genesis",
			Description: `
			Sanity test to check the beacon clients can start with deneb genesis.
			`,
			DenebGenesis: true,
			GenesisExecutionWithdrawalCredentialsShares: 1,
			WaitForFinality: true,
			NodeCount:       1,
			ValidatorCount:  5,
		},
	)
}

func Suite(clients clients.ClientsByRole) hivesim.Suite {
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

	suites.SuiteHydrate(&testSuite, clients, Tests, &execution_config.GenesisState{
		ForkName:         "deneb",
		BeaconConfig:     beaconConfig,
		NumValidators:    64,
		GenesisTimeDelay: 15,
		Genesis:          &genesis,
	})
	return testSuite
}
