package testnet

import (
	"github.com/ethereum/go-ethereum/common"
	"taiko/common/clients"
	"time"

	"github.com/ethereum/hive/hivesim"
)

type Testnet struct {
	*hivesim.T
	clients.Nodes

	// debug flag
	DevDebug  bool
	DevExpose bool
}

func (t *Testnet) GenesisTimeUnix() time.Time {
	return time.Unix(1742366600, 0)
}

func (t *Testnet) GenesisValidatorsRoot() [32]byte {
	return common.Hash{}
}

func StartTestnet(
	t *hivesim.T,
	clientGroups clients.ClientGroups,
	config *Config,
) *Testnet {
	prep, err := PrepareTestnet(config)
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
		prep.prepareGethNode(index, testnet, config, executionDef)
		prep.prepareBeaconNode(index, testnet, config, clientsByRole[clients.Beacon])
		prep.prepareValidatorClient(testnet, clientsByRole[clients.Validator])
		prep.prepareBlobScanClient(index, testnet, config, clientsByRole)

		// taiko-geth or taiko-reth
		taikoGeth := clientsByRole[clients.Eth2]
		if taikoGeth == nil {
			taikoGeth = clientsByRole[clients.Reth]
		}
		prep.prepareTaikoGethClient(index, testnet, config, taikoGeth)
		prep.prepareDriverClient(index, testnet, config, clientsByRole[clients.Driver])
		prep.prepareProposerClient(index, testnet, config, clientsByRole[clients.Proposer])
		prep.prepareProverClient(index, testnet, config, clientsByRole[clients.Prover])

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
