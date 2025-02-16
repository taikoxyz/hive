package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
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
	// set preconf environment variables.
	params.SetEnvParams("PRECONFIRMATION_SERVER_PORT", fmt.Sprintf("%d", clients.PreconfServerPort))
	params.SetEnvParams("PRECONFIRMATION_SERVER_SIGNATURE_CHECK", "true")
	return r.BaseTestSpec.GetTestnetConfig()
}

func (r *PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	t.Nil(node.Start(), "cannot start node")

	// For Debug
	if r.IsDebug() {
		time.Sleep(time.Minute * 120)
	}

	driver := node.DriverClient
	//pacayaTokens, err := pacaya.NewPacayaClients(driver.L1, driver.L2)

	// Wait until l2 height touched pacaya fork.
	for driver.L2Head.Load().Number.Uint64() < driver.PacayaClients.ForkHeight {
		time.Sleep(1 * time.Second)
	}
}

func (r *PreconfTestSpec) buildPreconfBlock(
	t *hivesim.T,
	l2EthClient *clients.TaikoGethClient,
	driver *clients.DriverClient,
) {
	t.Logf("start insert new soft block")
	defer t.Logf("successfully insert new soft block")

	l2Block, err := driver.L2.BlockByNumber(context.Background(), nil)
	t.Nil(err, "l2 latest block")

	l2Num := l2Block.NumberU64()
	l1Head, txs, err := driver.BuildPreconfBlock(
		l2Num+1,
		nil,
	)
	t.Nil(err)
	r.anchorL1Head = l1Head

	// wait l2 node.
	t.Nil(l2EthClient.WaitTargetNumber(context.Background(), time.Second*60, l2Num+1))

	l2Block, err = driver.L2.BlockByNumber(context.Background(), nil)
	t.Nil(err, "l2 latest block")

	t.Equal(l2Num+1, l2Block.NumberU64(), "block number")

	// check txs count.
	t.Equal(txs.Len()+1, l2Block.Transactions().Len(), "transaction count")

	// check l1Origin variables.
	l1Origin, err := driver.L2.L1OriginByID(context.Background(), l2Block.Number())
	t.Nil(err, "l1 origin")

	t.Equal(l2Block.Number(), l1Origin.BlockID, "l1Origin's blockID")
	t.Equal(l2Block.Hash().String(), l1Origin.L2BlockHash.String(), "l1Origin's l2BlockHash")
	t.Equal(uint64(0), l1Origin.L1BlockHeight.Uint64(), "l1Origin's l1BlockHeight")
	t.Equal(common.Hash{}.String(), l1Origin.L1BlockHash.String(), "l1Origin's l1BlockHash")
}
