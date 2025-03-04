package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

func init() {
	Tests = append(Tests,
		&PreconfTestSpec{
			BaseTestSpec: suite_base.BaseTestSpec{
				Name:  "preconf",
				Debug: false,
			},
		},
	)
}

type PreconfTestSpec struct {
	suite_base.BaseTestSpec
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

	l2geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight)

	// stop the proposer.
	proposer.PauseClient()

	// For Debug
	if r.Debug {
		driver.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	for times := 0; times < 4; times++ {
		index := times % len(testnet.Nodes)
		driver = testnet.Nodes[index].DriverClient

		l2Header, anchorL1Header, _, err := preconferBlock(index, driver.Client, driver.PreconfServerURL(), 5)
		t.FailIfNotNil(err, "cannot preconfirmer proposer")

		// Verify latest preconf block.
		verifyL2Chain(t, true, testnet.Nodes, l2Header)

		// propose txs.
		_, err = proposeBlock(ctx, driver.Envs, driver.Client, anchorL1Header)
		t.FailIfNotNil(err, "cannot propose txs")

		// Verify latest propose block.
		verifyL2Chain(t, false, testnet.Nodes, l2Header)
	}
}

func verifyL2Chain(t *hivesim.T, isPreconf bool, nodes []*clients.Node, l2Header *types.Header) {
	l2cli := nodes[0].DriverClient.Client.L2
	l1Origin, err := l2cli.L1OriginByID(context.Background(), l2Header.Number)
	t.FailIfNotNil(err, "cannot get l1 origin by id")

	l1HeadOrigin, err := l2cli.HeadL1Origin(context.Background())
	t.FailIfNotNil(err, "cannot get head l1 origin")

	for _, node := range nodes {
		l2geth := node.L2EthClient
		l2geth.WaitLatestNumber(context.Background(), time.Second*30, l2Header.Number.Uint64())

		actualHeader, err := l2geth.EthClient.HeaderByNumber(context.Background(), l2Header.Number)
		t.FailIfNotNil(err, "cannot get header by number")
		t.Equal(l2Header.Hash().String(), actualHeader.Hash().String(), fmt.Sprintf("header hash not equal, node id: %d", node.Index))

		l2cli = node.DriverClient.L2
		l1o, err := l2cli.L1OriginByID(context.Background(), l2Header.Number)
		t.FailIfNotNil(err, "cannot get l1 origin by id")

		l1ho, err := l2cli.HeadL1Origin(context.Background())
		t.FailIfNotNil(err, "cannot get head l1 origin")

		if isPreconf {
			t.Equal(true, l1Origin.L1BlockHeight == nil)
			t.Equal(l1Origin.L1BlockHash.String(), common.Hash{}.String())
		} else {
			t.Equal(l1Origin.L1BlockHeight.Uint64(), l1o.L1BlockHeight.Uint64())
			t.Equal(l1Origin.L1BlockHash.String(), l1o.L1BlockHash.String())
		}
		t.Equal(l1Origin.BlockID.Uint64(), l1o.BlockID.Uint64())
		t.Equal(l1Origin.L2BlockHash.String(), l1o.L2BlockHash.String())

		t.Equal(l1HeadOrigin.L1BlockHeight.Uint64(), l1ho.L1BlockHeight.Uint64())
		t.Equal(l1HeadOrigin.L1BlockHash.String(), l1ho.L1BlockHash.String())
		t.Equal(l1HeadOrigin.BlockID.Uint64(), l1ho.BlockID.Uint64())
		t.Equal(l1HeadOrigin.L2BlockHash.String(), l1ho.L2BlockHash.String())
	}
}
