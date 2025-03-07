package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/hive/hivesim"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/params"
	"testing"
	"time"
)

var (
	anvil *AnvilClient
)

func init() {
	var err error
	rpccli, err = rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:        "ws://localhost:8545",
		L2Endpoint:        "ws://localhost:6046",
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:  "http://localhost:6051",
		JwtSecret:         "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
	})
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
			EthClient: rpccli.L1,
		},
	}
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
