package testnet

import (
	"context"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"net"
	"taiko/common/clients"
	consensus_config "taiko/common/config/consensus"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	exec_client "github.com/marioevz/eth-clients/clients/execution"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	execution_config "taiko/common/config/execution"
)

type Testnet struct {
	*hivesim.T
	clients.Nodes

	genesisTime           common.Timestamp
	genesisValidatorsRoot common.Root

	// Consensus chain configuration
	spec *consensus_config.Spec
	// Execution chain configuration and genesis info
	executionGenesis *execution_config.ExecutionGenesis
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

func (t *Testnet) GenesisTime() common.Timestamp {
	// return time.Unix(int64(t.genesisTime), 0)
	return t.genesisTime
}

func (t *Testnet) GenesisTimeUnix() time.Time {
	return time.Unix(int64(t.genesisTime), 0)
}

func (t *Testnet) GenesisBeaconState() state.BeaconState {
	return t.executionGenesis.GenesisState
}

func (t *Testnet) GenesisValidatorsRoot() common.Root {
	return t.genesisValidatorsRoot
}

func (t *Testnet) ExecutionGenesis() *core.Genesis {
	return t.executionGenesis.Genesis
}

func StartTestnet(
	parentCtx context.Context,
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
		testnet     = prep.createTestnet(t)
		genesisTime = testnet.GenesisTimeUnix()
	)
	t.Logf(
		"Created new testnet, genesis at %s (%s from now)",
		genesisTime,
		time.Until(genesisTime),
	)

	var simulatorIP net.IP
	if simIPStr, err := t.Sim.ContainerNetworkIP(
		testnet.T.SuiteID,
		"bridge",
		"simulation",
	); err != nil {
		panic(err)
	} else {
		simulatorIP = net.ParseIP(simIPStr)
	}

	for nodeIndex, clientsByRole := range clientGroups {
		nodeClient := &clients.Node{
			Index:        nodeIndex,
			BeaconConfig: generateState.BeaconConfig,
			Genesis:      generateState.Genesis,
		}
		testnet.Nodes = append(testnet.Nodes, nodeClient)

		// Prepare clients for this node
		var (
			anvilDef     = clientsByRole[clients.Anvil]
			executionDef = clientsByRole[clients.Eth1]
		)

		if anvilDef != nil && executionDef != nil {
			t.Fatalf("Node %d has both anvil and execution clients, only one is allowed", nodeIndex)
		}

		// Prepare the client objects with all the information necessary to
		// eventually start
		nodeClient.AnvilClient = prep.prepareAnvilNode(
			config.Network,
			testnet,
			anvilDef,
		)
		nodeClient.L1EthClient = prep.prepareExecutionNode(
			testnet,
			executionDef,
			config.Eth1Consensus,
			clients.ExecutionClientConfig{
				ClientIndex: nodeIndex,
				JWTSecret:   ethcommon.HexToHash(config.JWTSecret),
				Network:     config.Network,
				ProxyConfig: &clients.ExecutionProxyConfig{
					Host: simulatorIP,
					Port: exec_client.PortEngineRPC + nodeIndex,
				},
			},
		)

		nodeClient.BeaconClient = prep.prepareBeaconNode(
			parentCtx,
			testnet,
			clientsByRole[clients.Beacon],
			&clients.BeaconClientConfig{
				ClientIndex:           nodeIndex,
				Spec:                  testnet.spec,
				GenesisValidatorsRoot: &testnet.genesisValidatorsRoot,
				GenesisTime:           &testnet.genesisTime,
				Network:               config.Network,
			},
			nodeClient.L1EthClient,
		)

		nodeClient.ValidatorClient = prep.prepareValidatorClient(
			testnet,
			clientsByRole[clients.Validator],
			nodeClient.BeaconClient,
		)

		nodeClient.L2EthClient = prep.prepareTaikoGethClient(
			config.Network,
			testnet,
			clientsByRole[clients.Eth2],
		)
		nodeClient.DriverClient = prep.prepareDriverClient(
			testnet,
			clientsByRole[clients.Driver],
			nodeClient.AnvilClient,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
		)
		nodeClient.ProposerClient = prep.prepareProposerClient(
			testnet,
			clientsByRole[clients.Proposer],
			nodeClient.AnvilClient,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
		)
		nodeClient.ProverClient = prep.prepareProverClient(
			testnet,
			clientsByRole[clients.Prover],
			nodeClient.AnvilClient,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
		)

		// Add rest of properties
		nodeClient.Logging = t
		nodeClient.Index = nodeIndex

		// Start the node clients if specified so
		// Connect to the network if specified
		if err = nodeClient.CreateNetwork(t, config.Network); err != nil {
			t.Fatalf("FAIL: Unable to connect to network: %v", err)
		}
		t.Logf("Starting node %d", nodeIndex)
		if err = nodeClient.Start(); err != nil {
			t.Fatalf("FAIL: Unable to start node %d: %v", nodeIndex, err)
		}
	}
	return testnet
}

func (t *Testnet) Stop() {
	for _, node := range t.Nodes {
		_ = node.Shutdown()
	}
}
