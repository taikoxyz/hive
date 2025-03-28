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

type ReorgParams struct {
	DelayTime   int64 // seconds
	DelayNumber int   // blocks
}

func (a *AnvilClient) StartRecordReorgPoints(ctx context.Context, l2cli *rpc.EthClient) {
	a.reorgCh = make(chan struct{})
	a.reorgCache = make(map[uint64]*L1BlockInfo)
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
	}
	a.reorgCh = nil
	a.reorgCache = nil
}

func (a *AnvilClient) Reorg(l2Number uint64, params *ReorgParams) {
	a.StopMining()
	defer a.StartMining()

	var (
		ctx      = context.Background()
		client   = a.EthClient
		l1Number uint64
		snapshot string
	)

	if params == nil {
		params = &ReorgParams{}
	}
	a.Logf("%s: start reorg l1chain, l2_number: %d, delay_time: %d, delay_number: %d", a.ClientType(), l2Number, params.DelayTime, params.DelayNumber)

	nums := maps.Keys(a.reorgCache)
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	for _, num := range nums {
		info := a.reorgCache[num]
		if info.L2Number > l2Number {
			break
		}
		l1Number = num
		snapshot = info.Snapshot
		a.Logf("%s: check reorg point, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number)
	}

	var (
		startTime  int64
		blockCount = params.DelayNumber
		txs        []*types.Transaction
	)
	for l1Number += 1; true; l1Number++ {
		block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(l1Number))
		if err != nil {
			break
		}
		if startTime == 0 {
			startTime = int64(block.Time())
		}
		txs = append(txs, block.Transactions()...)
		blockCount++
		delete(a.reorgCache, l1Number)
		a.Logf("%s: reorg l1 chain, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number)
	}
	if txs == nil || blockCount <= 0 || startTime+params.DelayTime <= 0 {
		a.Errorf("%s: no txs to reorg, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number)
		return
	}

	// Revert l1 chain.
	a.RevertSnapshot(snapshot)

	var curTxs []*types.Transaction
	for i := 0; i < blockCount; i++ {
		count := min((len(txs)+blockCount-1)/blockCount, len(txs))
		curTxs, txs = txs[:count], txs[count:]
		for _, tx := range curTxs {
			err := client.SendTransaction(ctx, tx)
			a.FailIfNotNil(err, fmt.Sprintf("failed to send tx %s, err: %v", tx.Hash().Hex(), err))
		}
		a.SetNextBlockTimestamp(uint64(startTime+params.DelayTime) + uint64(i)*a.SecondsPerSlot)
		a.MineBlock()
		a.Logf("%s: mint a new l1 block, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number+uint64(i))
	}
}
