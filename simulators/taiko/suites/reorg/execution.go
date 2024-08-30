package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

func (r ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	for i, node := range testnet.Nodes {
		if err := node.Start(); err != nil {
			t.Fatalf("ReorgTestSpec failed to start node[%d], err: %v", i, err)
		}
	}

	r.reorgAndVerifyFirstCluster(ctx, t, testnet.Nodes[0])

	r.VerifyL2Nodes(ctx, t, r.L2TargetNumber, testnet.Nodes)
}

func (r ReorgTestSpec) reorgAndVerifyFirstCluster(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		anvil    = node.AnvilClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
		prover   = node.ProverClient
		l2eth    = node.L2EthClient

		timeout            = time.Second * 60
		l2ReorgStartNumber = r.L2TargetNumber
		reorgDepth         = r.ReorgDepth
	)
	if reorgDepth == 0 {
		// random [5, 100) value
		reorgDepth = rand.Uint64N(100-10) + 10
	}

	if anvil == nil || driver == nil || proposer == nil || prover == nil || l2eth == nil {
		t.Fatalf("anvil, driver, proposer, prover or l2eth client is nil!")
	}

	if !anvil.IsRunning() || !l2eth.IsRunning() {
		t.Fatalf("anvil or l2eth node is not running!")
	}

	var (
		l2Client = l2eth.EthClient()
	)

	if err := l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber); err != nil {
		t.Fatalf("failed to wait %s latest number, err: %v", l2eth.ClientType(), err)
	}

	l2OriginHeader, err := l2Client.HeaderByNumber(ctx, new(big.Int).SetUint64(l2ReorgStartNumber))
	if err != nil {
		t.Fatalf("failed to get %s header by number, err: %v", l2eth.ClientType(), err)
	}

	// reorg
	snapshot, number := anvil.SetReorgPoint()
	if err := anvil.WaitLatestNumber(ctx, timeout, number+reorgDepth); err != nil {
		t.Fatalf("failed to wait %s for number, err: %v", anvil.ClientType(), err)
	}
	anvil.Reorg(snapshot)

	// Wait for the reorg to be processed.
	if err = l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber); err != nil {
		t.Fatalf("failed to wait %s latest number, err: %v", l2eth.ClientType(), err)
	}

	l2ReorgedHeader, err := l2Client.HeaderByNumber(ctx, new(big.Int).SetUint64(l2ReorgStartNumber))
	if err != nil {
		t.Fatalf("failed to get %s header by number, err: %v", l2eth.ClientType(), err)
	}

	// Verify the reorged header hash is equal to the origin header hash.
	if l2OriginHeader.Hash() != l2ReorgedHeader.Hash() {
		t.Fatalf("%s header hash %s is not equal to %s reorged header hash %s", l2eth.ClientType(), l2eth.ClientType(), l2OriginHeader.Hash().Hex(), l2ReorgedHeader.Hash().Hex())
	}
}
