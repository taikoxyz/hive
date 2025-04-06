package preconf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/hive/hivesim"
	"github.com/holiman/uint256"
	"golang.org/x/exp/slices"
	"math/rand/v2"
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
	cfg.Network = "network_preconf_ancients_12"

	return cfg
}

func (r AncientsTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
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

	p2pNode, err := clients.NewP2PNode(ctx, driver.Client, driver.Index, driver.GetPreconfP2PNode())
	t.FailIfNotNil(err, fmt.Sprintf("cannot create p2p node"))
	defer p2pNode.Close()

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

	// Create a batch of preconf blocks.
	var (
		batchSize  = rand.IntN(50-10) + 10
		l2Number   = l2geth.BlockNumber(ctx)
		sendBodies []*eth.ExecutionPayloadEnvelope
	)

	for i := 0; i < batchSize; i++ {
		requestBody, err := clients.BuildPreconfRequestBody(ctx, driver.Client, params.ChainAuths[1].PrivateKey, anvil.BlockNumber(ctx), nil)
		t.FailIfNotNil(err, "cannot build preconf request body")

		header, err := clients.SendPreconfBlock(driver.PreconfServerURL(), requestBody)
		t.FailIfNotNil(err, "cannot send preconf request")

		sendBody := &eth.ExecutionPayloadEnvelope{
			ExecutionPayload: &eth.ExecutionPayload{
				ParentHash:    header.ParentHash,
				FeeRecipient:  header.Coinbase,
				ExtraData:     header.Extra,
				PrevRandao:    eth.Bytes32(header.MixDigest),
				BlockNumber:   eth.Uint64Quantity(header.Number.Uint64()),
				GasLimit:      eth.Uint64Quantity(header.GasLimit),
				GasUsed:       eth.Uint64Quantity(header.GasUsed),
				Timestamp:     eth.Uint64Quantity(header.Time),
				BlockHash:     header.Hash(),
				BaseFeePerGas: eth.Uint256Quantity(*uint256.NewInt(requestBody.ExecutableData.BaseFeePerGas)),
				Transactions:  []eth.Data{requestBody.ExecutableData.Transactions},
			},
		}
		sendBodies = append(sendBodies, sendBody)

		body, _ := json.Marshal(sendBody)
		t.Logf("%s", string(body))

		time.Sleep(time.Second)
	}

	l2geth.RevertTaikoGeth(ctx, l2Number)

	// Waiting for p2p node is connected.
	t.FailIfNotNil(p2pNode.WaitConnected(ctx, time.Minute))
	time.Sleep(time.Minute)

	slices.Reverse(sendBodies)

	for _, requestBody := range sendBodies {
		t.FailIfNotNil(p2pNode.PublishL2Payload(ctx, requestBody))
		time.Sleep(time.Second)
	}

	// For DevDebug
	if testnet.DevDebug {
		driver.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	l2geth.WaitLatestNumber(ctx, time.Minute*3, l2Number.Uint64()+uint64(batchSize))
}
