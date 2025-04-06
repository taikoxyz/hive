package preconf

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"golang.org/x/sync/errgroup"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	"taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, &ReorgPreconfTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name: "reorg_preconf",
			},
		},
	})
}

type ReorgPreconfTestSpec struct {
	PreconfTestSpec
}

func (r *ReorgPreconfTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_reorg_preconf"

	return cfg
}

func (r *ReorgPreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
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

	// Test reorg preconf blocks.
	r.reorgPreconfBlocks(ctx, t, testnet)
}

// reorg the first cluster's preconf blocks and verify all the clusters' reorged preconf blocks.
func (r *ReorgPreconfTestSpec) reorgPreconfBlocks(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet) {
	var (
		node   = testnet.Nodes[0]
		anvil  = node.AnvilClient
		l2geth = node.L2EthClient

		l1Number = anvil.BlockNumber(ctx)
		l2Number = l2geth.BlockNumber(ctx)
	)

	var (
		batchSize      = rand.IntN(50-10) + 10
		preconfHeaders []*types.Header

		eg errgroup.Group
	)

	// create a batch of preconf blocks.
	// Reorg propose blocks.
	eg.Go(func() error {
		createPreconfBlocks(ctx, t, node, l1Number, l2Number, batchSize)
		return nil
	})

	eg.Go(func() error {
		preconfHeaders = createPreconfBlocks(ctx, t, node, l1Number.Sub(l1Number, big.NewInt(1)), l2Number, batchSize+5)
		return nil
	})

	t.FailIfNotNil(eg.Wait(), "cannot create a batch of preconf blocks")

	// verify all the clusters' reorged preconf blocks.
	for _, node := range testnet.Nodes {
		l2geth = node.L2EthClient
		for _, header := range preconfHeaders {
			l2geth.HeaderByHash(ctx, header.Hash())
		}
	}
}

func createPreconfBlocks(ctx context.Context, t *hivesim.T, node *clients.Node, l1Number, l2Number *big.Int, batchSize int) (headers []*types.Header) {
	var driver = node.DriverClient

	for i := 0; i < batchSize; i++ {
		requestBody, err := clients.BuildPreconfRequestBody(ctx, driver.Client, params.ChainAuths[driver.Index*2+1].PrivateKey, l1Number, l2Number)
		t.FailIfNotNil(err, "cannot build preconf request body")

		header, err := clients.SendPreconfBlock(driver.PreconfServerURL(), requestBody)
		t.FailIfNotNil(err, "cannot send preconf request")

		headers = append(headers, header)
		time.Sleep(time.Second)
	}

	return
}
