package suite_reorg

import suite_base "taiko/suites/base"

type ReorgTestSpec struct {
	L2ReorgStartNumber uint64
	ReorgDepth         uint64
	suite_base.BaseTestSpec
}
