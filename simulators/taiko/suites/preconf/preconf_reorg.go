package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"taiko/common/testnet"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, &ReorgTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name: "reorg",
			},
		},
	})
}

type ReorgTestSpec struct {
	PreconfTestSpec
}

func (r *ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_preconf_reorg"

	return cfg
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

	// Reorg propose blocks.
	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		// Start record reorg points.
		anvil.StartRecordReorgPoints(ctx, l2Geth.EthClient)

		l2Header, anchorL1Header, _, err := preconferBlock(index, driver.Client, driver.PreconfServerURL(), 5)
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

func (r *PreconfTestSpec) reorgPreconfBlocks(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
	var (
		node   = testnet.Nodes[0]
		driver = node.DriverClient
	)

	// Reorg propose blocks.

	l2Number, err := driver.L2.BlockNumber(ctx)
	t.FailIfNotNil(err, "cannot get l2 header by number")

	l2Header, anchorL1Header, _, err := preconferBlock(0, driver.Client, driver.PreconfServerURL(), 5)
	t.FailIfNotNil(err, "cannot preconfirmer proposer")

	// Verify latest preconf block.
	verifyL2Chain(t, true, 0, testnet.Nodes, l2Header)

	preconfBlocks := make([]*types.Block, 0)
	for number := l2Number + 1; number <= l2Header.Number.Uint64(); number++ {
		block, err := driver.L2.BlockByNumber(ctx, big.NewInt(int64(number)))
		t.FailIfNotNil(err, fmt.Sprintf("cannot get preconf block by number %d", number))
		preconfBlocks = append(preconfBlocks, block)
	}

	// change the anchorL1Header time to reorg the preconf blocks.
	anchorL1Header.Time += 1

	// propose txs.
	_, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
	t.FailIfNotNil(err, "cannot propose txs")

	proposeBlocks := make([]*types.Block, 0)
	for number := l2Number + 1; number <= l2Header.Number.Uint64(); number++ {
		block, err := driver.L2.BlockByNumber(ctx, big.NewInt(int64(number)))
		t.FailIfNotNil(err, fmt.Sprintf("cannot get propose block by number: %d", number))
		proposeBlocks = append(proposeBlocks, block)
	}

	for i := 0; i < len(proposeBlocks); i++ {
		t.NotEqual(proposeBlocks[i].Hash().String(), preconfBlocks[i].Hash().String())
		t.Equal(proposeBlocks[i].TxHash().String(), preconfBlocks[i].TxHash().String())
	}
}
