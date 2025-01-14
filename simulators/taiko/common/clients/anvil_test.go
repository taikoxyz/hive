package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"taiko/params"
	"testing"
)

func TestAnvil(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)
	anvil := &AnvilClient{
		SecondsPerSlot: 3,
		reorgCache:     make(map[uint64]*L1BlockInfo),
		EthNode: &EthNode{
			EthClient: client,
		},
	}

	ctx := context.Background()

	anvil.StopMining()
	for _, tx := range params.ContractTxs {
		assert.NoError(t, anvil.EthClient.SendTransaction(ctx, tx))
	}
	anvil.MineBlock()

	anvil.StartMining()
}
