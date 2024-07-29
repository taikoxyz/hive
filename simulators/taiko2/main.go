package main

import (
	"github.com/ethereum/hive/hivesim"
	"taiko2/common/clients"
	suite_base "taiko2/suites/base"
	suite_builder "taiko2/suites/builder"
	suite_blobs_gossip "taiko2/suites/p2p/gossip/blobs"
	suite_reorg "taiko2/suites/reorg"
	suite_sync "taiko2/suites/sync"
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
	// Mark suites for execution
	hivesim.MustRunSuite(sim, suite_base.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_sync.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_builder.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_reorg.Suite(clientsByRole))
	hivesim.MustRunSuite(sim, suite_blobs_gossip.Suite(clientsByRole))
}
