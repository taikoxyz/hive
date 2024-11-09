package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"math/rand/v2"
	"strings"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

type PreconfTestSpec struct {
	suite_base.BaseTestSpec
}

var (
	l1HeadCh            = make(chan *types.Header, 1)
	l2BlockCh           = make(chan *types.Block, 1)
	curL1OriginCh       = make(chan *rawdb.L1Origin, 1)
	canonicalL1OriginCh = make(chan *rawdb.L1Origin, 1)
)

func (r PreconfTestSpec) GetTestnetConfig() *testnet.Config {
	params.SetEnvParams("SOFT_BLOCK_SERVER_PORT", fmt.Sprintf("%d", clients.SoftBlockServerPort))
	return r.BaseTestSpec.GetTestnetConfig()
}

func (r PreconfTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	node := testnet.Nodes[0]
	if err := node.Start(); err != nil {
		t.Fatalf("cannot start node: %v", err)
	}

	driver := node.DriverClient
	l1Cli, err := ethclient.Dial(node.AnvilClient.WSURL())
	if err != nil {
		t.Fatalf("cannot connect l1 client: %v", err)
	}
	l2Cli, err := ethclient.Dial(node.L2EthClient.WSURL())
	if err != nil {
		t.Fatalf("cannot dial l2 client: %v", err)
	}

	l1HCh := make(chan *types.Header, 2)
	l2HeadCh := make(chan *types.Header, 2)
	l1Sub, err := l1Cli.SubscribeNewHead(context.Background(), l1HeadCh)
	t.Nil(err)
	defer l1Sub.Unsubscribe()
	l2Sub, err := l2Cli.SubscribeNewHead(context.Background(), l2HeadCh)
	t.Nil(err)
	defer l2Sub.Unsubscribe()

	var closeCh = make(chan struct{})
	go func() {
		for {
			select {
			case <-closeCh:
				return
			case head := <-l1HCh:
				select {
				case <-l1HeadCh:
					l1HeadCh <- head
				default:
					l1HeadCh <- head
				}
			case head := <-l2HeadCh:
				block, err := l2Cli.BlockByHash(ctx, head.Hash())
				t.Nil(err)
				select {
				case <-l2BlockCh:
					l2BlockCh <- block
				default:
					l2BlockCh <- block
				}
				origin, err := l2Cli.HeadL1Origin(ctx)
				t.Nil(err)
				select {
				case <-canonicalL1OriginCh:
					canonicalL1OriginCh <- origin
				default:
					canonicalL1OriginCh <- origin
				}
				origin, err = l2Cli.L1OriginByID(ctx, head.Number)
				t.Nil(err)
				select {
				case <-curL1OriginCh:
					curL1OriginCh <- origin
				default:
					curL1OriginCh <- origin
				}
			}
		}
	}()
	defer close(closeCh)

	// wait l2 node.
	t.Nil(node.L2EthClient.WaitTargetNumber(ctx, time.Second*60, 5))

	var batchID uint64
	for range 10 {
		time.Sleep(time.Millisecond * 200)
		if rand.Int()/2 == 0 {
			batchID = 0
			r.insertNewSoftBlock(t, l1Cli, l2Cli, driver)
		} else {
			batchID++
			r.replaceLatestSoftBlock(t, l1Cli, l2Cli, driver, batchID)
		}
	}

}

func (r PreconfTestSpec) insertNewSoftBlock(
	t *hivesim.T,
	l1cli *ethclient.Client,
	l2cli *ethclient.Client,
	driver *clients.DriverClient,
) {
	l2Block := <-l2BlockCh
	l2Num := l2Block.NumberU64()
	t.Logf("insertNewSoftBlock, l2 number: %d", l2Num)
	txsCount, err := driver.BuildSoftBlock(
		t,
		l1cli,
		l2cli,
		l2Num+1,
		0,
		false,
		false,
	)
	t.Nil(err)

	l2Block = <-l2BlockCh
	t.Equal(l2Num+1, l2Block.NumberU64())

	// check txs count.
	t.Equal(txsCount+1, l2Block.Transactions().Len())

	l1Head := <-l1HeadCh
	// check l1Origin variables.
	l1Origin := <-curL1OriginCh
	t.Nil(err, "cannot get l1 origin")
	t.Equal(l2Block.Number(), l1Origin.BlockID)
	t.Equal(l2Block.Hash(), l1Origin.L2BlockHash)
	t.Equal(l1Head.Number, l1Origin.L1BlockHeight)
	t.Equal(l1Head.Hash(), l1Origin.L1BlockHash)
	t.Equal(0, l1Origin.BatchID.Uint64())
	t.Equal(false, l1Origin.EndOfBlock)
	t.Equal(false, l1Origin.EndOfPreconf)
}

func (r PreconfTestSpec) replaceLatestSoftBlock(
	t *hivesim.T,
	l1cli *ethclient.Client,
	l2cli *ethclient.Client,
	driver *clients.DriverClient,
	batchID uint64,
) {
	preBlock := <-l2BlockCh
	t.Logf("replaceLatestSoftBlock, l2 number: %d", preBlock.NumberU64())

	txsCount, err := driver.BuildSoftBlock(
		t,
		l1cli,
		l2cli,
		preBlock.NumberU64(),
		batchID,
		false,
		false,
	)

	l1Origin := <-curL1OriginCh
	if l1Origin.BatchID != nil {
		t.FailIfNotNil(err, "cannot build soft block")
	} else if err == nil || !strings.Contains(err.Error(), "batch ID mismatch") {
		t.Fatalf("if current block is not soft block then the error should contain 'batch ID mismatch', err: %v", err)
	}

	curBlock := <-l2BlockCh
	l1Head := <-l1HeadCh

	t.Equal(preBlock.NumberU64(), curBlock.NumberU64())
	t.Equal(preBlock.Transactions().Len()+txsCount, curBlock.Transactions().Len())

	l1Origin = <-curL1OriginCh
	t.Equal(curBlock.Number(), l1Origin.BlockID)
	t.Equal(curBlock.Hash(), l1Origin.L2BlockHash)
	t.Equal(l1Head.Number, l1Origin.L1BlockHeight)
	t.Equal(l1Head.Hash(), l1Origin.L1BlockHash)
	t.Equal(batchID, l1Origin.BatchID.Uint64())
	t.Equal(false, l1Origin.EndOfBlock)
	t.Equal(false, l1Origin.EndOfPreconf)
}

func (r PreconfTestSpec) randomRemoveSoftBlock(
	t *hivesim.T,
	driver *clients.DriverClient,
) {
	canonicalL1Origin := <-canonicalL1OriginCh
	curL1Origin1 := <-curL1OriginCh
	l2Num := rand.Uint64N(curL1Origin1.BlockID.Uint64()-canonicalL1Origin.BlockID.Uint64()) +
		canonicalL1Origin.BlockID.Uint64()
	// remove soft blocks.
	t.FailIfNotNil(driver.RemoveSoftBlocks(t, l2Num))

	// wait a new current l1Origin
	curL1Origin2 := <-curL1OriginCh
	l2Block := <-l2BlockCh
	t.Equal(l2Num, l2Block.NumberU64())
	t.Equal(curL1Origin2.BlockID, l2Block.Number())

	t.Equal(curL1Origin1.BatchID, curL1Origin2.BatchID)
	t.Equal(curL1Origin1.L1BlockHeight, curL1Origin2.L1BlockHeight)
	t.Equal(curL1Origin1.L1BlockHash, curL1Origin2.L1BlockHash)
	t.Equal(curL1Origin1.L2BlockHash, curL1Origin2.L2BlockHash)
	t.Equal(curL1Origin1.EndOfBlock, curL1Origin2.EndOfBlock)
	t.Equal(curL1Origin1.EndOfPreconf, curL1Origin2.EndOfPreconf)
	t.Equal(curL1Origin1.Preconfer, curL1Origin2.Preconfer)
}

func (r PreconfTestSpec) testInsertSoftBlocksAfterEOB(
	t *hivesim.T,
	driver *clients.DriverClient,
) {

}
