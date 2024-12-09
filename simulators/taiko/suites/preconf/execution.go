package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
	"math/big"
	"math/rand/v2"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

type PreconfTestSpec struct {
	suite_base.BaseTestSpec

	batchID uint64
	l1Head  *types.Header
}

func (r *PreconfTestSpec) GetTestnetConfig() *testnet.Config {
	params.SetEnvParams("SOFT_BLOCK_SERVER_PORT", fmt.Sprintf("%d", clients.SoftBlockServerPort))
	return r.BaseTestSpec.GetTestnetConfig()
}

func (r *PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	t.Nil(node.Start(), "cannot start node")

	driver := node.DriverClient
	taikoL1 := driver.TaikoL1

	// wait l2 node.
	safeNum := uint64(5)
	eventCh := make(chan *bindings.TaikoL1ClientBlockProposedV2, 3)
	sub, err := taikoL1.WatchBlockProposedV2(nil, eventCh, nil)
	t.Nil(err)
	defer sub.Unsubscribe()
	for event := range eventCh {
		if event.BlockId.Uint64() == safeNum {
			// pause proposer
			node.ProposerClient.PauseClient()
			break
		}
	}
	t.Nil(node.L2EthClient.WaitTargetNumber(ctx, time.Second*60, safeNum))

	for loops := 0; loops < rand.IntN(10)+10; loops++ {
		for i := 0; i < 4; i++ {
			time.Sleep(time.Second)
			if i%2 == 0 {
				r.insertNewSoftBlock(t, node.L2EthClient, driver)
			} else {
				r.replaceLatestSoftBlock(t, driver)
			}
		}

		r.randomRemoveSoftBlock(t, driver)

		// propose safe block.
		r.proposeTxListsForSoftBlocks(t, node.L2EthClient, node.ProposerClient)

		// Propose all the pending txs.
		r.proposeTxListForPendingTxs(t, node.L2EthClient, node.ProposerClient)
	}
}

func (r *PreconfTestSpec) insertNewSoftBlock(
	t *hivesim.T,
	l2EthClient *clients.TaikoGethClient,
	driver *clients.DriverClient,
) {
	r.batchID = 0
	t.Logf("start insert new soft block")
	defer t.Logf("successfully insert new soft block")

	l2Block, err := driver.L2.BlockByNumber(context.Background(), nil)
	t.Nil(err, "l2 latest block")

	l2Num := l2Block.NumberU64()
	l1Head, txs, err := driver.BuildSoftBlock(
		l2Num+1,
		r.batchID,
		false,
		false,
		nil,
	)
	t.Nil(err)
	r.l1Head = l1Head

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
	t.Equal(r.batchID, l1Origin.BatchID.Uint64(), "l1Origin's batchID")
	t.Equal(false, l1Origin.EndOfBlock, "l1Origin's endOfBlock")
	t.Equal(false, l1Origin.EndOfPreconf, "l1Origin's endOfPreconf")
	t.Nil(driver.StateError(), "state error")
}

func (r *PreconfTestSpec) replaceLatestSoftBlock(
	t *hivesim.T,
	driver *clients.DriverClient,
) {
	r.batchID++
	t.Logf("start replace latest soft block")
	defer t.Logf("successfully replace latest soft block")

	preHead := driver.L2Head.Load()
	preBlock, err := driver.L2.BlockByNumber(context.Background(), preHead.Number)
	t.Nil(err, "l2 latest block before append soft block")

	_, txs, err := driver.BuildSoftBlock(
		preBlock.NumberU64(),
		r.batchID,
		false,
		false,
		r.l1Head,
	)
	t.Nil(err)

	curBlock, err := driver.L2.BlockByNumber(context.Background(), nil)
	t.Nil(err, "l2 latest block after append soft block")

	t.Equal(preBlock.NumberU64(), curBlock.NumberU64(), "block number")
	t.Equal(preBlock.Transactions().Len()+txs.Len(), curBlock.Transactions().Len(), "transaction count")

	l1Origin, err := driver.L2.L1OriginByID(context.Background(), curBlock.Number())
	t.Nil(err, "l1 origin")

	t.Equal(r.batchID, l1Origin.BatchID.Uint64(), "batchID")
	t.Equal(curBlock.Number().Uint64(), l1Origin.BlockID.Uint64(), "l1Origin's blockID")
	t.Equal(curBlock.Hash().String(), l1Origin.L2BlockHash.String(), "l1Origin's l2BlockHash")
	t.Equal(uint64(0), l1Origin.L1BlockHeight.Uint64(), "l1Origin's l1BlockHeight")
	t.Equal(common.Hash{}.String(), l1Origin.L1BlockHash.String(), "l1Origin's l1BlockHash")
	t.Equal(false, l1Origin.EndOfBlock, "l1Origin's endOfBlock")
	t.Equal(false, l1Origin.EndOfPreconf, "l1Origin's endOfPreconf")
	t.Nil(driver.StateError(), "state error")
}

func (r *PreconfTestSpec) randomRemoveSoftBlock(
	t *hivesim.T,
	driver *clients.DriverClient,
) {
	t.Logf("start remove soft blocks")
	defer t.Logf("successfully remove soft blocks")

	canonicalL1Origin := driver.CanonicalL1Origin.Load()
	latestNum, err := driver.L2.BlockNumber(context.Background())
	t.Nil(err, "l2 latest block number")

	l2Num := rand.Uint64N(latestNum-canonicalL1Origin.BlockID.Uint64()) + canonicalL1Origin.BlockID.Uint64()

	t.Logf("canonical number: %d, remove number: %d, latest number: %d", canonicalL1Origin.BlockID.Uint64(), l2Num, latestNum)

	curL1Origin, err := driver.L2.L1OriginByID(context.Background(), big.NewInt(int64(l2Num)))
	t.Nil(err, "l1 origin")

	// remove soft blocks.
	t.Nil(driver.RemoveSoftBlocks(l2Num))
	time.Sleep(time.Millisecond * 500)

	l2Head := driver.L2Head.Load()
	// wait a new current l1Origin
	curL1Origin2 := driver.LatestL1Origin.Load()

	t.Equal(l2Num, l2Head.Number.Uint64(), "l2 number")
	t.Equal(curL1Origin2.BlockID, l2Head.Number, "current l1Origin's blockID")

	t.Equal(curL1Origin.BatchID, curL1Origin2.BatchID, "l1Origin's batchID")
	t.Equal(curL1Origin.L1BlockHeight.Uint64(), curL1Origin2.L1BlockHeight.Uint64(), "l1Origin's l1BlockHeight")
	t.Equal(curL1Origin.L1BlockHash.String(), curL1Origin2.L1BlockHash.String(), "l1Origin's l1BlockHash")
	t.Equal(curL1Origin.L2BlockHash.String(), curL1Origin2.L2BlockHash.String(), "l1Origin's l2BlockHash")
	t.Equal(curL1Origin.EndOfBlock, curL1Origin2.EndOfBlock, "l1Origin's endOfBlock")
	t.Equal(curL1Origin.EndOfPreconf, curL1Origin2.EndOfPreconf, "l1Origin's endOfPreconf")
	t.Equal(curL1Origin.Preconfer.String(), curL1Origin2.Preconfer.String(), "l1Origin's preconfer")
	t.Nil(driver.StateError(), "state error")
}

func (r *PreconfTestSpec) proposeTxListsForSoftBlocks(
	t *hivesim.T,
	l2EthClient *clients.TaikoGethClient,
	proposerClient *clients.ProposerClient,
) {
	t.Logf("start propose tx lists")
	defer t.Logf("successfully start propose tx lists")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	canonicalL1Origin, txs, err := proposerClient.ProposeTxLists(ctx, t)
	t.Nil(err)

	// No soft blocks need to be proposed
	if canonicalL1Origin == nil {
		return
	}

	err = l2EthClient.WaitTargetNumber(ctx, time.Second*60, canonicalL1Origin.BlockID.Uint64()+uint64(len(txs)))
	t.Nil(err)

	var (
		l1Number    *big.Int
		l1Hash      common.Hash
		startNumber = canonicalL1Origin.BlockID.Uint64() + 1
	)
	for _, txLst := range txs {
		block, err := proposerClient.L2.BlockByNumber(ctx, big.NewInt(int64(startNumber)))
		t.Nil(err)

		l1Origin, err := proposerClient.L2.L1OriginByID(ctx, big.NewInt(int64(startNumber)))
		t.Nil(err)

		t.Equal(txLst.Len(), block.Transactions().Len())
		t.Equal(block.NumberU64(), l1Origin.BlockID.Uint64(), "l1Origin's blockID")
		t.Equal(block.Hash().String(), l1Origin.L2BlockHash.String(), "l1Origin's l2BlockHash")
		t.True(l1Origin.BatchID == nil, "l1Origin's batchID")
		t.True(l1Origin.L1BlockHeight != nil, "l1Origin's l1BlockHeight")
		if l1Number == nil {
			l1Number = l1Origin.L1BlockHeight
			l1Hash = l1Origin.L1BlockHash
		} else {
			t.Equal(l1Number.Uint64(), l1Origin.L1BlockHeight.Uint64(), "l1Origin's l1BlockHeight")
			t.Equal(l1Hash.String(), l1Origin.L1BlockHash.String(), "l1Origin's l1BlockHash")
		}
	}
}

func (r *PreconfTestSpec) proposeTxListForPendingTxs(
	t *hivesim.T,
	l2EthClient *clients.TaikoGethClient,
	proposerClient *clients.ProposerClient,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// propose pending txs.
	t.Nil(proposerClient.ProposeOp(ctx), "propose for pending txs")

	// Wait until all the txs are proposed.

	for {
		timeAfter := time.After(time.Second * 10)
		select {
		case <-timeAfter:
			return
		case number := <-proposerClient.ProposedBlockID:
			_ = l2EthClient.WaitTargetNumber(ctx, time.Second*10, number.Uint64())
		}
	}
}
