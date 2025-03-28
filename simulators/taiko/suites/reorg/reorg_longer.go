package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/rand/v2"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		ReorgLongerTestSpec{
			ReorgTestSpec{
				BaseTestSpec: suite_base.BaseTestSpec{
					Name:           "reorg_long",
					L2TargetNumber: 13,
				},
			},
		},
	)
}

type ReorgLongerTestSpec struct {
	ReorgTestSpec
}

func (r ReorgLongerTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if testnet.DevDebug {
		node.ProposerClient.Shutdown()
		node.ProverClient.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	r.params = &clients.ReorgParams{
		DelayTime:   0,
		DelayNumber: rand.IntN(2) + 2,
	}

	r.reorg(ctx, t, node)
}
