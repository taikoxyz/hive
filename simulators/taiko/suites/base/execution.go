package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"taiko/params"
	"time"
)

func (ts BaseTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		nodes    = testnet.Nodes
		target   = ts.L2TargetNumber
		proposer = nodes[0].ProposerClient
		prover   = nodes[0].ProverClient
	)
	if target == 0 {
		target = rand.Uint64N(50-10) + 10
	}
	t.Logf("BaseTestSpec target number: %d", target)

	if err := nodes[0].Start(); err != nil {
		t.Fatalf("BaseTestSpec failed to start 0 node: %v", err)
	}
	for i, node := range nodes[1:] {
		if err := node.L2EthClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	for range time.Tick(time.Second) {
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID >= proposer.PacayaClients.ForkHeight {
			break
		}
		t.Nil(prover.VerifyBlocks(params.L1Auths[0]), "failed to verify blocks")
	}

	// For Debug
	if ts.IsDebug() {
		proposer.PauseClient()
		time.Sleep(time.Minute * 120)
	}

	// Start the other cluster's l2eth and driver nodes.
	for i, node := range nodes[1:] {
		if err := node.DriverClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	// Verify all l2eth nodes.
	ts.verifyL2Nodes(ctx, t, prover.GetLastVerifiedBlockId(ctx), nodes)
}

func waitL2LatestNumber(ctx context.Context, t *hivesim.T, targetNumber uint64, l2Eth *clients.TaikoGethClient) {
	// Verify l2eth node run successfully.
	err := l2Eth.WaitLatestNumber(ctx, time.Second*60, targetNumber)
	if err != nil {
		t.Fatalf("failed to get latest l2geth number, number: %d, err: %v", targetNumber, err)
	}
}

func (ts BaseTestSpec) verifyL2Nodes(ctx context.Context, t *hivesim.T, latestVerified uint64, nodes []*clients.Node) {
	var (
		index  int
		l2cli  = nodes[0].L2EthClient.EthClient
		prover = nodes[0].ProverClient
	)

	// todo: storeForcedInclusion test.

	t.FailIfNotNil(prover.WaitLatestVerifiedNumber(ctx, time.Second*60, latestVerified+1), "failed to wait latest verified number")

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

		for num := uint64(1); num <= target; num++ {
			l1Origin, err := l2client.L1OriginByID(ctx, big.NewInt(0).SetUint64(num))
			if err != nil {
				t.Fatalf("unexpect error when get l1origin, number: %d, err: %v", num, err)
			}
			if l1Origin == nil {
				t.Fatalf("l1Origin should not be null, number: %d", num)
			}
		}
	}
}
