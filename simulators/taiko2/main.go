package main

import (
	"encoding/json"
	"fmt"
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

	data, _ := json.Marshal(clientTypes)
	fmt.Println("Client types available: ", string(data))

	clientsByRole := clients.ClientsByRole(clientTypes)
	if clientsByRole == nil {
		panic("failed to create clients by role")
	}
	// Mark suites for execution

	data, _ = json.Marshal(clientsByRole)
	fmt.Println("clients by role: ", string(data))

	hivesim.MustRunSuite(sim, suite_base.Suite(clientsByRole))
	//hivesim.MustRunSuite(sim, suite_sync.Suite(clientsByRole))
	//hivesim.MustRunSuite(sim, suite_builder.Suite(clientsByRole))
	//hivesim.MustRunSuite(sim, suite_reorg.Suite(clientsByRole))
	//hivesim.MustRunSuite(sim, suite_blobs_gossip.Suite(clientsByRole))
}
