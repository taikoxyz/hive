package blob

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	"taiko/suites"
	suite_base "taiko/suites/base"
	"time"
)

type BlobTestSpec struct {
	TestL1Beacon   bool
	TestBlobServer bool
	suite_base.BaseTestSpec
}

func (r BlobTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()

	// enable blob tx.
	params.SetEnvParams("L1_BLOB_ALLOWED", "true")

	params.SetEnvParams("TEST_L1_BEACON", fmt.Sprintf("%v", r.TestL1Beacon))
	params.SetEnvParams("TEST_BLOB_SERVER", fmt.Sprintf("%v", r.TestBlobServer))

	return cfg
}

func (r BlobTestSpec) DebugTestSpec() suites.TestSpec {
	return BlobTestSpec{
		r.TestL1Beacon,
		r.TestBlobServer,
		r.BaseTestSpec.DebugTestSpec().(suite_base.BaseTestSpec),
	}
}

func (r BlobTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if r.IsDebug() {
		time.Sleep(time.Minute * 60)
	}

	var (
		l2eth        = node.L2EthClient
		timeout      = time.Second * 60
		targetNumber = uint64(10)
	)

	if err := l2eth.WaitLatestNumber(ctx, timeout, targetNumber); err != nil {
		t.Fatalf("%s: can't wait %s touch target height, err: %v", r.Name, l2eth.ClientType(), err)
	}
}
