package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	tn "taiko/common/testnet"
)

func (ts BaseTestSpec) fullSyncVerify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet, target uint64) {
	var (
		nodes = testnet.Nodes
		anvil = nodes[0].AnvilClient
	)
	// Get latest verified number.
	latestVerified := anvil.GetTaikoDataSlotB(ctx).LastVerifiedBlockId
	t.Logf("fullSyncVerify: latestVerified: %d", latestVerified)

	for _, node := range nodes[1:] {
		waitL2LatestNumber(ctx, t, target, node.L2EthClient)

		// verify to all the unverified blocks.
		for num := latestVerified + 1; num <= target; num++ {
			l1Origin, err := node.L2EthClient.L1OriginByID(ctx, new(big.Int).SetUint64(latestVerified+1))
			if err != nil {
				t.Fatalf("unexpect error when get l1origin, number: %d, err: %v", latestVerified+1, err)
			}
			if l1Origin == nil {
				t.Fatalf("l1Origin should not be null, number: %d", latestVerified+1)
			}
		}

		for num := latestVerified; num > 0; num-- {
			l1Origin, err := node.L2EthClient.L1OriginByID(ctx, new(big.Int).SetUint64(num))
			if err == nil || err.Error() != "not found" {
				t.Logf("unexpect error when get l1origin, number: %d, l1Origin is empty: %v, err: %v", num, l1Origin == nil, err)
			}
			if l1Origin != nil {
				t.Fatalf("l1Origin should be null, number: %d", num)
			}
		}
	}
}
