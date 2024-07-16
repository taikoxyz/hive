package main

import (
	"github.com/ethereum/hive/hivesim"
)

const (
	envFile = "/taiko/.env"
	network = "hive_taiko_network"
)

var networkCreated = make(map[hivesim.SuiteID]bool)

// createNetwork ensures there is a separate network to be able to send the client traffic
// from two separate IP addrs.
func createAndConnectNetwork(t *hivesim.T, container string) {
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true
	}

	if err := t.Sim.ConnectContainer(t.SuiteID, network, container); err != nil {
		t.Fatal("can't connect container to network:", err)
	}
}
