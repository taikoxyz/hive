package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		ReorgFinalizeTestSpec{
			ReorgTestSpec{BaseTestSpec: suite_base.BaseTestSpec{Name: "reorg_finalize"}},
		},
	)
}

type ReorgFinalizeTestSpec struct {
	ReorgTestSpec
}

func (r ReorgFinalizeTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node    = testnet.Nodes[0]
		anvil   = node.AnvilClient
		l2eth   = node.L2EthClient
		prover  = node.ProverClient
		propoer = node.ProposerClient
	)

	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if testnet.DevDebug {
		node.ProposerClient.Shutdown()
		node.ProverClient.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	r.reorgStart = 13
	r.reorgDeep = 5
	r.params = &clients.ReorgParams{
		DelayTime:   1,
		DelayNumber: 0,
	}

	// Start recording l1 chain.
	anvil.StartRecordReorgPoints(ctx, l2eth.EthClient)

	// Verify blocks.
	for range time.Tick(time.Second) {
		prover.VerifyBlocks(params.L1Auths[0])
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID >= propoer.PacayaClients.ForkHeight-1 {
			node.ProverClient.Shutdown()
			break
		}
	}

	// Set reorg point before the fork height.
	r.reorgStart = prover.GetLastVerifiedBlockId(ctx) / 2

	r.reorg(ctx, t, node)

	// verify finalize
	prover.WaitLatestVerifiedNumber(ctx, time.Second*5, r.reorgStart)
}
