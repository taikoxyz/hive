package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/libp2p/go-libp2p/core/peer"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

type PreconfTestSpec struct {
	suite_base.BaseTestSpec

	anchorL1Head *types.Header
}

func (r *PreconfTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.BaseTestSpec.GetTestnetConfig()
	preConfig := cfg.CreateConfig

	var indexd = make(map[int]bool)
	cfg.CreateConfig = func(index int, nodes clients.Nodes) (hivesim.Params, error) {
		if indexd[index] {
			return params.ClusterEnvs[index], nil
		}
		indexd[index] = true

		envs, err := preConfig(index, nodes)
		if err != nil {
			return nil, err
		}

		envs["PRECONFIRMATION_SERVER_PORT"] = fmt.Sprintf("%d", clients.PreconfServerPort)
		envs["PRECONFIRMATION_SERVER_SIGNATURE_CHECK"] = "true"
		envs["PRECONFIRMATION_P2P_DISCOVERY_PATH"] = "memory"
		envs["PRECONFIRMATION_P2P_PEERSTORE_PATH"] = "memory"
		envs["PRECONFIRMATION_P2P_SEQUENCER_KEY"] = envs["L1_PROPOSER_PRIV_KEY"]
		envs["PRECONFIRMATION_P2P_PRIV_RAW"] = params.ChainAuths[index].Key

		var staticPeers string
		for i := 0; i < index; i++ {
			sk, err := parsePriv(params.ChainAuths[i].Key)
			if err != nil {
				return nil, err
			}

			idB, err := peer.IDFromPublicKey(sk.GetPublic())
			if err != nil {
				return nil, err
			}

			staticPeers = fmt.Sprintf("%s,/ip4/%s/tcp/%d/p2p/%s", staticPeers, nodes[i].DriverClient.NetworkIP(), clients.PreconfP2pPort, idB)
		}
		envs["PRECONFIRMATION_P2P_STATIC"] = staticPeers

		params.ClusterEnvs[index] = envs

		return envs, nil
	}

	return cfg
}

func (r *PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	// Start all the cluster's nodes.
	for _, node := range testnet.Nodes {
		t.Nil(node.Start(), "cannot start L2EthClient")
	}

	var (
		node     = testnet.Nodes[0]
		driver   = node.DriverClient
		proposer = node.ProposerClient
		l2geth   = node.L2EthClient
	)

	t.Nil(l2geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight))

	// stop the proposer.
	proposer.PauseClient()

	// For Debug
	if r.IsDebug() {
		driver.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		l2Header, anchorL1Header, err := preconferBlock(index, driver.Client, driver.PreconfServerURL(), 5)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// Verify latest preconf block.
		verifyL2Chain(t, testnet.Nodes, l2Header)

		// propose txs.
		_, _, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
		t.FailIfNotNil(err, "cannot propose txs")

		// todo: check l1Origin and head l1Origin, brefore and after
		// todo: multi times reorg and preconf
		// Verify latest propose block.
		verifyL2Chain(t, testnet.Nodes, l2Header)
	}
}

func verifyL2Chain(t *hivesim.T, nodes []*clients.Node, l2Header *types.Header) {
	for _, node := range nodes {
		l2geth := node.L2EthClient
		t.FailIfNotNil(l2geth.WaitLatestNumber(context.Background(), time.Second*30, l2Header.Number.Uint64()), "cannot get latest number")

		actualHeader, err := l2geth.EthClient.HeaderByNumber(context.Background(), l2Header.Number)
		t.FailIfNotNil(err, "cannot get header by number")
		t.Equal(l2Header.Hash().String(), actualHeader.Hash().String(), fmt.Sprintf("header hash not equal, node id: %d", node.Index))
	}
}

func verifyL1Origin(t *hivesim.T) {}
