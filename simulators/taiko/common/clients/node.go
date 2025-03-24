package clients

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/pkg/errors"
	"taiko/common/utils"
	"taiko/params"
)

const (
	EthHttpPort       = 8545
	EthWSPort         = 8546
	EthEngineRPC      = 8551
	GethP2PPort       = 30303
	PreconfServerPort = 7001
	PreconfP2pPort    = 9222
	BeaconPort        = 3500
	BlobscanAPIPort   = 3001
)

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

	BlobScanClient *BlobScanClient
}

func (n *Node) Logf(format string, values ...interface{}) {
	if l := n.Logging; l != nil {
		l.Logf(format, values...)
	}
}

// Starts all clients included in the bundle
func (n *Node) Start() error {
	n.Logf("Starting validator client bundle %d", n.Index)
	if n.AnvilClient != nil {
		if err := n.AnvilClient.Start(); err != nil {
			return errors.Wrap(err, n.AnvilClient.ClientType())
		}
	}
	if n.L1EthClient != nil {
		if err := n.L1EthClient.Start(); err != nil {
			return err
		}
	}
	if n.BeaconClient != nil {
		if err := n.BeaconClient.Start(); err != nil {
			return err
		}
	}
	if n.ValidatorClient != nil {
		if err := n.ValidatorClient.Start(); err != nil {
			return err
		}
	}

	if n.L2EthClient != nil {
		if err := n.L2EthClient.Start(); err != nil {
			return errors.Wrap(err, n.L2EthClient.ClientType())
		}
	}

	// Deploy contracts if needed
	if n.Index == 0 && (n.L1EthClient != nil || n.AnvilClient != nil) {
		l1API := EthExposeAPI(n.AnvilClient)
		if n.L1EthClient != nil {
			l1API = n.L1EthClient
		}

		n.Logf("Deploying contracts in %s node, url: %s\n", l1API.ClientType(), l1API.HttpURL())
		if err := utils.DeployContracts(params.EnvParams(), l1API.HTTPClient()); err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s: failed to deploy contracts", l1API.ClientType()))
		}
		n.Logf("Deployed contracts in %s node, url: %s\n", l1API.ClientType(), l1API.HttpURL())
	}

	// Start mining.
	if n.AnvilClient != nil {
		n.AnvilClient.StartMining()
	}

	if n.BlobScanClient != nil {
		if err := n.BlobScanClient.Start(); err != nil {
			return err
		}
	}

	// Only start the first cluster's driver.
	if n.DriverClient != nil {
		if err := n.DriverClient.Start(); err != nil {
			return errors.Wrap(err, n.DriverClient.ClientType())
		}
	}
	if n.ProposerClient != nil {
		if err := n.ProposerClient.Start(); err != nil {
			return errors.Wrap(err, n.ProposerClient.ClientType())
		}
	}
	if n.ProverClient != nil {
		if err := n.ProverClient.Start(); err != nil {
			return errors.Wrap(err, n.ProverClient.ClientType())
		}
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
	if n.ProverClient != nil && n.ProverClient.IsRunning() {
		if err := n.ProverClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.ProposerClient != nil && n.ProposerClient.IsRunning() {
		if err := n.ProposerClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.DriverClient != nil && n.DriverClient.IsRunning() {
		if err := n.DriverClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.L2EthClient != nil && n.L2EthClient.IsRunning() {
		if err := n.L2EthClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.ValidatorClient != nil && n.ValidatorClient.IsRunning() {
		if err := n.ValidatorClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.BeaconClient != nil && n.BeaconClient.IsRunning() {
		if err := n.BeaconClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.L1EthClient != nil && n.L1EthClient.IsRunning() {
		if err := n.L1EthClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.AnvilClient != nil && n.AnvilClient.IsRunning() {
		if err := n.AnvilClient.Shutdown(); err != nil {
			return err
		}
	}
	if n.BlobScanClient != nil && n.BlobScanClient.IsRunning() {
		if err := n.BlobScanClient.Shutdown(); err != nil {
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
func (all Nodes) L1EthClients() ExecutionClients {
	en := make(ExecutionClients, 0)
	for _, n := range all {
		if n.L1EthClient != nil {
			en = append(en, n.L1EthClient)
		}
	}
	return en
}

func (all Nodes) L2EthClients() L2EthClients {
	en := make(L2EthClients, 0)
	for _, n := range all {
		if n.L2EthClient != nil {
			en = append(en, n.L2EthClient)
		}
	}
	return en
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
