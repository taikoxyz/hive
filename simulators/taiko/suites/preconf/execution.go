package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/utils"

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

	// driver preconf p2p config.
	params.SetEnvParams("PRECONFIRMATION_P2P_DISCOVERY_PATH", "memory")
	params.SetEnvParams("PRECONFIRMATION_P2P_PEERSTORE_PATH", "memory")
	params.SetEnvParams("PRECONFIRMATION_P2P_PRIV_RAW", utils.RandomHash().String())
	params.SetEnvParams("PRECONFIRMATION_P2P_SEQUENCER_KEY", params.ParamByKey("L1_PROPOSER_PRIV_KEY"))
	params.SetEnvParams("PRECONFIRMATION_P2P_NO_DISCOVERY", "true")

	return r.BaseTestSpec.GetTestnetConfig()
}

func (r *PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	// Start all the cluster's nodes.
	for _, node := range testnet.Nodes {
		t.Nil(node.Start(), "cannot start L2EthClient")
	}

	var (
		node     = testnet.Nodes[0]
		driver   = node.DriverClient
		proposer = node.ProposerClient
		l2geth   = node.L2EthClient
	)

	t.Nil(l2geth.WaitLatestNumber(ctx, time.Second*30, 13))

	// stop the proposer.
	proposer.PauseClient()

	// For Debug
	if r.IsDebug() {
		driver.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	for times := 10; times > 0; times-- {
		l2Header, err := preconferProposer(driver.Client, driver.PreconfServerURL(), 5)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// verify all the l2 geth nodes.
		verifyL2Chain(t, testnet.Nodes, l2Header)
	}
}

func verifyL2Chain(t *hivesim.T, nodes []*clients.Node, l2Header *types.Header) {
	for _, node := range nodes {
		l2geth := node.L2EthClient
		t.FailIfNotNil(l2geth.WaitLatestNumber(context.Background(), time.Second*30, l2Header.Number.Uint64()), "cannot get latest number")

		actualHeader, err := l2geth.EthClient.HeaderByNumber(context.Background(), l2Header.Number)
		t.FailIfNotNil(err, "cannot get header by number")
		t.Equal(l2Header.Hash(), actualHeader.Hash(), "header hash mismatch")
	}
}
