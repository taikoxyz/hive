package suites

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"strings"
	"taiko2/common/clients"
	consensus_config "taiko2/common/config/consensus"
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
	GetValidatorKeys(string) consensus_config.ValidatorsSetupDetails
}

// Add all tests to the suite
func SuiteHydrate(
	suite *hivesim.Suite,
	c *clients.ClientDefinitionsByRole,
	tests []TestSpec,
	generateState *execution_config.GenesisState,
) {
	mnemonic := "couple kiwi radio river setup fortune hunt grief buddy forward perfect empty slim wear bounce drift execute nation tobacco dutch chapter festival ice fog"

	clientCombinations := c.Combinations()
	for _, test := range tests {
		test := test
		suite.Add(hivesim.TestSpec{
			Name: fmt.Sprintf(
				"%s-%s",
				test.GetName(),
				strings.Join(clientCombinations.ClientTypes(), "-"),
			),
			DisplayName: test.GetDisplayName(),
			Description: test.GetDescription().Format(),
			Run: func(t *hivesim.T) {
				t.Logf("Starting test: %s", test.GetName())
				defer t.Logf("Finished test: %s", test.GetName())
				keys := test.GetValidatorKeys(mnemonic)
				env := &testnet.Environment{
					Clients:    c,
					Validators: keys,
				}
				config := test.GetTestnetConfig(clientCombinations)

				// Create the testnet
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				testnet := testnet.StartTestnet(ctx, t, env, config, generateState)
				if testnet == nil {
					t.Fatalf("failed to start testnet")
				}
				defer testnet.Stop()

				time.Sleep(time.Second * 200)
			},
		},
		)
	}
}
