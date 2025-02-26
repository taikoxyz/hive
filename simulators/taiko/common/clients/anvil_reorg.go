package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"time"
)

func (a *AnvilClient) StartRecordReorgPoints(ctx context.Context, l2cli *rpc.EthClient) {
	if a.reorgCh != nil {
		return
	}
	a.reorgCh = make(chan struct{})

	var (
		l1Number uint64
		l1Cli    = a.EthClient
	)

	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
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
	a.StopMining()
	defer a.StartMining()

	var (
		ctx      = context.Background()
		client   = a.EthClient
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
	for l1Number += 1; true; {
		block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(l1Number))
		if err != nil {
			break
		}
		blocks = append(blocks, block)
		delete(a.reorgCache, l1Number)
		l1Number++
	}

	a.RevertSnapshot(snapshot)

	for _, block := range blocks {
		for _, tx := range block.Transactions() {
			err := client.SendTransaction(ctx, tx)
			a.FailIfNotNil(err, fmt.Sprintf("failed to send tx %s, err: %v", tx.Hash().Hex(), err))
		}
		a.MineBlock()
	}

	a.FailIfNotNil(a.WaitLatestNumber(ctx, time.Second*30, l1Number), fmt.Sprintf("cannot wait for l1 number: %d", l1Number))
}
