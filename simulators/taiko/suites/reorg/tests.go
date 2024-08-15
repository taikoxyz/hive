package suite_reorg

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"taiko/common/clients"
	execution_config "taiko/common/config/execution"
	taparams "taiko/params"
	"taiko/suites"
	suite_base "taiko/suites/base"
)

var testSuite = hivesim.Suite{
	Name:        "taiko-deneb-reorg",
	DisplayName: "Deneb reorg",
	Description: "",
	Location:    "suites/reorg",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		ReorgTestSpec{
			BaseTestSpec: suite_base.BaseTestSpec{
				Name:        "test-deneb-reorg",
				DisplayName: "Deneb Reorg",
				Description: `
			Reorg l1eth and test taiko-client work normally.
			`,
				DenebGenesis: true,
				GenesisExecutionWithdrawalCredentialsShares: 1,
				WaitForFinality: true,
				NodeCount:       1,
				ValidatorCount:  5,
			},
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
		GenesisTimeDelay: 3,
		Genesis:          &genesis,
	})
	return testSuite
}
