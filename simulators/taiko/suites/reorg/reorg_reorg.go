package suite_reorg

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
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
		time.Sleep(time.Minute * 2)
	}
	r.params = &clients.ReorgParams{
		DelayTime: 1,
	}

	r.reorg(ctx, t, node)
}

func (r ReorgTestSpec) reorg(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		anvil    = node.AnvilClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
		l2eth    = node.L2EthClient

		timeout                   = time.Minute * 3
		l2ReorgStartNumber        = r.L2TargetNumber
		reorgDeep          uint64 = 5
	)
	if l2ReorgStartNumber == 0 {
		// random [10, 50) value
		l2ReorgStartNumber = rand.Uint64N(50-10) + 10
	}
	t.Logf("%s: start reorgAndVerifyFirstCluster, target number: %d", r.Name, l2ReorgStartNumber)

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber)

	// Start recording reorg points.
	anvil.StartRecordReorgPoints(ctx, l2eth.EthClient)

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber+reorgDeep)

	// pause driver, proposer, prover
	driver.PauseClient()
	proposer.PauseClient()

	// Reorg l1 eth chain.
	anvil.Reorg(l2ReorgStartNumber, r.params)

	// unpause driver, proposer, prover
	driver.UnpauseClient()
	proposer.UnpauseClient()

	l2eth.WaitLatestNumber(ctx, timeout, l2ReorgStartNumber*2)

	// Verify l1Origins.
	verifyL1Origin(ctx, t, l2ReorgStartNumber, reorgDeep, l2eth.EthClient)
}

func verifyL1Origin(ctx context.Context, t *hivesim.T, l2ReorgStartNumber, reorgDeep uint64, l1client *rpc.EthClient) {
	var l1Headers = map[uint64]*types.Header{}
	for number := l2ReorgStartNumber + 1; number <= l2ReorgStartNumber+reorgDeep; number++ {
		header, err := l1client.HeaderByNumber(ctx, big.NewInt(int64(number)))
		t.FailIfNotNil(err, "failed to get l1 origin")
		l1Headers[number] = header
	}

	// Verify l1Origins.
	for number := l2ReorgStartNumber + 1; number <= l2ReorgStartNumber+5; number++ {
		l1Origin, err := l1client.L1OriginByID(ctx, big.NewInt(int64(number)))
		t.FailIfNotNil(err, "failed to get l1 origin")

		l1Number := l1Origin.L1BlockHeight.Uint64()
		t.Equal(l1Headers[l1Number].Hash().String(), l1Origin.L1BlockHash.String(), fmt.Sprintf("l1Origin content is not right, l1BlockHeight: %d, l1BlockHash: %s", l1Number, l1Headers[l1Number].Hash().String()))
	}
}
