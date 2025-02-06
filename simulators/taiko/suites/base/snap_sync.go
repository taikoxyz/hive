package suite_base

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"strings"
	tn "taiko/common/testnet"
	"taiko/params"
)

func (ts BaseTestSpec) snapSyncVerify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet, target uint64) {
	// Verify snap sync.
	for _, node := range testnet.Nodes[1:] {
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
		}

		if target < 65 {
			continue
		}

		l2cli := node.L2EthClient.EthClient
		for num := target - 65; num <= target; num++ {
			_, err := l2cli.BalanceAt(ctx, common.HexToAddress(params.ParamByKey("TAIKO_INBOX")), new(big.Int).SetUint64(num))
			if num == target-65 && (err == nil || !strings.Contains(err.Error(), "missing trie node")) {
				t.Fatalf("balance %d should be missing trie node, err: %v", num, err)
			}
			if num > target-65 && err != nil {
				t.Fatalf("should not appear error when get balance numer: %d, err: %v", num, err)
			}
		}
	}
}
