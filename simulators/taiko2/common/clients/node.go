package clients

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/marioevz/eth-clients/clients"
	"github.com/marioevz/eth-clients/clients/execution"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"strings"
	"taiko2/common/utils"
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
	// Clients that comprise the node
	ExecutionClient *ExecutionClient
	BeaconClient    *BeaconClient
	ValidatorClient *ValidatorClient
	// Whether this node can be used to query test verification information
	Verification bool
}

func (n *Node) Logf(format string, values ...interface{}) {
	if l := n.Logging; l != nil {
		l.Logf(format, values...)
	}
}

// Starts all clients included in the bundle
func (n *Node) Start() error {
	n.Logf("Starting validator client bundle %d", n.Index)
	if n.ExecutionClient != nil {
		if err := n.ExecutionClient.Start(); err != nil {
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
	return nil
}

var networkCreated = make(map[hivesim.SuiteID]bool)

func (n *Node) CreateOrConnectNetwork(t *hivesim.T, network string) error {
	n.Logf("Creating or connecting to network %s", network)
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true

	}
	if n.ExecutionClient != nil {
		client := n.ExecutionClient.HiveClient()
		if err := t.Sim.ConnectContainer(t.SuiteID, network, client.Container); err != nil {
			t.Fatalf("can't connect %s container to network, err = %v", client.Type, err)
		}
	}
	if n.BeaconClient != nil {
		client := n.BeaconClient.HiveClient()
		if err := t.Sim.ConnectContainer(t.SuiteID, network, client.Container); err != nil {
			t.Fatalf("can't connect %s container to network, err = %v", client.Type, err)
		}
	}
	if n.ValidatorClient != nil {
		client := n.ValidatorClient.HiveClient()
		if err := t.Sim.ConnectContainer(t.SuiteID, network, client.Container); err != nil {
			t.Fatalf("can't connect %s container to network, err = %v", client.Type, err)
		}
	}
	return nil
}

func (n *Node) Shutdown() error {
	if err := n.ExecutionClient.Shutdown(); err != nil {
		return err
	}
	if err := n.BeaconClient.Shutdown(); err != nil {
		return err
	}
	if err := n.ValidatorClient.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (n *Node) ClientNames() string {
	var name string
	if n.ExecutionClient != nil {
		name = n.ExecutionClient.ClientType()
	}
	if n.BeaconClient != nil {
		name = fmt.Sprintf("%s/%s", name, n.BeaconClient.ClientName())
	}
	return name
}

func (n *Node) IsRunning() bool {
	return n.ExecutionClient.IsRunning() && n.BeaconClient.IsRunning() && n.ValidatorClient.IsRunning()
}

// Node cluster operations
type Nodes []*Node

// Return all execution clients, even the ones not currently running
func (all Nodes) ExecutionClients() ExecutionClients {
	en := make(ExecutionClients, 0)
	for _, n := range all {
		if n.ExecutionClient != nil {
			en = append(en, n.ExecutionClient)
		}
	}
	return en
}

// Return all proxy pointers, even the ones not currently running
func (all Nodes) Proxies() execution.Proxies {
	ps := make(execution.Proxies, 0)
	for _, n := range all {
		if n.ExecutionClient != nil {
			ps = append(ps, n.ExecutionClient)
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

// Return subset of nodes which are marked as verification nodes
func (all Nodes) VerificationNodes() Nodes {
	// If none is set as verification, then all are verification nodes
	var any bool
	for _, n := range all {
		if n.Verification {
			any = true
			break
		}
	}
	if !any {
		return all
	}

	res := make(Nodes, 0)
	for _, n := range all {
		if n.Verification {
			res = append(res, n)
		}
	}
	return res
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
			if strings.Contains(n.ExecutionClient.ClientType(), filter) {
				ret = append(ret, n)
				break
			}
		}
	}
	return ret
}

func (all Nodes) RemoveNodeAsVerifier(id int) error {
	if id >= len(all) {
		return fmt.Errorf("node %d does not exist", id)
	}
	var any bool
	for _, n := range all {
		if n.Verification {
			any = true
			break
		}
	}
	if any {
		all[id].Verification = false
	} else {
		// If no node is set as verifier, we will set all other nodes as verifiers then
		for i := range all {
			all[i].Verification = (i != id)
		}
	}
	return nil
}

func (all Nodes) ByValidatorIndex(validatorIndex common.ValidatorIndex) *Node {
	for _, n := range all {
		if n.ValidatorClient.ContainsValidatorIndex(validatorIndex) {
			return n
		}
	}
	return nil
}
