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
	Name:        "taiko-reorg",
	DisplayName: "Deneb reorg",
	Location:    "suites/reorg",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		ReorgTestSpec{
			BaseTestSpec: suite_base.BaseTestSpec{
				Name:       "taiko-reorg",
				L2SyncMode: "full",
			},
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
