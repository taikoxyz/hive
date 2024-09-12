package blob

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
)

type BlobTestSpec struct {
	TestL1Beacon   bool
	TestBlobSocial bool
	TestBlobServer bool
	suite_base.BaseTestSpec
}

func (r BlobTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()

	params.SetEnvParams("TEST_L1_BEACON", fmt.Sprintf("%v", r.TestL1Beacon))
	params.SetEnvParams("TEST_BLOB_SOCIAL", fmt.Sprintf("%v", r.TestBlobSocial))
	params.SetEnvParams("TEST_BLOB_SERVER", fmt.Sprintf("%v", r.TestBlobServer))

	params.SetEnvParams("RUN_TESTS", "")
	// enable blob tx.
	params.SetEnvParams("L1_BLOB_ALLOWED", "true")
	// set a mock beacon api url.
	params.SetEnvParams("L1_BEACON", "mock_url")
	params.SetEnvParams("BLOB_SOCIAL_SCAN_ENDPOINT", "mock_url")
	params.SetEnvParams("BLOB_SERVER", "mock_url")

	return cfg
}

func (r BlobTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {

}
