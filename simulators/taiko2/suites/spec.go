package suites

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"taiko2/common/clients"
	"taiko2/common/config/execution"
	"taiko2/common/testnet"
	"taiko2/common/utils"
	"time"
)

var Deneb string = "deneb"

type TestSpec interface {
	GetName() string
	GetTestnetConfig(clients.NodeDefinitions) *testnet.Config
	GetDisplayName() string
	GetDescription() *utils.Description
}

// Add all tests to the suite
func SuiteHydrate(
	suite *hivesim.Suite,
	clients *clients.ClientDefinitionsByRole,
	tests []TestSpec,
	generateState *execution_config.GenesisState,
) {
	clientCombinations := clients.Combinations()
	for _, test := range tests {
		test := test
		suite.Add(hivesim.TestSpec{
			Name:        fmt.Sprintf("%s/%s", suite.Name, test.GetName()),
			DisplayName: test.GetDisplayName(),
			Description: test.GetDescription().Format(),
			Run: func(t *hivesim.T) {
				t.Logf("Starting test: %s", test.GetName())
				defer t.Logf("Finished test: %s", test.GetName())
				env := &testnet.Environment{
					Clients: clients,
				}

				t.Logf("Starting testnet with %d nodes", len(clientCombinations))
				config := test.GetTestnetConfig(clientCombinations)

				// Create the testnet
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				testnet := testnet.StartTestnet(ctx, t, env, config, generateState)
				if testnet == nil {
					t.Fatalf("failed to start testnet")
				}
				defer testnet.Stop()

				time.Sleep(time.Second * 50)
			},
		},
		)
	}
}
