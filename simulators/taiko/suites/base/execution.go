package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

func (ts BaseTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	for _, node := range testnet.Nodes {
		ts.verify(ctx, t, node)
	}
}

func (ts BaseTestSpec) verify(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		anvil = node.AnvilClient
		l1Eth = node.L1EthClient
		l2Eth = node.L2EthClient

		targetNumber = uint64(5)

		timeout = time.Second * 200
	)

	if anvil != nil {
		if !anvil.IsRunning() {
			t.Fatalf("anvil node is not running!")
		}
		// Verify l1eth node run successfully.
		err := anvil.WaitNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
	}
	if l1Eth != nil {
		if !l1Eth.IsRunning() {
			t.Fatalf("l1eth node is not running!")
		}
		// Verify l1eth node run successfully.
		err := l1Eth.WaitNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
	}

	if l2Eth != nil {
		if !l2Eth.IsRunning() {
			t.Fatalf("l2eth node is not running!")
		}
		// Verify l2eth node run successfully.
		err := l2Eth.WaitNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l2geth number, err: %v", err)
		}
	}
}
