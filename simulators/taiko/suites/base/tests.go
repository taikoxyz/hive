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
	Name:        "taiko-genesis",
	DisplayName: "Deneb Testnet",
	Description: `Collection of test vectors that use a L1EthClient+BeaconNode+ValidatorClient testnet for Cancun+Deneb.`,
	Location:    "suites/base",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		BaseTestSpec{
			Name:           "l2-full-sync",
			DisplayName:    "Deneb Genesis",
			NodeCount:      1,
			ValidatorCount: 5,
			L2SyncMode:     "full",
		},
		BaseTestSpec{
			Name:           "l2-snap-sync",
			DisplayName:    "Deneb Genesis",
			NodeCount:      1,
			ValidatorCount: 5,
			L2SyncMode:     "snap",
		},
	)
}

func Suite(clients clients.ClientGroups) hivesim.Suite {
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
