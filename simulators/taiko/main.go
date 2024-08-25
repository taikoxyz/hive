package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	suite_base "taiko/suites/base"
	suite_reorg "taiko/suites/reorg"
)

func main() {
	sim := hivesim.New()
	// From the simulator we can get all client types provided
	clientTypes, err := sim.ClientTypes()
	if err != nil {
		panic(err)
	}

	// TODO
	data, _ := json.Marshal(clientTypes)
	fmt.Println("client types: ", string(data))

	clientsByRole := clients.GetClientsByRole(clientTypes)
	if clientsByRole == nil {
		panic("failed to create clients by role")
	}

	// TODO
	data, _ = json.Marshal(clientsByRole)
	fmt.Println("clients by role: ", string(data))

	hivesim.MustRunSuite(sim, suite_base.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_reorg.Suite(clientsByRole))
}
