package blob

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
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
	cfg.Network = "taiko_blob_test"

	// enable blob tx.
	params.SetEnvParams("TEST_L1_BEACON", fmt.Sprintf("%v", r.TestL1Beacon))
	params.SetEnvParams("TEST_BLOB_SCAN", fmt.Sprintf("%v", r.TestBlobServer))

	cfgFun := cfg.CreateConfig
	cfg.CreateConfig = func(index int, nodes clients.Nodes) (hivesim.Params, error) {
		envs, err := cfgFun(index, nodes)
		if err != nil {
			return nil, err
		}

		var (
			firstNode      = nodes[0]
			beaconClient   = firstNode.BeaconClient
			blobscanClient = firstNode.BlobScanClient
		)

		// Allow proposer use blob transaction builder.
		envs["L1_BLOB_ALLOWED"] = "true"

		// These variables are set for blob tests.
		if envs["TEST_L1_BEACON"] == "true" {
			delete(envs, "RUN_TESTS")
			delete(envs, "L1_BEACON")
			envs["BLOB_SERVER"] = beaconClient.BeaconURL()
		}
		if envs["TEST_BLOB_SCAN"] == "true" {
			delete(envs, "L1_BEACON")
			envs["BLOB_SOCIAL_SCAN_ENDPOINT"] = blobscanClient.BlobAPIURL()
		}

		params.ClusterEnvs[index] = envs

		return envs, nil
	}

	return cfg
}

func (r BlobTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if r.Debug {
		node.DriverClient.Shutdown()
		time.Sleep(time.Minute * 60)
	}

	var (
		l2eth        = node.L2EthClient
		timeout      = time.Second * 60
		targetNumber = uint64(10)
	)

	l2eth.WaitLatestNumber(ctx, timeout, targetNumber)
}
