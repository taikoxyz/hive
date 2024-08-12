package consensus_config

import (
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

type Spec struct {
	//common.Phase0Preset      `json:",inline" yaml:",inline"`
	//common.AltairPreset      `json:",inline" yaml:",inline"`
	//common.BellatrixPreset   `json:",inline" yaml:",inline"`
	//common.CapellaPreset     `json:",inline" yaml:",inline"`
	//common.DenebPreset       `json:",inline" yaml:",inline"`
	params.BeaconChainConfig `json:",inline" yaml:",inline"`

	// Experimental, for bellatrix
	common.ExecutionEngine `json:"-" yaml:"-"`
}

func (spec *Spec) SlotToEpoch(s primitives.Slot) primitives.Epoch {
	return primitives.Epoch(s / spec.SlotsPerEpoch)
}
