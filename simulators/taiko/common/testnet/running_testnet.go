package testnet

import (
	"context"
	"encoding/hex"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"net"
	"taiko/common/clients"
	consensus_config "taiko/common/config/consensus"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	exec_client "github.com/marioevz/eth-clients/clients/execution"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	execution_config "taiko/common/config/execution"
	"taiko/common/utils"
)

var (
	JWT_SECRET, _ = hex.DecodeString("7365637265747365637265747365637265747365637265747365637265747365")
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
	// Consensus genesis state
	//eth2GenesisState common.BeaconState

	// Blobber
	//blobber *blobber.Blobber

	// Test configuration
	maxConsecutiveErrorsOnWaits int

	// Validators
	Validators      []*ethpb.Validator
	ValidatorGroups map[string]*utils.Validators
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
	env *Environment,
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

	testnet.Nodes = make(clients.Nodes, len(config.NodeDefinitions))

	// Init all client bundles
	for nodeIndex := range testnet.Nodes {
		testnet.Nodes[nodeIndex] = new(clients.Node)
	}

	// For each key partition, we start a client bundle that consists of:
	// - 1 execution client
	// - 1 beacon client
	// - 1 validator client,
	for nodeIndex, node := range config.NodeDefinitions {
		// Prepare clients for this node
		var (
			nodeClient = testnet.Nodes[nodeIndex]

			anvilDef     = env.Clients.ClientByNameAndRole(node.L1EthClient, "anvil")
			executionDef = env.Clients.ClientByNameAndRole(node.L1EthClient, "eth1")
			beaconDef    = env.Clients.ClientByNameAndRole(node.ConsensusClient, "beacon")
			validatorDef = env.Clients.ClientByNameAndRole(node.ValidatorClientName(), "validator")
			taikoGethDef = env.Clients.ClientByNameAndRole(node.L2EthClient, "taiko-geth")
			driverDef    = env.Clients.ClientByNameAndRole(node.DriverClient, "driver")
			proposerDef  = env.Clients.ClientByNameAndRole(node.ProposerClient, "proposer")
			proverDef    = env.Clients.ClientByNameAndRole(node.ProverClient, "prover")
			executionTTD = int64(0)
			beaconTTD    = int64(0)
		)

		if node.ExecutionClientTTD != nil {
			executionTTD = node.ExecutionClientTTD.Int64()
		} else if testnet.executionGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
			executionTTD = testnet.executionGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		}
		if node.BeaconNodeTTD != nil {
			beaconTTD = node.BeaconNodeTTD.Int64()
		} else if testnet.executionGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
			beaconTTD = testnet.executionGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		}

		// Prepare the client objects with all the information necessary to
		// eventually start
		if anvilDef != nil {
			nodeClient.AnvilClient = prep.prepareAnvilNode(
				parentCtx,
				config.Network,
				testnet,
				anvilDef,
			)
		}
		if executionDef != nil {
			nodeClient.L1EthClient = prep.prepareExecutionNode(
				parentCtx,
				testnet,
				executionDef,
				config.Eth1Consensus,
				node.Chain,
				clients.ExecutionClientConfig{
					ClientIndex:             nodeIndex,
					TerminalTotalDifficulty: executionTTD,
					Subnet:                  node.GetExecutionSubnet(),
					JWTSecret:               ethcommon.HexToHash(config.JWTSecret),
					Network:                 config.Network,
					ProxyConfig: &clients.ExecutionProxyConfig{
						Host:                   simulatorIP,
						Port:                   exec_client.PortEngineRPC + nodeIndex,
						TrackForkchoiceUpdated: false,
						LogEngineCalls:         env.LogEngineCalls,
					},
				},
			)
		}

		nodeClient.BeaconClient = prep.prepareBeaconNode(
			parentCtx,
			testnet,
			beaconDef,
			&clients.BeaconClientConfig{
				ClientIndex:             nodeIndex,
				TerminalTotalDifficulty: beaconTTD,
				Spec:                    testnet.spec,
				GenesisValidatorsRoot:   &testnet.genesisValidatorsRoot,
				GenesisTime:             &testnet.genesisTime,
				Subnet:                  node.GetConsensusSubnet(),
				Network:                 config.Network,
			},
			nodeClient.L1EthClient,
		)

		nodeClient.ValidatorClient = prep.prepareValidatorClient(
			parentCtx,
			testnet,
			validatorDef,
			nodeClient.BeaconClient,
			nodeIndex,
		)

		nodeClient.L2EthClient = prep.prepareTaikoGethClient(
			parentCtx,
			config.Network,
			testnet,
			taikoGethDef,
		)
		nodeClient.DriverClient = prep.prepareDriverClient(
			parentCtx,
			testnet,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
			driverDef,
		)
		nodeClient.ProverClient = prep.prepareProverClient(
			parentCtx,
			testnet,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
			proverDef,
		)
		nodeClient.ProposerClient = prep.prepareProposerClient(
			parentCtx,
			testnet,
			nodeClient.L1EthClient,
			nodeClient.BeaconClient,
			nodeClient.L2EthClient,
			proposerDef,
		)

		// Add rest of properties
		nodeClient.Logging = t
		nodeClient.Index = nodeIndex
		// Start the node clients if specified so
		if !node.DisableStartup {
			// Connect to the network if specified
			if err = nodeClient.CreateNetwork(t, config.Network); err != nil {
				t.Fatalf("FAIL: Unable to connect to network: %v", err)
			}
			t.Logf("Starting node %d", nodeIndex)
			if err = nodeClient.Start(); err != nil {
				t.Fatalf("FAIL: Unable to start node %d: %v", nodeIndex, err)
			}
		} else {
			t.Logf("Node %d startup disabled, skipping", nodeIndex)
		}
	}

	return testnet
}

func (t *Testnet) Stop() {
	for _, node := range t.Nodes {
		node.Shutdown()
	}
	for _, p := range t.Proxies().Running() {
		p.Cancel()
	}
}
