package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/rand/v2"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/common/utils"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		ReorgTestSpec{
			BaseTestSpec: suite_base.BaseTestSpec{
				Name:           "reorg",
				L2TargetNumber: 13,
			},
		},
	)
}

type ReorgTestSpec struct {
	suite_base.BaseTestSpec
}

func (r ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_reorg_test"
	return cfg
}

func (r ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if utils.GetenvBool("HIVE_DEBUG") {
		time.Sleep(time.Minute * 60)
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

	// Start recording reorg points.
	anvil.StartRecordReorgPoints(ctx, l2eth.EthClient)

	for range time.Tick(time.Second) {
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID >= proposer.PacayaClients.ForkHeight {
			break
		}
		t.Nil(prover.VerifyBlocks(params.L1Auths[0]), "failed to verify blocks")
	}

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber)

	if err := prover.WaitLatestVerifiedNumber(ctx, timeout, l2ReorgStartNumber/2); err != nil {
		t.Fatalf("%s: failed to wait LatestVerifiedNumber: %d, err: %v", anvil.ClientType(), l2ReorgStartNumber/2, err)
	}

	// Get reorg point.
	latestVerified := prover.GetLastVerifiedBlockId(ctx)
	t.Logf("get the l2chain latestVerified: %d", latestVerified)

	// pause driver, proposer, prover
	driver.PauseClient()
	proposer.PauseClient()
	prover.PauseClient()

	// Reorg l1 eth chain.
	anvil.Reorg(latestVerified)

	// unpause driver, proposer, prover
	driver.UnpauseClient()
	proposer.UnpauseClient()
	prover.UnpauseClient()

	if err := prover.WaitLatestVerifiedNumber(ctx, timeout, latestVerified+1); err != nil {
		t.Fatalf("%s: failed to wait LatestVerifiedNumber: %d, err: %v", anvil.ClientType(), latestVerified+1, err)
	}

}
