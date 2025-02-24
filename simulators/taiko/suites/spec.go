package suites

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/common/config/execution"
	"taiko/common/testnet"
)

var Deneb string = "deneb"

type TestSpec interface {
	GetName() string
	GetTestnetConfig() *testnet.Config
	GetDisplayName() string
	Verify(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet)
	DebugTestSpec() TestSpec
}

// SuiteHydrate Add all tests to the suite
func SuiteHydrate(
	suite *hivesim.Suite,
	clients clients.ClientGroups,
	tests []TestSpec,
	generateState *execution_config.GenesisState,
) {

	for _, test := range tests {
		test := test
		suite.Add(hivesim.TestSpec{
			Name:        test.GetName(),
			DisplayName: test.GetDisplayName(),
			Run: func(t *hivesim.T) {
				t.Logf("Starting test: %s", test.GetName())
				defer t.Logf("Finished test: %s", test.GetName())

				// Create the testnet
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				testnet := testnet.StartTestnet(t, clients, test.GetTestnetConfig(), generateState)
				if testnet == nil {
					t.Fatalf("failed to start testnet")
				}
				defer testnet.Stop()

				// Verify nodes.
				test.Verify(ctx, t, testnet)
			},
		})
	}
}
