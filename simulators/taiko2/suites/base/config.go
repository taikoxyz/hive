package suite_base

import (
	"fmt"
	execution_config "taiko2/common/config/execution"

	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"taiko2/common/clients"
	cl "taiko2/common/config/consensus"
	"taiko2/common/testnet"
	"taiko2/common/utils"
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

	// Actions
	ExitValidatorsShare int

	// Verifications
	EpochsAfterFork beacon.Epoch
	WaitForBlobs    bool
	WaitForFinality bool

	// Extra Gwei
	ExtraGwei beacon.Gwei
}

var (
	DEFAULT_VALIDATOR_COUNT uint64 = 64
)

func (ts BaseTestSpec) GetNodeCount() int {
	if ts.NodeCount > 0 {
		return ts.NodeCount
	}
	return 2
}

func (ts BaseTestSpec) GetValidatingNodeCount() int {
	if ts.ValidatingNodeCount > 0 {
		return ts.ValidatingNodeCount
	}
	return ts.GetNodeCount()
}

func (ts BaseTestSpec) GetTestnetConfig(
	allNodeDefinitions clients.NodeDefinitions,
) *testnet.Config {
	maxValidatingNodeIndex := ts.GetValidatingNodeCount()
	nodeDefinitions := make(clients.NodeDefinitions, 0)
	for i := 0; i < ts.GetNodeCount(); i++ {
		n := allNodeDefinitions[i%len(allNodeDefinitions)]
		if i < maxValidatingNodeIndex {
			n.ValidatorShares = 1
		} else {
			n.ValidatorShares = 0
		}
		nodeDefinitions = append(nodeDefinitions, n)
	}

	return &testnet.Config{
		Eth1Consensus: execution_config.ExecutionCliqueConsensus{
			CliquePrivateKey: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			CliqueAddress:    "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
		NodeDefinitions: nodeDefinitions,
		Network:         "taiko2_base_test",
		JWTSecret:       "7365637265747365637265747365637265747365637265747365637265747365",
		FeeReceipt:      "123463a4b065722e99115d6c222f267d9cabb524",
	}
}

func (ts BaseTestSpec) CanRun(clients.NodeDefinitions) bool {
	// Base test specs can always run
	return true
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
	if ts.WaitForFinality {
		desc.Add(utils.CategoryVerificationsConsensusClient, `
		- After all other verifications are done, the beacon chain is able to finalize the current epoch`)
	}

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

func (ts BaseTestSpec) GetValidatorKeys(
	mnemonic string,
) cl.ValidatorsSetupDetails {
	keySrc := &cl.MnemonicsKeySource{
		From:     0,
		To:       ts.GetValidatorCount(),
		Mnemonic: mnemonic,
	}
	keys, err := keySrc.Keys()
	if err != nil {
		panic(err)
	}

	for index, key := range keys {
		// All validators have idiosyncratic balance amounts to identify them.
		// Also include a high amount in order to guarantee withdrawals.
		key.ExtraInitialBalance = beacon.Gwei((index+1)*1000000) + ts.ExtraGwei

		if ts.GenesisExecutionWithdrawalCredentialsShares > 0 &&
			(index%ts.GenesisExecutionWithdrawalCredentialsShares) == 0 {
			key.WithdrawalCredentialType = beacon.ETH1_ADDRESS_WITHDRAWAL_PREFIX
			key.WithdrawalExecAddress = beacon.Eth1Address{byte(index + 0x100)}
		}
		if ts.GenesisExitedShares > 1 && (index%ts.GenesisExitedShares) == 1 {
			key.Exited = true
		}
		if ts.GenesisSlashedShares > 2 &&
			(index%ts.GenesisSlashedShares) == 2 {
			key.Slashed = true
		}
		fmt.Printf(
			"INFO: Validator %d, extra_gwei=%d, exited=%v, slashed=%v, key_type=%d\n",
			index,
			key.ExtraInitialBalance,
			key.Exited,
			key.Slashed,
			key.WithdrawalCredentialType,
		)
	}

	return keys
}
