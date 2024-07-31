package testnet

import (
	"math/big"
	execution_config "taiko2/common/config/execution"

	blobber_config "github.com/marioevz/blobber/config"
	"taiko2/common/clients"
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)

	MAINNET_SLOT_TIME int64 = 12
	MINIMAL_SLOT_TIME int64 = 6

	// Clients that support the minimal slot time in hive
	MINIMAL_SLOT_TIME_CLIENTS = []string{
		"lighthouse",
		"teku",
		"prysm",
		"lodestar",
		"grandine",
	}
)

type Config struct {
	// Node configurations to launch. Each node as a proportional share of
	// validators.
	NodeDefinitions clients.NodeDefinitions             `json:"node_definitions,omitempty"`
	Eth1Consensus   execution_config.ExecutionConsensus `json:"eth1_consensus,omitempty"`

	// Execution Layer specific config
	InitialBaseFeePerGas *big.Int `json:"initial_base_fee_per_gas,omitempty"`

	// Consensus Layer specific config
	DisablePeerScoring bool `json:"disable_peer_scoring,omitempty"`

	// Builders
	EnableBuilders bool `json:"enable_builders,omitempty"`
	//BuilderOptions []mock_builder.Option `json:"builder_options,omitempty"`

	// Blobber
	EnableBlobber  bool                    `json:"enable_blobber,omitempty"`
	BlobberOptions []blobber_config.Option `json:"blobber_options,omitempty"`

	// Network
	Network string
}

// Choose a configuration value. `b` takes precedence
func choose(a, b *big.Int) *big.Int {
	if b != nil {
		return new(big.Int).Set(b)
	}
	if a != nil {
		return new(big.Int).Set(a)
	}
	return nil
}

// Join two configurations. `b` takes precedence
func (a *Config) Join(b *Config) *Config {
	c := Config{}

	// EL config
	c.InitialBaseFeePerGas = choose(
		a.InitialBaseFeePerGas,
		b.InitialBaseFeePerGas,
	)

	if b.NodeDefinitions != nil {
		c.NodeDefinitions = b.NodeDefinitions
	} else {
		c.NodeDefinitions = a.NodeDefinitions
	}

	if b.Eth1Consensus != nil {
		c.Eth1Consensus = b.Eth1Consensus
	} else {
		c.Eth1Consensus = a.Eth1Consensus
	}

	c.EnableBuilders = b.EnableBuilders || a.EnableBuilders

	return &c
}
