package suite_base

import (
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/suites"
)

var testSuite = hivesim.Suite{
	Name:        "base",
	DisplayName: "Deneb Testnet",
	Description: `Collection of test vectors that use a L1EthClient+BeaconNode+ValidatorClient testnet for Cancun+Deneb.`,
	Location:    "suites/base",
}

var Tests = make([]suites.TestSpec, 0)

func Suite(clients clients.ClientGroups) hivesim.Suite {
	suites.SuiteHydrate(&testSuite, clients, Tests)
	return testSuite
}
