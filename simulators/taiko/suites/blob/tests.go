package blob

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
	Name:        "driver-blob-session-timeout",
	DisplayName: "driver blob session timeout",
	Location:    "suites/blob",
}

var Tests = make([]suites.TestSpec, 0)

func init() {
	Tests = append(Tests,
		BlobTestSpec{
			TestL1Beacon: true,
			BaseTestSpec: suite_base.BaseTestSpec{
				Name: "driver-l1-beacon",
			},
		},
		BlobTestSpec{
			TestBlobSocial: true,
			BaseTestSpec: suite_base.BaseTestSpec{
				Name: "driver-blob-social",
			},
		},
		BlobTestSpec{
			TestBlobServer: true,
			BaseTestSpec: suite_base.BaseTestSpec{
				Name: "driver-blob-server",
			},
		})
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
