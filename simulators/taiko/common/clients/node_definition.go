package clients

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

// Describe a node setup, which consists of:
// - Execution Client
// - Beacon Client
// - Validator Client
type NodeDefinition struct {
	ConsensusClient string `json:"consensus_client"`
	ValidatorClient string `json:"validator_client"`

	// L2chain client Types
	L2EthClient    string `json:"taiko_geth_client"`
	ProposerClient string `json:"proposer_client"`
	ProverClient   string `json:"prover_client"`
	DriverClient   string `json:"driver_client"`

	// Execution Config
	ExecutionClientTTD *big.Int       `json:"execution_client_ttd,omitempty"`
	Chain              []*types.Block `json:"chain,omitempty"`

	// Beacon Config
	BeaconNodeTTD *big.Int `json:"beacon_node_ttd,omitempty"`

	// Node Config
	DisableStartup bool `json:"disable_startup"`

	// Subnet Configuration
	ExecutionSubnet string `json:"execution_subnet"`
	ConsensusSubnet string `json:"consensus_subnet"`
	Subnet          string `json:"subnet"`
}

func (n *NodeDefinition) ValidatorClientName() string {
	if n.ValidatorClient == "" {
		return beaconNodeToValidator(n.ConsensusClient)
	}
	return n.ValidatorClient
}

func (n *NodeDefinition) GetExecutionSubnet() string {
	if n.ExecutionSubnet != "" {
		return n.ExecutionSubnet
	}
	if n.Subnet != "" {
		return n.Subnet
	}
	return ""
}

func (n *NodeDefinition) GetConsensusSubnet() string {
	if n.ConsensusSubnet != "" {
		return n.ConsensusSubnet
	}
	if n.Subnet != "" {
		return n.Subnet
	}
	return ""
}

func beaconNodeToValidator(name string) string {
	name, branch, hasBranch := strings.Cut(name, "_")
	name = strings.TrimSuffix(name, "-bn")
	validator := name + "-vc"
	if hasBranch {
		validator += "_" + branch
	}
	return validator
}
