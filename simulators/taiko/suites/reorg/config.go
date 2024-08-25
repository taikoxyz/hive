package suite_reorg

import (
	"taiko/common/testnet"
	suite_base "taiko/suites/base"
)

type ReorgTestSpec struct {
	L2ReorgStartNumber uint64
	ReorgDepth         uint64
	suite_base.BaseTestSpec
}

func (r ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_reorg_test"
	return cfg
}
