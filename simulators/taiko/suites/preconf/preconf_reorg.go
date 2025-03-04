package preconf

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, &ReorgTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name:  "reorg",
				Debug: false,
			},
		},
	})
}

type ReorgTestSpec struct {
	PreconfTestSpec
}

func (r *ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	return r.PreconfTestSpec.GetTestnetConfig()
}

func (r *ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
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

	l2geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight)

	// stop the proposer.
	proposer.PauseClient()

	// For Debug
	if r.Debug {
		driver.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	// Test reorg propose blocks.
	r.reorgProposeBlocks(ctx, t, testnet)

	// Test reorg preconf blocks.
	r.reorgPreconfBlocks(ctx, t, testnet)
}

func (r *ReorgTestSpec) reorgProposeBlocks(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
	var (
		node   = testnet.Nodes[0]
		anvil  = node.AnvilClient
		l2Geth = node.L2EthClient
		driver = node.DriverClient
	)

	// Start record reorg points.
	anvil.StartRecordReorgPoints(ctx, l2Geth.EthClient)
	defer anvil.StopRecordReorgPoints()

	// Reorg propose blocks.
	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		l2Header, anchorL1Header, _, err := preconferBlock(index, driver.Client, driver.PreconfServerURL(), 5)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// Verify latest preconf block.
		verifyL2Chain(t, true, testnet.Nodes, l2Header)

		// propose txs.
		_, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
		t.FailIfNotNil(err, "cannot propose txs")

		// Verify latest propose block.
		verifyL2Chain(t, false, testnet.Nodes, l2Header)

		// Reorg to the specified l2 block.
		anvil.Reorg(l2Header.Number.Uint64() - 1)

		// Verify latest propose block.
		verifyL2Chain(t, false, testnet.Nodes, l2Header)
	}
}

func (r *PreconfTestSpec) reorgPreconfBlocks(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
	var (
		node   = testnet.Nodes[0]
		driver = node.DriverClient
	)

	// Reorg propose blocks.

	l2Header, anchorL1Header, _, err := preconferBlock(0, driver.Client, driver.PreconfServerURL(), 5)
	t.FailIfNotNil(err, "cannot preconfirmer proposer")

	// Verify latest preconf block.
	verifyL2Chain(t, true, testnet.Nodes, l2Header)

	// change the anchorL1Header time to reorg the preconf blocks.
	anchorL1Header.Time += 1

	// propose txs.
	_, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
	t.FailIfNotNil(err, "cannot propose txs")
}

func reorgPreconfBlocksByTime() {}
