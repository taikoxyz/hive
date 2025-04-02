package preconf

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, &ReorgProposeTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name: "reorg_propose",
			},
		},
	})
}

type ReorgProposeTestSpec struct {
	PreconfTestSpec
}

func (r *ReorgProposeTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_reorg_propose"

	return cfg
}

func (r *ReorgProposeTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
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

	l2geth.WaitLatestNumber(ctx, time.Minute*3, proposer.PacayaClients.ForkHeight-1)

	// stop the proposer.
	proposer.Shutdown()

	// For DevDebug
	if testnet.DevDebug {
		driver.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	// Test reorg propose blocks.
	r.reorgProposeBlocks(ctx, t, testnet)
}

func (r *ReorgProposeTestSpec) reorgProposeBlocks(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
	var (
		node   = testnet.Nodes[0]
		anvil  = node.AnvilClient
		l2Geth = node.L2EthClient
		driver = node.DriverClient
	)

	// Reorg propose blocks.
	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		// Start record reorg points.
		anvil.StartRecordReorgPoints(ctx, l2Geth.EthClient)

		l2Header, anchorL1Header, err := preconferBlock(params.ChainAuths[index*2+1].PrivateKey, driver.Client, driver.PreconfServerURL(), 5, nil)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// Verify latest preconf block.
		verifyL2Chain(t, true, index, testnet.Nodes, l2Header)

		// propose txs.
		_, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
		t.FailIfNotNil(err, "cannot propose txs")

		// Verify latest propose block.
		verifyL2Chain(t, false, index, testnet.Nodes, l2Header)

		// Reorg to the specified l2 block.
		anvil.Reorg(l2Header.Number.Uint64()-1, nil)

		// Verify latest propose block.
		verifyL2Chain(t, false, index, testnet.Nodes, l2Header)
	}
}
