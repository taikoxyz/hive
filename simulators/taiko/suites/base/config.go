package suite_base

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	execution_config "taiko/common/config/execution"
	"taiko/common/testnet"
	"taiko/params"
	"taiko/suites"
)

type BaseTestSpec struct {
	// Spec
	Name        string
	DisplayName string
	Description string

	// L2 eth target number
	L2TargetNumber uint64

	// driver config
	IsGuardian bool
	BeaconSync bool

	Debug bool
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
		Debug:      ts.Debug,
		CreateConfig: func(index int, nodes clients.Nodes) (hivesim.Params, error) {
			var (
				firstNode    = nodes[0]
				node         = nodes[index]
				anvilClient  = firstNode.AnvilClient
				l1Client     = firstNode.L1EthClient
				beaconClient = firstNode.BeaconClient
				blobClient   = firstNode.BlobScanClient
				l2Client     = node.L2EthClient
				envs         = params.EnvParams()
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

			// This variables are set for blob tests.
			if envs["TEST_L1_BEACON"] == "true" {
				delete(envs, "RUN_TESTS")
			} else if envs["TEST_BLOB_SERVER"] == "true" {
				delete(envs, "L1_BEACON")
				envs["BLOB_SERVER"] = blobClient.BlobAPIURL()
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

func (ts BaseTestSpec) IsDebug() bool {
	return ts.Debug
}

func (ts BaseTestSpec) GetDisplayName() string {
	return ts.DisplayName
}

func (ts BaseTestSpec) DebugTestSpec() suites.TestSpec {
	return BaseTestSpec{
		fmt.Sprintf("Debug/%s", ts.Name),
		ts.DisplayName,
		ts.Description,
		ts.L2TargetNumber,
		ts.IsGuardian,
		ts.BeaconSync,
		true,
	}
}
