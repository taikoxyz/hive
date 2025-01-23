package testnet

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"taiko/common/clients"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/config/consensus"
	"taiko/common/config/execution"
)

type Testnet struct {
	*hivesim.T
	clients.Nodes

	genesisTime           uint64
	genesisValidatorsRoot common.Hash

	// Consensus chain configuration
	spec *consensus_config.Spec
	// Execution chain configuration and genesis info
	executionGenesis *execution_config.ExecutionGenesis

	// debug flag
	Debug bool
}

type ActiveSpec struct {
	*consensus_config.Spec
}

const slotsTolerance primitives.Slot = 2

func (spec *ActiveSpec) EpochTimeoutContext(
	parent context.Context,
	epochs primitives.Epoch,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		parent,
		time.Duration(
			uint64((spec.SlotsPerEpoch*primitives.Slot(epochs))+slotsTolerance)*
				uint64(spec.SecondsPerSlot),
		)*time.Second,
	)
}

func (spec *ActiveSpec) SlotTimeoutContext(
	parent context.Context,
	slots primitives.Slot,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		parent,
		time.Duration(
			uint64(slots+slotsTolerance)*
				uint64(spec.SecondsPerSlot))*time.Second,
	)
}

func (spec *ActiveSpec) EpochsTimeout(epochs primitives.Epoch) <-chan time.Time {
	return time.After(
		time.Duration(
			uint64(
				spec.SlotsPerEpoch*primitives.Slot(epochs),
			)*uint64(
				spec.SecondsPerSlot,
			),
		) * time.Second,
	)
}

func (spec *ActiveSpec) SlotsTimeout(slots primitives.Slot) <-chan time.Time {
	return time.After(
		time.Duration(
			uint64(slots)*uint64(spec.SecondsPerSlot),
		) * time.Second,
	)
}

func (t *Testnet) Spec() *ActiveSpec {
	return &ActiveSpec{
		Spec: t.spec,
	}
}

func (t *Testnet) GenesisTimeUnix() time.Time {
	return time.Unix(int64(t.genesisTime), 0)
}

func (t *Testnet) GenesisBeaconState() state.BeaconState {
	return t.executionGenesis.GenesisState
}

func (t *Testnet) GenesisValidatorsRoot() [32]byte {
	return t.genesisValidatorsRoot
}

func (t *Testnet) ExecutionGenesis() *core.Genesis {
	return t.executionGenesis.Genesis
}

func StartTestnet(
	t *hivesim.T,
	clientGroups clients.ClientGroups,
	config *Config,
	generateState *execution_config.GenesisState,
) *Testnet {
	prep, err := PrepareTestnet(config, generateState)
	if err != nil {
		t.Fatalf("FAIL: Unable to prepare testnet: %v", err)
	}
	var (
		testnet     = prep.createTestnet(t, config)
		genesisTime = testnet.GenesisTimeUnix()
	)
	t.Logf(
		"Created new testnet, genesis at %s (%s from now)",
		genesisTime,
		time.Until(genesisTime),
	)

	for index, clientsByRole := range clientGroups {
		nodeClient := &clients.Node{
			Logging: t,
			Index:   index,
			Genesis: generateState.Genesis,
		}
		testnet.Nodes = append(testnet.Nodes, nodeClient)

		// Prepare clients for this node
		var (
			anvilDef     = clientsByRole[clients.Anvil]
			executionDef = clientsByRole[clients.Eth1]
		)

		if anvilDef != nil && executionDef != nil {
			t.Fatalf("Node %d has both anvil and execution clients, only one is allowed", index)
		}

		// Prepare the client objects with all the information necessary to
		// eventually start
		prep.prepareAnvilNode(testnet, config, anvilDef)
		prep.prepareExecutionNode(index, testnet, config, executionDef, config.Eth1Consensus)
		prep.prepareBeaconNode(testnet, config, clientsByRole[clients.Beacon])
		prep.prepareValidatorClient(testnet, clientsByRole[clients.Validator])
		prep.prepareBlobScanClient(index, testnet, config, clientsByRole)

		// taiko-geth or taiko-reth
		taikoGeth := clientsByRole[clients.Eth2]
		if taikoGeth == nil {
			taikoGeth = clientsByRole[clients.Reth]
		}
		prep.prepareTaikoGethClient(index, testnet, config, taikoGeth)
		prep.prepareDriverClient(index, testnet, config, clientsByRole[clients.Driver])
		prep.prepareProposerClient(index, testnet, clientsByRole[clients.Proposer])
		prep.prepareProverClient(index, testnet, clientsByRole[clients.Prover])

		// Start the node clients if specified so
		// Connect to the network if specified
		if err = nodeClient.CreateNetwork(t, config.Network); err != nil {
			t.Fatalf("FAIL: Unable to connect to network: %v", err)
		}
		t.Logf("node %d is ready to use!", index)
	}
	return testnet
}

func (t *Testnet) Stop() {
	for _, node := range t.Nodes {
		_ = node.Shutdown()
	}
}
