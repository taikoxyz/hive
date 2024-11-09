package suite_base

import (
	"fmt"
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

	// l2 geth syncmode
	L2SyncMode string

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

	syncMode := "snap"
	if ts.L2SyncMode != "" {
		syncMode = ts.L2SyncMode
	}

	return &testnet.Config{
		Eth1Consensus: execution_config.ExecutionCliqueConsensus{
			CliquePrivateKey: "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622",
			CliqueAddress:    "123463a4B065722E99115D6c222f267d9cABb524",
		},
		Network:    "taiko_base_test",
		LogLevel:   3,
		JWTSecret:  "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
		FeeReceipt: "a0Ee7A142d267C1f36714E4a8F75612F20a79720",
		L2SyncMode: syncMode,
		BeaconSync: ts.BeaconSync,
		Debug:      ts.Debug,
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
		ts.L2SyncMode,
		ts.IsGuardian,
		ts.BeaconSync,
		true,
	}
}
