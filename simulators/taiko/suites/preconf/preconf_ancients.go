package preconf

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"time"
)

func init() {
	preconf := PreconfAncientsTestSpec{}
	preconf.Name = "ancients"
	Tests = append(Tests, preconf)
}

type PreconfAncientsTestSpec struct {
	PreconfTestSpec
}

func (r PreconfAncientsTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_preconf_ancients"

	return cfg
}

func (r PreconfAncientsTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node     = testnet.Nodes[0]
		proposer = node.ProposerClient
		l2geth   = node.L2EthClient
	)

	// Start all the cluster's nodes.
	for _, node := range testnet.Nodes {
		t.Nil(node.Start(), "cannot start L2EthClient")
	}

	l2geth.WaitLatestNumber(ctx, time.Minute*3, proposer.PacayaClients.ForkHeight-1)

	// stop the proposer.
	t.FailIfNotNil(proposer.Shutdown())

	// Create a batch of preconf request bodies.

}
