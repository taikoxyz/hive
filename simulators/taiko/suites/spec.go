package suites

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	"taiko/common/testnet"
)

var Deneb string = "deneb"

type TestSpec interface {
	GetName() string
	GetTestnetConfig() *testnet.Config
	GetDisplayName() string
	Verify(ctx context.Context, t *hivesim.T, testnet *testnet.Testnet)
}

// SuiteHydrate Add all tests to the suite
func SuiteHydrate(
	suite *hivesim.Suite,
	clients clients.ClientGroups,
	tests []TestSpec,
) {

	for _, test := range tests {
		test := test
		suite.Add(hivesim.TestSpec{
			Name:        test.GetName(),
			DisplayName: test.GetDisplayName(),
			Run: func(t *hivesim.T) {
				t.Logf("Starting test: %s", test.GetName())
				defer t.Logf("Finished test: %s", test.GetName())

				testnet := testnet.StartTestnet(t, clients, test.GetTestnetConfig())
				if testnet == nil {
					t.Fatalf("failed to start testnet")
				}
				defer testnet.Stop()

				// Verify nodes.
				test.Verify(context.Background(), t, testnet)
			},
		})
	}
}
