package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	tn "taiko/common/testnet"
	"time"
)

func (r ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	if len(testnet.Nodes) != 1 {
		t.Fatalf("testnet nodes count is not 1, got: %v", len(testnet.Nodes))
	}
	var (
		node     = testnet.Nodes[0]
		anvil    = node.AnvilClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
		prover   = node.ProverClient
		l2eth    = node.L2EthClient
		l2Number = big.NewInt(10)
	)
	if anvil == nil || driver == nil || proposer == nil || prover == nil || l2eth == nil {
		t.Fatalf("anvil, driver, proposer, prover or l2eth client is nil!")
	}

	if !anvil.IsRunning() || !l2eth.IsRunning() {
		t.Fatalf("anvil or l2eth node is not running!")
	}

	snapshots := make(map[uint64]string)
	stopCh := make(chan struct{})
	go func() {
		number := uint64(1)
		for {
			select {
			// Collected enough snapshots, stop and return.
			case <-stopCh:
				return
			default:
				err := anvil.WaitNumber(ctx, time.Second*100, number)
				if err != nil {
					t.Fatalf("failed to verify anvil number, err: %v", err)
				}
				snapshots[number] = anvil.SetL1Snapshot()
				number++
			}
		}
	}()

	if err := l2eth.WaitNumber(ctx, time.Second*100, l2Number.Uint64()); err != nil {
		t.Fatalf("failed to verify l2eth number, err: %v", err)
	}
	close(stopCh)

	l2Client := l2eth.EthClient()
	originHeader, err := l2Client.HeaderByNumber(ctx, l2Number)
	if err != nil {
		t.Fatalf("failed to get header by number, err: %v", err)
	}

	l1Origin, err := l2eth.L1OriginByID(ctx, l2Number)
	if err != nil {
		t.Fatalf("failed to get l1 origin by id, err: %v", err)
	}

	// Reorg l1geth.
	snapshot := snapshots[l1Origin.L1BlockHeight.Uint64()-1]
	anvil.RevertL1Snapshot(snapshot)

	if err := l2eth.WaitNumber(ctx, time.Second*100, l2Number.Uint64()); err != nil {
		t.Fatalf("failed to verify l2eth number, err: %v", err)
	}

	reorgedHeader, err := l2Client.HeaderByNumber(ctx, l2Number)
	if err != nil {
		t.Fatalf("failed to get l2eth header by number, err: %v", err)
	}

	if originHeader.Hash() != reorgedHeader.Hash() {
		t.Fatalf("reorg failed, origin header hash: %v, reorged header hash: %v", originHeader.Hash(), reorgedHeader.Hash())
	}
}
