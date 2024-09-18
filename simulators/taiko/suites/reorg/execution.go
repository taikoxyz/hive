package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/rand/v2"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

func (r ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	r.reorgAndVerifyFirstCluster(ctx, t, node)
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
	)
	if l2ReorgStartNumber == 0 {
		// random [10, 50) value
		l2ReorgStartNumber = rand.Uint64N(50-10) + 10
	}
	t.Logf("%s: start reorgAndVerifyFirstCluster, target number: %d", r.Name, l2ReorgStartNumber)

	if anvil == nil || driver == nil || proposer == nil || /*prover == nil ||*/ l2eth == nil {
		t.Fatalf("anvil, driver, proposer, prover or l2eth client is nil!")
	}

	if !anvil.IsRunning() || !l2eth.IsRunning() {
		t.Fatalf("anvil or l2eth node is not running!")
	}

	// Start recording reorg points.
	anvil.SetReorgPoint(ctx, l2eth.EthClient())

	if err := l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber); err != nil {
		t.Fatalf("%s: failed to wait latest number, err: %v", l2eth.ClientType(), err)
	}

	if err := anvil.WaitLatestVerifiedNumber(ctx, timeout, l2ReorgStartNumber/2); err != nil {
		t.Fatalf("%s: failed to wait LatestVerifiedNumber: %d, err: %v", anvil.ClientType(), l2ReorgStartNumber/2, err)
	}

	// Get reorg point.
	latestVerified := anvil.GetTaikoDataSlotB(ctx).LastVerifiedBlockId
	t.Logf("%s: latestVerified: %d", l2eth.ClientType(), latestVerified)

	// pause driver, proposer, prover
	driver.PauseClient()
	proposer.PauseClient()
	prover.PauseClient()

	// Reorg l1 eth chain.
	anvil.Reorg(latestVerified + 1)

	// unpause driver, proposer, prover
	driver.UnpauseClient()
	proposer.UnpauseClient()
	prover.UnpauseClient()

	if err := anvil.WaitLatestVerifiedNumber(ctx, timeout, latestVerified+1); err != nil {
		t.Fatalf("%s: failed to wait LatestVerifiedNumber: %d, err: %v", anvil.ClientType(), latestVerified+1, err)
	}

}
