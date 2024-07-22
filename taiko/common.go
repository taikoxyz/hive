package taiko

import "github.com/ethereum/hive/hivesim"

var networkCreated = make(map[hivesim.SuiteID]bool)

// CreateOrConnectNetwork ensures there is a separate network to be able to send the client traffic
// from two separate IP addrs.
func CreateOrConnectNetwork(t *hivesim.T, container, network string) {
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
