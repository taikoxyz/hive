package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	tn "taiko2/common/testnet"
	"time"
)

var Deneb = "deneb"

func (ts BaseTestSpec) Verify(t *hivesim.T, ctx context.Context, testnet *tn.Testnet) {
	node0 := testnet.Nodes[0]

	if l1EthClient := node0.L1EthClient; l1EthClient != nil {
		err := l1EthClient.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
	}
	if l2EthClient := node0.L2EthClient; l2EthClient != nil {
		err := l2EthClient.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l2geth number, err: %v", err)
		}
	}
}
