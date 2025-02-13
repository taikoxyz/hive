package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/stretchr/testify/assert"
	"taiko/params"
	"testing"
	"time"
)

var (
	anvil *AnvilClient
)

func init() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}
	anvil = &AnvilClient{
		SecondsPerSlot: 3,
		reorgCache:     make(map[uint64]*L1BlockInfo),
		EthNode: &EthNode{
			HiveManagedClient: &HiveManagedClient{
				T: &hivesim.T{},
				HiveClientDefinition: &hivesim.ClientDefinition{
					Name: "anvil",
				},
			},
			EthClient: client,
		},
	}
	_ = anvil.InitL1Clients()
}

func TestAnvil(t *testing.T) {

	ctx := context.Background()

	anvil.StopMining()
	for _, tx := range params.ContractTxs {
		assert.NoError(t, anvil.EthClient.SendTransaction(ctx, tx))
	}
	anvil.MineBlock()

	anvil.StartMining()
}

func TestVerified(t *testing.T) {
	assert.NoError(t, anvil.WaitLatestVerifiedNumber(context.Background(), time.Second*60, 3))
}
