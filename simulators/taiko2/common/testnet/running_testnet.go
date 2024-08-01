package testnet

import (
	"context"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"net"
	"taiko2/common/clients"
	consensus_config "taiko2/common/config/consensus"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/hive/hivesim"
	exec_client "github.com/marioevz/eth-clients/clients/execution"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	execution_config "taiko2/common/config/execution"
	"taiko2/common/utils"
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
	prep, err := PrepareTestnet(env, config, generateState)
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

			executionDef = env.Clients.ClientByNameAndRole(
				node.ExecutionClientName(),
				"eth1",
			)
			beaconDef = env.Clients.ClientByNameAndRole(
				node.ConsensusClientName(),
				"beacon",
			)
			validatorDef = env.Clients.ClientByNameAndRole(
				node.ValidatorClientName(),
				"validator",
			)
			executionTTD = int64(0)
			beaconTTD    = int64(0)
		)

		if executionDef == nil || beaconDef == nil || validatorDef == nil {
			t.Fatalf("FAIL: Unable to get client")
		}
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
		nodeClient.ExecutionClient = prep.prepareExecutionNode(
			parentCtx,
			testnet,
			executionDef,
			config.Eth1Consensus,
			node.Chain,
			clients.ExecutionClientConfig{
				ClientIndex:             nodeIndex,
				TerminalTotalDifficulty: executionTTD,
				Subnet:                  node.GetExecutionSubnet(),
				JWTSecret:               ethcommon.FromHex(config.JWTSecret),
				ProxyConfig: &clients.ExecutionProxyConfig{
					Host:                   simulatorIP,
					Port:                   exec_client.PortEngineRPC + nodeIndex,
					TrackForkchoiceUpdated: false,
					LogEngineCalls:         env.LogEngineCalls,
				},
			},
		)

		if node.ConsensusClient != "" {
			nodeClient.BeaconClient = prep.prepareBeaconNode(
				parentCtx,
				testnet,
				beaconDef,
				&clients.BeaconClientConfig{
					ClientIndex:             nodeIndex,
					BeaconAPIPort:           clients.PortBeaconAPI,
					TerminalTotalDifficulty: beaconTTD,
					Spec:                    testnet.spec,
					GenesisValidatorsRoot:   &testnet.genesisValidatorsRoot,
					GenesisTime:             &testnet.genesisTime,
					Subnet:                  node.GetConsensusSubnet(),
				},
				nodeClient.ExecutionClient,
			)

			nodeClient.ValidatorClient = prep.prepareValidatorClient(
				parentCtx,
				testnet,
				validatorDef,
				nodeClient.BeaconClient,
				nodeIndex,
			)
		}

		// Add rest of properties
		nodeClient.Logging = t
		nodeClient.Index = nodeIndex
		nodeClient.Verification = node.TestVerificationNode
		// Start the node clients if specified so
		if !node.DisableStartup {
			t.Logf("Starting node %d", nodeIndex)
			if err := nodeClient.Start(); err != nil {
				t.Fatalf("FAIL: Unable to start node %d: %v", nodeIndex, err)
			}
			// Connect to the network if specified
			if err := nodeClient.CreateOrConnectNetwork(t, fmt.Sprintf("%s_%d", config.Network, nodeIndex)); err != nil {
				t.Fatalf("FAIL: Unable to connect to network: %v", err)
			}
		} else {
			t.Logf("Node %d startup disabled, skipping", nodeIndex)
		}
	}

	return testnet
}

func (t *Testnet) Stop() {
	for _, p := range t.Proxies().Running() {
		p.Cancel()
	}
}

func (t *Testnet) ValidatorClientIndex(pk [48]byte) (int, error) {
	for i, v := range t.ValidatorClients() {
		if v.ContainsKey(pk) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("key not found in any validator client")
}
