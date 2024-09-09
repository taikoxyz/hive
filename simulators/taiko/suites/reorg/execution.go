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
	t.Logf("%s: start reorgAndVerifyFirstCluster, target number: %d, reorg depth: %d", r.Name, l2ReorgStartNumber, reorgDepth)

	if anvil == nil || driver == nil || proposer == nil || /*prover == nil ||*/ l2eth == nil {
		t.Fatalf("anvil, driver, proposer, prover or l2eth client is nil!")
	}

	if !anvil.IsRunning() || !l2eth.IsRunning() {
		t.Fatalf("anvil or l2eth node is not running!")
	}

	if err := l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber); err != nil {
		t.Fatalf("%s: failed to wait latest number, err: %v", l2eth.ClientType(), err)
	}

	// Wait until taiko-geth touch verified number.
	l1Number, l2LatestVerifiedNumber, err := anvil.WaitVerifiedNumber(ctx, timeout)
	if err != nil {
		t.Fatalf("failed to wait %s for verified number, err: %v", anvil.ClientType(), err)
	}

	if err := l2eth.WaitLatestNumber(ctx, timeout, l2LatestVerifiedNumber+reorgDepth); err != nil {
		t.Fatalf("%s: failed to wait latest number, err: %v", l2eth.ClientType(), err)
	}

	l2OriginHeader, err := l2eth.EthClient().HeaderByNumber(ctx, new(big.Int).SetUint64(l2LatestVerifiedNumber))
	if err != nil {
		t.Fatalf("failed to get %s header by number, err: %v", l2eth.ClientType(), err)
	}

	// pause driver, proposer, prover
	driver.PauseClient()
	proposer.PauseClient()
	prover.PauseClient()

	// Reorg l1 eth chain.
	anvil.Reorg(l1Number)

	// unpause driver, proposer, prover
	driver.UnpauseClient()
	proposer.UnpauseClient()
	prover.UnpauseClient()

	// Wait for the reorg to be processed.
	if err = l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber+reorgDepth); err != nil {
		t.Fatalf("failed to wait %s latest number, err: %v", l2eth.ClientType(), err)
	}

	l2ReorgedHeader, err := l2eth.EthClient().HeaderByNumber(ctx, new(big.Int).SetUint64(l2LatestVerifiedNumber))
	if err != nil {
		t.Fatalf("failed to get %s header by number, err: %v", l2eth.ClientType(), err)
	}

	// Verify the reorged header hash is equal to the origin header hash.
	if l2OriginHeader.Hash() != l2ReorgedHeader.Hash() {
		t.Fatalf("%s header hash %s is not equal to %s reorged header hash %s", l2eth.ClientType(), l2eth.ClientType(), l2OriginHeader.Hash().Hex(), l2ReorgedHeader.Hash().Hex())
	}
}
