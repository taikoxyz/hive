package suite_base

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	execution_config "taiko/common/config/execution"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/common/utils"
	"taiko/params"
	"time"
)

func init() {
	Tests = append(Tests,
		BaseTestSpec{
			Name:           "fullsync",
			BeaconSync:     true,
			L2TargetNumber: 13,
		},
	)
}

type BaseTestSpec struct {
	// Spec
	Name        string
	DisplayName string

	// L2 eth target number
	L2TargetNumber uint64

	// driver config
	IsGuardian bool
	BeaconSync bool
}

func (ts BaseTestSpec) GetTestnetConfig() *testnet.Config {
	params.SetEnvParams("IS_GUARDIAN", fmt.Sprintf("%v", ts.IsGuardian))
	if !ts.IsGuardian {
		params.SetEnvParams("GUARDIAN_PROVER_MINORITY", "")
		params.SetEnvParams("GUARDIAN_PROVER_MAJORITY", "")
		params.SetEnvParams("GUARDIAN_PROVER_CONTRACT", "")
	}

	return &testnet.Config{
		Eth1Consensus: execution_config.ExecutionCliqueConsensus{
			CliquePrivateKey: "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622",
			CliqueAddress:    "123463a4B065722E99115D6c222f267d9cABb524",
		},
		Network:    "taiko_base_test",
		LogLevel:   3,
		FeeReceipt: "a0Ee7A142d267C1f36714E4a8F75612F20a79720",
		BeaconSync: ts.BeaconSync,
		Debug:      utils.GetenvBool("HIVE_TAIKO_DEBUG"),
		CreateConfig: func(index int, nodes clients.Nodes) (hivesim.Params, error) {
			var (
				firstNode      = nodes[0]
				node           = nodes[index]
				anvilClient    = firstNode.AnvilClient
				l1Client       = firstNode.L1EthClient
				beaconClient   = firstNode.BeaconClient
				blobscanClient = firstNode.BlobScanClient
				l2Client       = node.L2EthClient
				envs           = params.EnvParams()
			)

			envs["L2_AUTH"] = l2Client.EngineURL()
			envs["L2_HTTP"] = l2Client.HttpURL()
			envs["L2_WS"] = l2Client.WSURL()
			if anvilClient != nil {
				envs["L1_HTTP"] = anvilClient.HttpURL()
				envs["L1_WS"] = anvilClient.WSURL()
				envs["L1_BEACON"] = anvilClient.HttpURL()
			}
			if l1Client != nil {
				envs["L1_HTTP"] = l1Client.HttpURL()
				envs["L1_WS"] = l1Client.WSURL()
				envs["L1_BEACON"] = beaconClient.BeaconURL()
			}

			// Allow proposer use blob transaction builder.
			envs["L1_BLOB_ALLOWED"] = "true"
			// These variables are set for blob tests.
			if envs["TEST_L1_BEACON"] == "true" {
				delete(envs, "RUN_TESTS")
				//envs["BLOB_SERVER"] = beaconClient.BeaconURL()
			}
			if envs["TEST_BLOB_SCAN"] == "true" {
				delete(envs, "L1_BEACON")
				envs["BLOB_SOCIAL_SCAN_ENDPOINT"] = blobscanClient.BlobAPIURL()
			}

			envs["L1_PROPOSER_PRIV_KEY"] = params.ChainAuths[index*2+1].Key
			envs["L2_SUGGESTED_FEE_RECIPIENT"] = params.ChainAuths[index*2+1].Address.String()
			envs["L1_PROVER_PRIV_KEY"] = params.ChainAuths[index*2+2].Key

			params.ClusterEnvs[index] = envs

			return envs, nil
		},
	}
}

func (ts BaseTestSpec) GetName() string {
	return ts.Name
}

func (ts BaseTestSpec) GetDisplayName() string {
	return ts.DisplayName
}

func (ts BaseTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		nodes    = testnet.Nodes
		target   = ts.L2TargetNumber
		proposer = nodes[0].ProposerClient
		prover   = nodes[0].ProverClient
	)
	if target == 0 {
		target = rand.Uint64N(50-10) + 10
	}
	t.Logf("BaseTestSpec target number: %d", target)

	if err := nodes[0].Start(); err != nil {
		t.Fatalf("BaseTestSpec failed to start 0 node: %v", err)
	}
	for i, node := range nodes[1:] {
		if err := node.L2EthClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	for range time.Tick(time.Second) {
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID >= proposer.PacayaClients.ForkHeight {
			break
		}
		t.Nil(prover.VerifyBlocks(params.L1Auths[0]), "failed to verify blocks")
	}

	// For Debug
	if utils.GetenvBool("HIVE_TAIKO_DEBUG") {
		proposer.PauseClient()
		time.Sleep(time.Minute * 120)
	}

	// Start the other cluster's l2eth and driver nodes.
	for i, node := range nodes[1:] {
		if err := node.DriverClient.Start(); err != nil {
			t.Fatalf("BaseTestSpec failed to start %d node: %v", i, err)
		}
	}

	// Verify all l2eth nodes.
	ts.verifyL2Nodes(ctx, t, prover.GetLastVerifiedBlockId(ctx), nodes)
}

func (ts BaseTestSpec) verifyL2Nodes(ctx context.Context, t *hivesim.T, latestVerified uint64, nodes []*clients.Node) {
	var (
		index  int
		l2cli  = nodes[0].L2EthClient.EthClient
		prover = nodes[0].ProverClient
	)

	// todo: storeForcedInclusion test.

	t.FailIfNotNil(prover.WaitLatestVerifiedNumber(ctx, time.Second*10, latestVerified+1), fmt.Sprintf("failed to wait latest verified number: %d", latestVerified+1))

	header, err := l2cli.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatalf("failed to get l2 block number: %v", err)
	}
	target := header.Number.Uint64()

	for i, node := range nodes[1:] {
		l2client := node.L2EthClient
		client := l2client.EthClient

		hd, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(target))
		if err != nil {
			t.Fatalf("failed to get header from [%d]:%s, err: %v", i, node.L1EthClient.ClientType(), err)
		}
		if header.Hash() != hd.Hash() {
			t.Fatalf("the %d number of %s's hash are different, [%d]:%s != [%d]:%s",
				target,
				l2client.ClientType(),
				index, header.Hash().String(),
				i, hd.Hash().String(),
			)
		}

		for num := uint64(1); num <= target; num++ {
			l1Origin, err := l2client.L1OriginByID(ctx, big.NewInt(0).SetUint64(num))
			if err != nil {
				t.Fatalf("unexpect error when get l1origin, number: %d, err: %v", num, err)
			}
			if l1Origin == nil {
				t.Fatalf("l1Origin should not be null, number: %d", num)
			}
		}
	}
}
