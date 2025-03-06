package blob

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/common/utils"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		BlobTestSpec{
			TestL1Beacon: true,
			BaseTestSpec: suite_base.BaseTestSpec{
				Name: "blob-l1-beacon",
			},
		},
		BlobTestSpec{
			TestBlobServer: true,
			BaseTestSpec: suite_base.BaseTestSpec{
				Name: "blob-server",
			},
		},
	)
}

type BlobTestSpec struct {
	TestL1Beacon   bool
	TestBlobServer bool
	suite_base.BaseTestSpec
}

func (r BlobTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	cfg.Network = "taiko_blob_test"

	// enable blob tx.
	params.SetEnvParams("TEST_L1_BEACON", fmt.Sprintf("%v", r.TestL1Beacon))
	params.SetEnvParams("TEST_BLOB_SERVER", fmt.Sprintf("%v", r.TestBlobServer))

	return cfg
}

func (r BlobTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if utils.GetenvBool("HIVE_DEBUG") {
		node.DriverClient.Shutdown()
		node.ProposerClient.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	var (
		l2eth        = node.L2EthClient
		targetNumber = uint64(13)
	)

	l2eth.WaitLatestNumber(ctx, time.Second*60, targetNumber)
}
