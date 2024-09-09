package suite_base

import (
	execution_config "taiko/common/config/execution"
	"taiko/common/testnet"
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

	IsGuardian bool

	Debug bool
}

func (ts BaseTestSpec) GetTestnetConfig() *testnet.Config {
	return &testnet.Config{
		Eth1Consensus: execution_config.ExecutionCliqueConsensus{
			CliquePrivateKey: "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622",
			CliqueAddress:    "123463a4B065722E99115D6c222f267d9cABb524",
		},
		Network:    "taiko_base_test",
		LogLevel:   3,
		JWTSecret:  "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
		FeeReceipt: "a0Ee7A142d267C1f36714E4a8F75612F20a79720",
		L2SyncMode: ts.L2SyncMode,
		BeaconSync: true,
		IsGuardian: ts.IsGuardian,
		Debug:      ts.Debug,
	}
}

func (ts BaseTestSpec) GetName() string {
	return ts.Name
}

func (ts BaseTestSpec) GetDisplayName() string {
	return ts.DisplayName
}
