package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/marioevz/eth-clients/clients"
	"github.com/marioevz/eth-clients/clients/execution"
	"strings"
	"taiko/common/utils"
)

const (
	EthHttpPort  = 8545
	EthWSPort    = 8546
	EthEngineRPC = 8551
	BeaconPort   = 3500
)

type Client interface {
	clients.Client
	HiveClient() *hivesim.Client
}

// A node bundles together:
// - Running Execution client
// - Running Beacon client
// - Running Validator client
// Contains a flag that marks a node that can be used to query
// test verification information.
type Node struct {
	// Logging interface for all the events that happen in the node
	Logging utils.Logging
	// Index of the node in the network/testnet
	Index int
	// L1 chain clients that comprise the node
	AnvilClient     *AnvilClient
	L1EthClient     *ExecutionClient
	BeaconClient    *BeaconClient
	ValidatorClient *ValidatorClient

	// L2 chain clients that comprise the node
	L2EthClient    *TaikoGethClient
	DriverClient   *DriverClient
	ProposerClient *ProposerClient
	ProverClient   *ProverClient
}

func (n *Node) Logf(format string, values ...interface{}) {
	if l := n.Logging; l != nil {
		l.Logf(format, values...)
	}
}

// Starts all clients included in the bundle
func (n *Node) Start() error {
	n.Logf("Starting validator client bundle %d", n.Index)
	if n.L1EthClient != nil {
		if err := n.L1EthClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No execution client started")
	}
	if n.BeaconClient != nil {
		if err := n.BeaconClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No beacon client started")
	}
	if n.ValidatorClient != nil {
		if err := n.ValidatorClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No validator client started")
	}

	if n.L2EthClient != nil {
		if err := n.L2EthClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No taiko geth client started")
	}

	// Deploy contracts if needed
	if n.L1EthClient != nil && (n.DriverClient != nil || n.ProposerClient != nil || n.ProverClient != nil) {
		if err := n.L1EthClient.DeployContracts(context.Background()); err != nil {
			return err
		}
	}

	if n.DriverClient != nil {
		if err := n.DriverClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No driver client started")
	}
	if n.ProposerClient != nil {
		if err := n.ProposerClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No proposer client started")
	}
	if n.ProverClient != nil {
		if err := n.ProverClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No prover client started")
	}
	return nil
}

var networkCreated = make(map[hivesim.SuiteID]bool)

func (n *Node) CreateNetwork(t *hivesim.T, network string) error {
	n.Logf("Creating network %s", network)
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true
	}
	return nil
}

func (n *Node) Shutdown() error {
	if n.ProverClient != nil {
		if err := n.ProverClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.ProposerClient != nil {
		if err := n.ProposerClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.DriverClient != nil {
		if err := n.DriverClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.L2EthClient != nil {
		if err := n.L2EthClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.ValidatorClient != nil {
		if err := n.ValidatorClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.BeaconClient != nil {
		if err := n.BeaconClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.L1EthClient != nil {
		if err := n.L1EthClient.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) ClientNames() string {
	var name string
	if n.L1EthClient != nil {
		name = n.L1EthClient.ClientType()
	}
	if n.BeaconClient != nil {
		name = fmt.Sprintf("%s/%s", name, n.BeaconClient.ClientName())
	}
	return name
}

func (n *Node) IsRunning() bool {
	return n.L1EthClient.IsRunning() && n.BeaconClient.IsRunning() && n.ValidatorClient.IsRunning()
}

// Node cluster operations
type Nodes []*Node

// Return all execution clients, even the ones not currently running
func (all Nodes) ExecutionClients() ExecutionClients {
	en := make(ExecutionClients, 0)
	for _, n := range all {
		if n.L1EthClient != nil {
			en = append(en, n.L1EthClient)
		}
	}
	return en
}

// Return all proxy pointers, even the ones not currently running
func (all Nodes) Proxies() execution.Proxies {
	ps := make(execution.Proxies, 0)
	for _, n := range all {
		if n.L1EthClient != nil {
			ps = append(ps, n.L1EthClient)
		}
	}
	return ps
}

// Return all beacon clients, even the ones not currently running
func (all Nodes) BeaconClients() BeaconClients {
	bn := make(BeaconClients, 0)
	for _, n := range all {
		if n.BeaconClient != nil {
			bn = append(bn, n.BeaconClient)
		}
	}
	return bn
}

// Return all validator clients, even the ones not currently running
func (all Nodes) ValidatorClients() ValidatorClients {
	vc := make(ValidatorClients, 0)
	for _, n := range all {
		if n.ValidatorClient != nil {
			vc = append(vc, n.ValidatorClient)
		}
	}
	return vc
}

// Return subset of nodes that are currently running
func (all Nodes) Running() Nodes {
	res := make(Nodes, 0)
	for _, n := range all {
		if n.IsRunning() {
			res = append(res, n)
		}
	}
	return res
}

func (all Nodes) FilterByCL(filters []string) Nodes {
	ret := make(Nodes, 0)
	for _, n := range all {
		for _, filter := range filters {
			if strings.Contains(n.BeaconClient.ClientName(), filter) {
				ret = append(ret, n)
				break
			}
		}
	}
	return ret
}

func (all Nodes) FilterByEL(filters []string) Nodes {
	ret := make(Nodes, 0)
	for _, n := range all {
		for _, filter := range filters {
			if strings.Contains(n.L1EthClient.ClientType(), filter) {
				ret = append(ret, n)
				break
			}
		}
	}
	return ret
}
