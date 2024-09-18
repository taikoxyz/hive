package suite_reorg

import (
	"taiko/common/testnet"
	suite_base "taiko/suites/base"
)

type ReorgTestSpec struct {
	suite_base.BaseTestSpec
}

func (r ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_reorg_test"
	return cfg
}
