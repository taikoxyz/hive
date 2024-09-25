package suite_reorg

import (
	"taiko/common/testnet"
	"taiko/suites"
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

func (r ReorgTestSpec) DebugTestSpec() suites.TestSpec {
	return ReorgTestSpec{
		r.BaseTestSpec.DebugTestSpec().(suite_base.BaseTestSpec),
	}
}
