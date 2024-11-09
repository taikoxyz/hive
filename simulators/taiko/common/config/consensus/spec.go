package consensus_config

import (
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
)

type Spec struct {
	//common.Phase0Preset      `json:",inline" yaml:",inline"`
	//common.AltairPreset      `json:",inline" yaml:",inline"`
	//common.BellatrixPreset   `json:",inline" yaml:",inline"`
	//common.CapellaPreset     `json:",inline" yaml:",inline"`
	//common.DenebPreset       `json:",inline" yaml:",inline"`
	params.BeaconChainConfig `json:",inline" yaml:",inline"`
}

func (spec *Spec) SlotToEpoch(s primitives.Slot) primitives.Epoch {
	return primitives.Epoch(s / spec.SlotsPerEpoch)
}
