package suite_base

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"math/rand/v2"
	"os"
	"taiko/common/clients"
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
}

func (ts BaseTestSpec) GetTestnetConfig() *testnet.Config {
	params.SetEnvParams("IS_GUARDIAN", fmt.Sprintf("%v", ts.IsGuardian))
	if !ts.IsGuardian {
		params.SetEnvParams("GUARDIAN_PROVER_MINORITY", "")
		params.SetEnvParams("GUARDIAN_PROVER_MAJORITY", "")
		params.SetEnvParams("GUARDIAN_PROVER_CONTRACT", "")
	}

	return &testnet.Config{
		Network:   "taiko_base_test",
		LogLevel:  os.Getenv("HIVE_LOGLEVEL"),
		DevDebug:  utils.GetenvBool("HIVE_DEV_DEBUG"),
		DevExpose: utils.GetenvBool("HIVE_DEV_EXPOSE"),
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

			if firstNode.AnvilClient != nil {
				envs["HIVE_L1_NODE"] = "anvil"
			}
			if firstNode.L1EthClient != nil {
				envs["HIVE_L1_NODE"] = "geth"
			}

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
			// These variables are set for blob tests.
			if envs["TEST_L1_BEACON"] == "true" {
				delete(envs, "RUN_TESTS")
				envs["L1_BLOB_ALLOWED"] = "true"
			}
			if envs["TEST_BLOB_SERVER"] == "true" {
				delete(envs, "L1_BEACON")
				envs["L1_BLOB_ALLOWED"] = "true"
				envs["BLOB_SERVER"] = blobscanClient.BlobAPIURL()
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
		prover.VerifyBlocks(params.L1Auths[0])
	}

	// For DevDebug
	if testnet.DevDebug {
		proposer.PauseClient()
		time.Sleep(time.Hour * 2)
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

	prover.WaitLatestVerifiedNumber(ctx, time.Minute*3, latestVerified+1)

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
