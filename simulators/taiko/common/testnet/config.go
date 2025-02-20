package testnet

import (
	"math/big"
	execution_config "taiko/common/config/execution"
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
	Eth1Consensus execution_config.ExecutionConsensus `json:"eth1_consensus,omitempty"`

	// Execution Layer specific config
	InitialBaseFeePerGas *big.Int `json:"initial_base_fee_per_gas,omitempty"`

	// Consensus Layer specific config
	DisablePeerScoring bool `json:"disable_peer_scoring,omitempty"`

	// Builders
	EnableBuilders bool `json:"enable_builders,omitempty"`
	//BuilderOptions []mock_builder.Option `json:"builder_options,omitempty"`

	// Network
	Network string

	// 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail
	LogLevel int

	FeeReceipt string

	// driver variables, false: beacon sync, true: blob sync
	BeaconSync bool

	// driver variables
	SoftBlockServerPort uint64

	// open debug flag
	Debug bool
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

	if b.Eth1Consensus != nil {
		c.Eth1Consensus = b.Eth1Consensus
	} else {
		c.Eth1Consensus = a.Eth1Consensus
	}

	c.EnableBuilders = b.EnableBuilders || a.EnableBuilders

	return &c
}
