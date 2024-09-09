package suite_base

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
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
	)
	if target == 0 {
		target = rand.Uint64N(50-10) + 10
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
	waitL2LatestNumber(ctx, t, target, nodes[0].L2EthClient)

	// Start the other cluster's l2eth and driver nodes.
	for i, node := range nodes[1:] {
		if err := node.DriverClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	// Verify snap sync.
	for _, node := range nodes[1:] {
		if node.L2EthClient == nil || node.DriverClient == nil {
			continue
		}

		waitL2LatestNumber(ctx, t, target+1, node.L2EthClient)

		l1Origin, err := node.L2EthClient.L1OriginByID(ctx, new(big.Int).SetUint64(target+1))
		if err != nil {
			t.Fatalf("unexpect error when get l1origin, number: %d, err: %v", target+1, err)
		}
		if l1Origin == nil {
			t.Fatalf("l1Origin should not be null, number: %d", target+1)
		}

		for num := target; num > 0; num-- {
			l1Origin, err := node.L2EthClient.L1OriginByID(ctx, new(big.Int).SetUint64(num))
			if err == nil || err.Error() != "not found" {
				t.Fatalf("unexpect error when get l1origin, number: %d, err: %v", num, err)
			}
			if l1Origin != nil {
				t.Fatalf("l1Origin should be null, number: %d", num)
			}
			//time.Sleep(time.Millisecond * 100)
		}
	}

	// Verify all l2eth nodes.
	ts.VerifyL2Nodes(ctx, t, target, nodes)
}

func waitL2LatestNumber(ctx context.Context, t *hivesim.T, targetNumber uint64, l2Eth *clients.TaikoGethClient) {
	var timeout = time.Second * 60

	if l2Eth != nil && l2Eth.IsRunning() {
		// Verify l2eth node run successfully.
		err := l2Eth.WaitLatestNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l2geth number, err: %v", err)
		}
		t.Logf("%s node is running successfully", l2Eth.ClientType())
	}
}

func (ts BaseTestSpec) VerifyL2Nodes(ctx context.Context, t *hivesim.T, target uint64, nodes []*clients.Node) {
	var (
		header *types.Header
		index  int
	)
	for i, node := range nodes {
		l2client := node.L2EthClient
		if l2client == nil || !l2client.IsRunning() {
			t.Logf("L2EthClient[%d] client is empty(%v)", i, l2client == nil)
			continue
		}

		waitL2LatestNumber(ctx, t, target, l2client)

		client := l2client.EthClient()
		hd, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(target))
		if err != nil {
			t.Fatalf("failed to get header from [%d]:%s, err: %v", i, node.L1EthClient.ClientType(), err)
		}
		if header == nil {
			header, index = hd, i
		} else if header.Hash() != hd.Hash() {
			t.Fatalf("the %d number of %s's hash are different, [%d]:%s != [%d]:%s",
				target,
				l2client.ClientType(),
				index, header.Hash().String(),
				i, hd.Hash().String(),
			)
		}
	}
}
