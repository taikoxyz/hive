package suite_reorg

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"taiko/common/clients"
	execution_config "taiko/common/config/execution"
	taparams "taiko/params"
	"taiko/suites"
)

var testSuite = hivesim.Suite{
	Name:        "reorg",
	DisplayName: "Deneb reorg",
	Location:    "suites/reorg",
}

var Tests = make([]suites.TestSpec, 0)

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
		GenesisTimeDelay: 8,
		Genesis:          &genesis,
	})
	return testSuite
}
