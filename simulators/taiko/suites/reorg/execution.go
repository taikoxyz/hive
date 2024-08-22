package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
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

		timeout  = time.Second * 60
		l2Number = uint64(10)
	)
	if anvil == nil || driver == nil || proposer == nil || prover == nil || l2eth == nil {
		t.Fatalf("anvil, driver, proposer, prover or l2eth client is nil!")
	}

	if !anvil.IsRunning() || !l2eth.IsRunning() {
		t.Fatalf("anvil or l2eth node is not running!")
	}

	if err := l2eth.WaitNumber(ctx, timeout, l2Number); err != nil {
		t.Fatalf("failed to verify l2eth number, err: %v", err)
	}

	// reorg
	snapshot, number := anvil.SetReorgPoint()
	err := anvil.WaitNumber(ctx, timeout, number+5)
	if err != nil {
		t.Fatalf("failed to wait %s for number, err: %v", anvil.ClientType(), err)
	}
	anvil.Reorg(snapshot)

	// TODO
	time.Sleep(time.Second * 300)
}
