package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"taiko/params"
	"testing"
)

func TestAnvil(t *testing.T) {
	rpcCli, err := rpc.Dial("http://localhost:8545")
	assert.NoError(t, err)
	anvil := &AnvilClient{
		SecondsPerSlot: 3,
		reorgCache:     make(map[string]*L1BlockInfo),
		EthNode: &EthNode{
			client: rpcCli,
		},
	}

	ctx := context.Background()
	client := anvil.EthClient()

	anvil.StopMining()
	for _, tx := range params.ContractTxs {
		assert.NoError(t, client.SendTransaction(ctx, tx))
	}
	anvil.MineBlock()

	anvil.StartMining()
}
