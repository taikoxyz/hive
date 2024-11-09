package consensus_config

import (
	"errors"
	"github.com/ethereum/hive/hivesim"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"math/big"
	"taiko/common/config"
)

type ConsensusConfig struct {
	ValidatorCount                  *big.Int `json:"validator_count,omitempty"`
	KeyTranches                     *big.Int `json:"key_tranches,omitempty"`
	SlotsPerEpoch                   *big.Int `json:"slots_per_epoch,omitempty"`
	SlotTime                        *big.Int `json:"slot_time,omitempty"`
	SafeSlotsToImportOptimistically *big.Int `json:"safe_slots_to_import_optimistically,omitempty"`
	ExtraShares                     *big.Int `json:"extra_shares,omitempty"`
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
func (a *ConsensusConfig) Join(b *ConsensusConfig) *ConsensusConfig {
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}
	c := ConsensusConfig{}
	c.ValidatorCount = choose(a.ValidatorCount, b.ValidatorCount)
	c.KeyTranches = choose(a.KeyTranches, b.KeyTranches)
	c.SlotsPerEpoch = choose(a.SlotsPerEpoch, b.SlotsPerEpoch)
	c.SlotTime = choose(a.SlotTime, b.SlotTime)
	c.SafeSlotsToImportOptimistically = choose(
		a.SafeSlotsToImportOptimistically,
		b.SafeSlotsToImportOptimistically,
	)
	c.ExtraShares = choose(a.ExtraShares, b.ExtraShares)
	return &c
}

func LoadChainConfig(path string) (*params.BeaconChainConfig, error) {
	return params.UnmarshalConfigFile(path, nil)
}

func StateBundle(state state.BeaconState) (hivesim.StartOption, error) {
	type MinimumSSZMarshal interface {
		MarshalSSZ() ([]byte, error)
	}
	marshalFn := func(o interface{}) ([]byte, error) {
		marshaler, ok := o.(MinimumSSZMarshal)
		if !ok {
			return nil, errors.New("not a marshaler")
		}
		return marshaler.MarshalSSZ()
	}
	encoded, err := marshalFn(state)
	if err != nil {
		return nil, err
	}
	return hivesim.WithDynamicFile(
		"/hive/input/genesis.ssz",
		config.BytesSource(encoded),
	), nil
}

func BuildSpec(
	configFile string,
) (*Spec, error) {
	cfg, err := params.UnmarshalConfigFile(configFile, nil)
	if err != nil {
		return nil, err
	}

	return &Spec{
		BeaconChainConfig: *cfg,
	}, nil
}

func ConfigBundle(content []byte) (hivesim.StartOption, error) {
	return hivesim.Bundle(
		hivesim.WithDynamicFile(
			"/hive/input/config.yaml",
			config.BytesSource(content),
		),
	), nil
}
