package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

type PreconfTestSpec struct {
	suite_base.BaseTestSpec

	anchorL1Head *types.Header
}

func (r *PreconfTestSpec) GetTestnetConfig() *testnet.Config {
	// set preconf environment variables.
	params.SetEnvParams("PRECONFIRMATION_SERVER_PORT", fmt.Sprintf("%d", clients.PreconfServerPort))
	params.SetEnvParams("PRECONFIRMATION_SERVER_SIGNATURE_CHECK", "true")
	return r.BaseTestSpec.GetTestnetConfig()
}

func (r *PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	t.Nil(node.Start(), "cannot start node")

	var (
		proposer = node.ProposerClient
		driver   = node.DriverClient
		l2geth   = node.L2EthClient
	)
	if driver == nil || proposer == nil {
		t.Errorf("cannot get driver or proposer client")
		return
	}

	t.Nil(l2geth.WaitLatestNumber(ctx, time.Second*30, 13))
	proposer.PauseClient()
	driver.PauseClient()

	// For Debug
	if r.IsDebug() {
		time.Sleep(time.Minute * 120)
	}
}
