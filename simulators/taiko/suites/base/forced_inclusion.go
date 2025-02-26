package suite_base

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"time"
)

type ForcedInclusionTestSpec struct {
	BaseTestSpec
}

func (f ForcedInclusionTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := f.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_forced_inclusion_test"

	return cfg
}

func (f ForcedInclusionTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node     = testnet.Nodes[0]
		l2Geth   = node.L2EthClient
		proposer = node.ProposerClient
	)

	if err := node.Start(); err != nil {
		t.Fatalf("BaseTestSpec failed to start 0 node: %v", err)
	}

	// For Debug
	if f.Debug {
		time.Sleep(time.Minute * 120)
	}

	// Wait until over the pacaya hardfork number.
	t.FailIfNotNil(l2Geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight), fmt.Sprintf("cannot wait for pacaya hardfork number: %d", proposer.PacayaClients.ForkHeight))

}
