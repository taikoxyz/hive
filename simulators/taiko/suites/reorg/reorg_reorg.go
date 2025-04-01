package suite_reorg

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		ReorgTestSpec{
			BaseTestSpec: suite_base.BaseTestSpec{
				Name:           "reorg_reorg",
				L2TargetNumber: 13,
			},
		},
	)
}

type ReorgTestSpec struct {
	suite_base.BaseTestSpec

	// reorg params
	params *clients.ReorgParams
}

func (r ReorgTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_reorg_test"
	return cfg
}

func (r ReorgTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if testnet.DevDebug {
		node.DriverClient.Shutdown()
		node.ProverClient.Shutdown()
		time.Sleep(time.Second * 10)
	}
	r.params = &clients.ReorgParams{
		DelayTime: 1,
	}

	r.reorg(ctx, t, node)
}

func (r ReorgTestSpec) reorg(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		anvil  = node.AnvilClient
		prover = node.ProverClient
		l2eth  = node.L2EthClient

		timeout                   = time.Minute * 3
		l2ReorgStartNumber        = r.L2TargetNumber
		reorgDeep          uint64 = 5
	)
	t.Logf("%s: start reorgAndVerifyFirstCluster, target number: %d", r.Name, l2ReorgStartNumber)

	// Start recording reorg points.
	anvil.StartRecordReorgPoints(ctx, l2eth.EthClient)

	for range time.Tick(time.Second) {
		prover.VerifyBlocks(params.L1Auths[0])
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID > 0 {
			node.ProverClient.Shutdown()
			break
		}
	}

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber+reorgDeep)

	// Reorg l1 eth chain.
	anvil.Reorg(l2ReorgStartNumber, r.params)

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber+reorgDeep+3)

	// Verify l1Origins.
	t.FailIfNotNil(clients.VerifyL1Origin(ctx, l2ReorgStartNumber, anvil.EthClient, l2eth.EthClient))
}
