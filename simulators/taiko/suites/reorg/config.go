package suite_reorg

import (
	"taiko/common/testnet"
	suite_base "taiko/suites/base"
)

type ReorgTestSpec struct {
	ReorgDepth uint64
	suite_base.BaseTestSpec
}

func (r ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_reorg_test"
	// reorg test use blob sync.
	cfg.BeaconSync = false

	return cfg
}
