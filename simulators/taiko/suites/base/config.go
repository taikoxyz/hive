package suite_base

import (
	"fmt"
	execution_config "taiko/common/config/execution"

	"taiko/common/testnet"
	"taiko/common/utils"
)

type BaseTestSpec struct {
	// Spec
	Name        string
	DisplayName string
	Description string

	// Testnet Nodes
	NodeCount           int
	ValidatingNodeCount int

	// Beacon Chain
	ValidatorCount uint64
	DenebGenesis   bool

	// Genesis Validators Configuration
	// (One every Nth validator, 1 means all validators, 2 means half, etc...)
	GenesisExecutionWithdrawalCredentialsShares int
	GenesisExitedShares                         int
	GenesisSlashedShares                        int
}

var (
	DEFAULT_VALIDATOR_COUNT uint64 = 64
)

func (ts BaseTestSpec) GetNodeCount() int {
	if ts.NodeCount > 0 {
		return ts.NodeCount
	}
	return 1
}

func (ts BaseTestSpec) GetValidatingNodeCount() int {
	if ts.ValidatingNodeCount > 0 {
		return ts.ValidatingNodeCount
	}
	return ts.GetNodeCount()
}

func (ts BaseTestSpec) GetTestnetConfig() *testnet.Config {
	return &testnet.Config{
		Eth1Consensus: execution_config.ExecutionCliqueConsensus{
			CliquePrivateKey: "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622",
			CliqueAddress:    "123463a4B065722E99115D6c222f267d9cABb524",
		},
		Network:    "taiko_base_test",
		LogLevel:   3,
		JWTSecret:  "7365637265747365637265747365637265747365637265747365637265747365",
		FeeReceipt: "a0Ee7A142d267C1f36714E4a8F75612F20a79720",
		ProposerKeyList: []string{
			"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
			"0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		},
		ProverKeyList: []string{
			"0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
			"0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
			"0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
			"0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
			"0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
		},
	}
}

func (ts BaseTestSpec) GetName() string {
	return ts.Name
}

func (ts BaseTestSpec) GetDisplayName() string {
	return ts.DisplayName
}

func (ts BaseTestSpec) GetDescription() *utils.Description {
	desc := utils.NewDescription(ts.Description)

	// Add the testnet config description
	desc.Add(utils.CategoryTestnetConfiguration, fmt.Sprintf(`
	  - Node Count: %d
	  - Validating Node Count: %d
	  - Validator Key Count: %d
	  - Validator Key per Node: %d`,
		ts.GetNodeCount(),
		ts.GetValidatingNodeCount(),
		ts.GetValidatorCount(),
		ts.GetValidatorCount()/uint64(ts.GetValidatingNodeCount()),
	))
	if ts.DenebGenesis {
		desc.Add(utils.CategoryTestnetConfiguration, "\n- Genesis Fork: Deneb")
	} else {
		desc.Add(utils.CategoryTestnetConfiguration, "\n- Genesis Fork: Capella")
	}
	execCredentialCount := ts.GetExecutionWithdrawalCredentialCount()
	blsCredentialCount := ts.GetValidatorCount() - execCredentialCount
	desc.Add(utils.CategoryTestnetConfiguration, fmt.Sprintf(`
	  - Execution Withdrawal Credentials Count: %d
	  - BLS Withdrawal Credentials Count: %d`,
		execCredentialCount,
		blsCredentialCount,
	))

	// Add the verifications description
	desc.Add(utils.CategoryVerificationsExecutionClient, `
	  - Blob (type-3) transactions are included in the blocks`)
	desc.Add(utils.CategoryVerificationsConsensusClient, `
	  - For each blob transaction on the execution chain, the blob sidecars are available for the beacon block at the same height
	  - The beacon block lists the correct commitments for each blob`)

	return desc
}

func (ts BaseTestSpec) GetValidatorCount() uint64 {
	if ts.ValidatorCount != 0 {
		return ts.ValidatorCount
	}
	return DEFAULT_VALIDATOR_COUNT
}

func (ts BaseTestSpec) GetExecutionWithdrawalCredentialCount() uint64 {
	if ts.GenesisExecutionWithdrawalCredentialsShares != 0 {
		return ts.GetValidatorCount() / uint64(ts.GenesisExecutionWithdrawalCredentialsShares)
	}
	return 0
}
