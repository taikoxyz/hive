package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"golang.org/x/exp/maps"
	"math/big"
	"sort"
	"time"
)

func (a *AnvilClient) StartRecordReorgPoints(ctx context.Context, l2cli *rpc.EthClient) {
	a.reorgCh = make(chan struct{})
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

				l1Num, err := a.EthClient.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get %s latest number, err: %v", a.ClientType(), err)
				}
				if a.reorgCache[l1Num] == nil {
					a.reorgCache[l1Num] = &L1BlockInfo{
						L2Number: l2Num,
						Snapshot: a.SetSnapshot(),
					}
					a.Logf("record reorg point, l2_number: %d, l1_number: %d", l2Num, l1Num)
				}
			}
		}
	}()
}

func (a *AnvilClient) StopRecordReorgPoints() {
	if a.reorgCh != nil {
		close(a.reorgCh)
		a.reorgCh = nil
	}
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

	nums := maps.Keys(a.reorgCache)
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	for _, num := range nums {
		info := a.reorgCache[num]
		if info.L2Number > l2Number {
			break
		}
		l1Number = num
		snapshot = info.Snapshot
		a.Logf("check reorg point, l2_number: %d, l1_number: %d", l2Number, l1Number)
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
		a.Logf("reorg l1chain block, l2_number: %d, l1_number: %d, hash: %s", l2Number, block.NumberU64(), block.Hash().Hex())
		for _, tx := range block.Transactions() {
			err := client.SendTransaction(ctx, tx)
			a.FailIfNotNil(err, fmt.Sprintf("failed to send tx %s, err: %v", tx.Hash().Hex(), err))
		}
		a.MineBlock()
	}

	l1Number -= 1

	a.WaitLatestNumber(ctx, time.Minute*3, l1Number)
}
