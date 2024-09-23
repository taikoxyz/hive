package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

func (a *AnvilClient) StartRecordReorgPoints(ctx context.Context, l2cli *ethclient.Client) {

	a.reorgCache = make(map[uint64]*L1BlockInfo)
	a.reorgCh = make(chan struct{})

	var (
		l1Number uint64
		l1Cli    = a.EthClient()
	)
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-a.reorgCh:
				return
			case <-tick.C:
				l2Num, err := l2cli.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get l2node latest number, err: %v", err)
				}

				l1Num, err := l1Cli.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get %s latest number, err: %v", a.ClientType(), err)
				}
				if l1Num > l1Number {
					l1Number = l1Num
					a.reorgCache[l1Num] = &L1BlockInfo{
						L2Number: l2Num,
						Snapshot: a.SetSnapshot(),
					}
				}
			}
		}
	}()
}

func (a *AnvilClient) Reorg(l2Number uint64) {
	close(a.reorgCh)
	a.StopMining()
	defer a.StartMining()

	var (
		ctx      = context.Background()
		client   = a.EthClient()
		l1Number uint64
		snapshot string
	)

	for num, info := range a.reorgCache {
		if info.L2Number == l2Number {
			if l1Number == 0 || num < l1Number {
				l1Number = num
				snapshot = info.Snapshot
			}
		}
	}

	blocks := make([]*types.Block, 0)
	for num := l1Number + 1; true; num++ {
		block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(num))
		if err != nil {
			break
		}
		blocks = append(blocks, block)
	}

	a.RevertSnapshot(snapshot)
	//a.SetNextBlockTimestamp(info.Header.Time + a.SecondsPerSlot + 1)
	//a.SetNextBlockTimestamp(blocks[0].Time() + 1)

	for _, block := range blocks {
		for _, tx := range block.Transactions() {
			if err := client.SendTransaction(ctx, tx); err != nil {
				a.Fatalf("failed to send tx %s, err: %v", tx.Hash().Hex(), err)
			}
		}
		//if err := a.WaitLatestNumber(ctx, time.Second*60, block.NumberU64()); err != nil {
		//	a.Fatalf("reorg failed to wait latest number, err: %v", err)
		//}
		//time.Sleep(time.Duration(a.SecondsPerSlot) * time.Second)
		a.MineBlock()
	}

	a.reorgCh = nil
	a.reorgCache = nil
}
