package suite_base

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	taparams "github.com/ethereum/hive/taiko/params"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"taiko2/common/clients"
	"taiko2/common/config/execution"
	"taiko2/suites"
)

var testSuite = hivesim.Suite{
	Name:        "eth2-deneb-testnet",
	DisplayName: "Deneb Testnet",
	Description: `Collection of test vectors that use a ExecutionClient+BeaconNode+ValidatorClient testnet for Cancun+Deneb.`,
	Location:    "suites/base",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		BaseTestSpec{
			Name:         "test-deneb-fork",
			DisplayName:  "Deneb Fork",
			Description:  `Sanity test to check the fork transition to deneb.`,
			DenebGenesis: false,
			GenesisExecutionWithdrawalCredentialsShares: 1,
			EpochsAfterFork:     1,
			ExitValidatorsShare: 10,
		},
		BaseTestSpec{
			Name:        "test-deneb-genesis",
			DisplayName: "Deneb Genesis",
			Description: `
			Sanity test to check the beacon clients can start with deneb genesis.
			`,
			DenebGenesis: true,
			GenesisExecutionWithdrawalCredentialsShares: 1,
			WaitForFinality:     true,
			ExitValidatorsShare: 10,
			NodeCount:           1,
			ValidatorCount:      5,
		},
	)
}

func Suite(c *clients.ClientDefinitionsByRole) hivesim.Suite {
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

	suites.SuiteHydrate(&testSuite, c, Tests, &execution_config.GenesisState{
		ForkName:         Deneb,
		BeaconConfig:     beaconConfig,
		NumValidators:    64,
		GenesisTimeDelay: 15,
		Genesis:          &genesis,
	})
	return testSuite
}
