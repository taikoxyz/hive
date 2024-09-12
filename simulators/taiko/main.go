package main

import (
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	suite_base "taiko/suites/base"
	suite_blob "taiko/suites/blob"
	suite_reorg "taiko/suites/reorg"
)

func main() {
	sim := hivesim.New()
	// From the simulator we can get all client types provided
	clientTypes, err := sim.ClientTypes()
	if err != nil {
		panic(err)
	}

	clientsByRole := clients.GetClientsByRole(clientTypes)
	if clientsByRole == nil {
		panic("failed to create clients by role")
	}

	hivesim.MustRunSuite(sim, suite_base.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_reorg.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_blob.Suite(clientsByRole))
}
