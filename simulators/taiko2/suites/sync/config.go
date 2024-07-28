package suite_sync

import (
	"taiko2/common/clients"
	"taiko2/common/testnet"
	suite_base "taiko2/suites/base"
)

type SyncTestSpec struct {
	suite_base.BaseTestSpec
}

func (ts SyncTestSpec) GetTestnetConfig(
	allNodeDefinitions clients.NodeDefinitions,
) *testnet.Config {
	tc := ts.BaseTestSpec.GetTestnetConfig(allNodeDefinitions)

	// We disable the start of the last node
	tc.NodeDefinitions[len(tc.NodeDefinitions)-1].DisableStartup = true

	return tc
}
