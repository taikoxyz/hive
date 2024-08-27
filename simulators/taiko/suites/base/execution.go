package suite_base

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"sync"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

func (ts BaseTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	wg := sync.WaitGroup{}
	wg.Add(len(testnet.Nodes))
	target := uint64(5)
	for _, node := range testnet.Nodes {
		node := node
		go func() {
			defer wg.Done()
			ts.verify(ctx, t, node, target)
		}()
	}
	wg.Wait()

	ts.VerifyL2Nodes(ctx, t, target, testnet.Nodes)
}

func (ts BaseTestSpec) verify(ctx context.Context, t *hivesim.T, node *clients.Node, targetNumber uint64) {
	var (
		anvil = node.AnvilClient
		l1Eth = node.L1EthClient
		l2Eth = node.L2EthClient

		timeout = time.Second * 60
	)

	if anvil != nil {
		if !anvil.IsRunning() {
			t.Fatalf("anvil node is not running!")
		}
		// Verify l1eth node run successfully.
		err := anvil.WaitLatestNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
		t.Logf("%s node is running successfully", anvil.ClientType())
	}
	if l1Eth != nil {
		if !l1Eth.IsRunning() {
			t.Fatalf("l1eth node is not running!")
		}
		// Verify l1eth node run successfully.
		err := l1Eth.WaitLatestNumber(ctx, timeout, targetNumber)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
		t.Logf("%s node is running successfully", l1Eth.ClientType())
	}

	if l2Eth != nil {
		if !l2Eth.IsRunning() {
			t.Fatalf("l2eth node is not running!")
		}
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
		if node.L2EthClient == nil {
			continue
		}
		client := node.L2EthClient.EthClient()
		hd, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(target))
		if err != nil {
			t.Fatalf("failed to get header from [%d]:%s, err: %v", i, node.L1EthClient.ClientType(), err)
		}
		if header == nil {
			header, index = hd, i
		} else if header.Hash() != hd.Hash() {
			t.Fatalf("the %d number of %s's hash are different, [%d]:%s != [%d]:%s",
				target,
				node.L2EthClient.ClientType(),
				index, header.Hash().String(),
				i, hd.Hash().String(),
			)
		}
	}
}
