package preconf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, CrossTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name: "cross",
			},
		},
	})
}

//Test two rollback preconf chains,they are 5 blocks and 6 blocks.
//
//start two cluster(anvil, l2geth, driver, proposer, prover):
//
//- Mint several blocks until over the pacaya fork.
//- Mint 5 preconf blocks and collect preconf requestbodies and revert l2geth.
//- Mint 6 preconf blocks and collect preconf requestbodies and revert l2geth.
//- create 2 p2pNodes(connect to the driver p2p node) and send requestbodies in reversed order.
//- Wait and check all the reverted preconf blocks are created again.

type CrossTestSpec struct {
	PreconfTestSpec
}

func (r CrossTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_preconf_cross"

	return cfg
}

func (r CrossTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node     = testnet.Nodes[0]
		anvil    = node.AnvilClient
		l2geth   = node.L2EthClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
		prover   = node.ProverClient
	)

	// Start all the cluster's nodes.
	for _, node := range testnet.Nodes {
		t.Nil(node.Start(), "cannot start L2EthClient")
	}

	for range time.Tick(time.Second) {
		prover.VerifyBlocks(params.L1Auths[0])
		lastVerifiedBlockID := prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID > 0 {
			prover.Shutdown()
			break
		}
	}

	l2geth.WaitLatestNumber(ctx, time.Minute*3, proposer.PacayaClients.ForkHeight-1)
	// stop the proposer.
	t.FailIfNotNil(proposer.Shutdown())

	var (
		batchSize = rand.IntN(50-10) + 10
		l1Number  = anvil.BlockNumber(ctx)
	)

	sendBodies0, l2Header0 := getPreconfBodies(ctx, t, node, l1Number.Sub(l1Number, big.NewInt(1)), batchSize)
	sendBodies1, l2Header1 := getPreconfBodies(ctx, t, node, l1Number, batchSize+5)

	slices.Reverse(sendBodies0)
	slices.Reverse(sendBodies1)

	data0, _ := json.Marshal(sendBodies0)
	data1, _ := json.Marshal(sendBodies1)
	t.Logf("the latest l2chain header, number: %d, hash: %s, preconf_bodies: %s", l2Header0.Number.Uint64(), l2Header0.Hash(), string(data0))
	t.Logf("the latest l2chain header, number: %d, hash: %s, preconf_bodies: %s", l2Header1.Number.Uint64(), l2Header1.Hash(), string(data1))

	// For DevDebug
	if testnet.DevDebug {
		driver.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	p2pNode, err := clients.NewP2PNode(ctx, driver.Client, driver.GetPreconfP2PNode())
	t.FailIfNotNil(err, fmt.Sprintf("cannot create p2p node"))
	defer p2pNode.Close()

	// Waiting for p2p node is connected.
	t.FailIfNotNil(p2pNode.WaitConnected(ctx, time.Minute))

	slices.Reverse(sendBodies0)

	var eg errgroup.Group

	eg.Go(func() error {
		driver := testnet.Nodes[0].DriverClient
		p2pNode, err := clients.NewP2PNode(ctx, driver.Client, driver.GetPreconfP2PNode())
		t.FailIfNotNil(err, fmt.Sprintf("cannot create p2p node"))
		defer p2pNode.Close()

		// Waiting for p2p node is connected.
		t.FailIfNotNil(p2pNode.WaitConnected(ctx, time.Minute))

		for _, requestBody := range sendBodies0 {
			t.FailIfNotNil(p2pNode.PublishL2Payload(ctx, requestBody))
		}
		return nil
	})

	eg.Go(func() error {
		driver := testnet.Nodes[1].DriverClient
		p2pNode, err := clients.NewP2PNode(ctx, driver.Client, driver.GetPreconfP2PNode())
		t.FailIfNotNil(err, fmt.Sprintf("cannot create p2p node"))
		defer p2pNode.Close()

		// Waiting for p2p node is connected.
		t.FailIfNotNil(p2pNode.WaitConnected(ctx, time.Minute))

		slices.Reverse(sendBodies1)
		for _, requestBody := range sendBodies1 {
			t.FailIfNotNil(p2pNode.PublishL2Payload(ctx, requestBody))
		}
		return nil
	})

	t.FailIfNotNil(eg.Wait())

	// Verify all the clusters.
	for _, nd := range testnet.Nodes {
		l2geth := nd.L2EthClient
		l2geth.WaitLatestNumber(ctx, time.Minute*3, l2Header1.Number.Uint64())

		actualL2Header := l2geth.HeaderByNumber(ctx, nil)
		t.Equal(l2Header1.Hash().String(), actualL2Header.Hash().String())
	}
}
