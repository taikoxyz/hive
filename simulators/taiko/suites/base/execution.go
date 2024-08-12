package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

var Deneb = "deneb"

func (ts BaseTestSpec) VerifyNodes(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	for _, node := range testnet.Nodes {
		ts.verify(ctx, t, node)
	}
}

func (ts BaseTestSpec) verify(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		l1Eth = node.L1EthClient
		l2Eth = node.L2EthClient
	)

	if l1Eth != nil {
		if !l1Eth.IsRunning() {
			t.Fatalf("l1eth node is not running!")
		}
		// Verify l1eth node run successfully.
		err := l1Eth.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
	}

	if l2Eth != nil {
		if !l2Eth.IsRunning() {
			t.Fatalf("l2eth node is not running!")
		}
		// Verify l2eth node run successfully.
		err := l2Eth.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l2geth number, err: %v", err)
		}
	}
}
