package main

import (
	"github.com/ethereum/hive/hivesim"
	"taiko2/common/clients"
	suite_base "taiko2/suites/base"
)

func main() {
	sim := hivesim.New()
	// From the simulator we can get all client types provided
	clientTypes, err := sim.ClientTypes()
	if err != nil {
		panic(err)
	}

	clientsByRole := clients.ClientsByRole(clientTypes)
	if clientsByRole == nil {
		panic("failed to create clients by role")
	}

	hivesim.MustRunSuite(sim, suite_base.Suite(clientsByRole))
}
