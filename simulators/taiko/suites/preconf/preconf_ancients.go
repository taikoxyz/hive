package preconf

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"math/big"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, AncientsTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name: "ancients",
			},
		},
	})
}

type AncientsTestSpec struct {
	PreconfTestSpec
}

func (r AncientsTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := r.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_preconf_ancients"

	return cfg
}

func (r AncientsTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node     = testnet.Nodes[0]
		anvil    = node.AnvilClient
		l2geth   = node.L2EthClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
	)

	// Start all the cluster's nodes.
	for _, node := range testnet.Nodes {
		t.Nil(node.Start(), "cannot start L2EthClient")
	}

	l2geth.WaitLatestNumber(ctx, time.Minute*3, proposer.PacayaClients.ForkHeight-1)

	// stop the proposer.
	t.FailIfNotNil(proposer.Shutdown())

	// Create a batch of preconf request bodies.
	var (
		batchSize     = 5
		l1Number      = anvil.BlockNumber(ctx)
		l2Number      = l2geth.BlockNumber(ctx)
		requestBodies []*preconfblocks.BuildPreconfBlockRequestBody

		//l2Header *types.Header
	)

	// Create a batch of preconf blocks.
	for index := 0; index < batchSize; index++ {
		l2Num := big.NewInt(l2Number.Int64() + int64(index) + 1)
		requestBody, err := clients.BuildPreconfRequestBody(ctx, driver.Client, params.ChainAuths[index*2+1].PrivateKey, l1Number, l2Num)
		t.FailIfNotNil(err, "cannot build preconf request body")

		//l2Header, err = clients.BuildPreconfBlock(driver.PreconfServerURL(), requestBody)
		//t.FailIfNotNil(err, "cannot build preconf block")

		time.Sleep(time.Second)

		driver.PublishL2Payload(ctx, requestBody)

		requestBodies = append(requestBodies, requestBody)
	}

	// For DevDebug
	if testnet.DevDebug {
		time.Sleep(time.Hour * 2)
	}
	//l2geth.RevertTaikoGeth()
}
