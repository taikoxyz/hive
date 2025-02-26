package preconf

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	"time"
)

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
		anvil    = node.AnvilClient
		l2Geth   = node.L2EthClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
		l2geth   = node.L2EthClient
	)

	t.Nil(l2geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight))

	// stop the proposer.
	proposer.PauseClient()

	// For Debug
	if r.Debug {
		driver.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	// Start record reorg points.
	anvil.StartRecordReorgPoints(ctx, l2Geth.EthClient)

	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		l2Header, anchorL1Header, err := preconferBlock(index, driver.Client, driver.PreconfServerURL(), 5)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// Verify latest preconf block.
		verifyL2Chain(t, true, testnet.Nodes, l2Header)

		// propose txs.
		_, _, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
		t.FailIfNotNil(err, "cannot propose txs")

		// Verify latest propose block.
		verifyL2Chain(t, false, testnet.Nodes, l2Header)

		// Reorg to the specified l2 block.
		anvil.Reorg(l2Header.Number.Uint64() - 1)

		// Verify latest propose block.
		verifyL2Chain(t, false, testnet.Nodes, l2Header)
	}
}
