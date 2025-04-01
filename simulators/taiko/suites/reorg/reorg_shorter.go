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
		ReorgShorterTestSpec{
			ReorgTestSpec{BaseTestSpec: suite_base.BaseTestSpec{Name: "reorg_shorter"}},
		},
	)
}

type ReorgShorterTestSpec struct {
	ReorgTestSpec
}

func (r ReorgShorterTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	prover := node.ProverClient

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
		DelayNumber: -2,
	}

	for range time.Tick(time.Second) {
		prover.VerifyBlocks(params.L1Auths[0])
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID > 0 {
			node.ProverClient.Shutdown()
			break
		}
	}

	r.reorg(ctx, t, node)
}
