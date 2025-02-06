package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

func (ts BaseTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		nodes  = testnet.Nodes
		target = ts.L2TargetNumber
		anvil  = nodes[0].AnvilClient
	)
	if target == 0 {
		if ts.L2SyncMode == "full" {
			target = rand.Uint64N(50-10) + 10
		} else {
			target = rand.Uint64N(130-65) + 65
		}
	}
	t.Logf("BaseTestSpec target number: %d, sync module: %s", target, ts.L2SyncMode)

	if err := nodes[0].Start(); err != nil {
		t.Fatalf("BaseTestSpec failed to start 0 node: %v", err)
	}
	for i, node := range nodes[1:] {
		if err := node.L2EthClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	// For Debug
	if ts.IsDebug() {
		time.Sleep(time.Minute * 120)
	}

	waitL2LatestNumber(ctx, t, target, nodes[0].L2EthClient)

	// Start the other cluster's l2eth and driver nodes.
	for i, node := range nodes[1:] {
		if err := node.DriverClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	// Verify the synced blocks.
	if ts.L2SyncMode == "full" {
		ts.fullSyncVerify(ctx, t, testnet, target)
	} else {
		ts.snapSyncVerify(ctx, t, testnet, target)
	}

	// Verify all l2eth nodes.
	ts.verifyL2Nodes(ctx, t, anvil.GetLastVerifiedBlockId(ctx), nodes)
}

func waitL2LatestNumber(ctx context.Context, t *hivesim.T, targetNumber uint64, l2Eth *clients.TaikoGethClient) {
	// Verify l2eth node run successfully.
	err := l2Eth.WaitLatestNumber(ctx, time.Second*60, targetNumber)
	if err != nil {
		t.Fatalf("failed to get latest l2geth number, err: %v", err)
	}
}

func waitLatestVerifiedNumber(ctx context.Context, t *hivesim.T, targetNumber uint64, anvil *clients.AnvilClient) {
	// Verify l2eth node run successfully.
	err := anvil.WaitLatestVerifiedNumber(ctx, time.Second*60, targetNumber)
	if err != nil {
		t.Fatalf("failed to get latest l2geth number, err: %v", err)
	}
}

func (ts BaseTestSpec) verifyL2Nodes(ctx context.Context, t *hivesim.T, latestVerified uint64, nodes []*clients.Node) {
	var (
		l2cli = nodes[0].L2EthClient.EthClient
		index int
	)

	waitLatestVerifiedNumber(ctx, t, latestVerified+1, nodes[0].AnvilClient)

	header, err := l2cli.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatalf("failed to get l2 block number: %v", err)
	}
	target := header.Number.Uint64()

	for i, node := range nodes[1:] {
		l2client := node.L2EthClient

		client := l2client.EthClient
		hd, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(target))
		if err != nil {
			t.Fatalf("failed to get header from [%d]:%s, err: %v", i, node.L1EthClient.ClientType(), err)
		}
		if header.Hash() != hd.Hash() {
			t.Fatalf("the %d number of %s's hash are different, [%d]:%s != [%d]:%s",
				target,
				l2client.ClientType(),
				index, header.Hash().String(),
				i, hd.Hash().String(),
			)
		}
	}
}
