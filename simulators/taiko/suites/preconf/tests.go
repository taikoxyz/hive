package preconf

import (
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/suites"
)

var testSuite = hivesim.Suite{
	Name:        "preconf",
	DisplayName: "driver preconf test",
	Location:    "suites/preconf",
}

var Tests = make([]suites.TestSpec, 0)

func Suite(clients clients.ClientGroups) hivesim.Suite {
	suites.SuiteHydrate(&testSuite, clients, Tests)
	return testSuite
}
