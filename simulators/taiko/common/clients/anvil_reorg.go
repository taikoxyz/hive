package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math"
	"math/big"
	"strings"
	"time"
)

type ReorgParams struct {
	DelayTime   int64 // seconds
	DelayNumber int   // blocks
}

func (a *AnvilClient) StopRecordReorgPoints() {
	if a.reorgCh != nil {
		close(a.reorgCh)
	}
	a.reorgCh = nil
	a.reorgPoints = nil
	a.l1Origins = nil
}

func (a *AnvilClient) StartRecordReorgPoints(ctx context.Context, l2cli *rpc.EthClient) {
	if a.reorgCh != nil {
		return
	}
	a.reorgCh = make(chan struct{})
	a.reorgPoints = make(map[uint64]string)
	a.l1Origins = make(map[uint64]*rawdb.L1Origin)

	go func() {
		tL1Origin := time.NewTicker(time.Second)
		defer tL1Origin.Stop()

		tL1Header := time.NewTicker(time.Second)
		defer tL1Header.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-a.reorgCh:
				return
			case <-tL1Header.C:
				l1Num, err := a.EthClient.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get %s latest number, err: %v", a.ClientType(), err)
				}
				if a.reorgPoints[l1Num] == "" {
					a.reorgPoints[l1Num] = a.SetSnapshot()
				}
			case <-tL1Origin.C:
				headL1Origin, err := l2cli.HeadL1Origin(ctx)
				if err != nil && !strings.Contains(err.Error(), "not found") {
					a.Fatalf("failed to get %s latest head l1 origin, err: %v", a.ClientType(), err)
				}
				if headL1Origin == nil {
					continue
				}

				l2Num := headL1Origin.BlockID.Uint64()
				if a.l1Origins[l2Num] == nil {
					a.l1Origins[l2Num] = headL1Origin
					a.Logf("record reorg point, l1_number: %d, l2_number: %d", headL1Origin.L1BlockHeight.Uint64(), l2Num)
				}
			}
		}
	}()
}

func (a *AnvilClient) Reorg(l2Number uint64, params *ReorgParams) {
	a.StopMining()
	defer a.StartMining()
	defer a.StopRecordReorgPoints()

	var (
		ctx    = context.Background()
		client = a.EthClient
		// reorg point.
		l1Number uint64
		snapshot string
	)

	if l1Origin := a.l1Origins[l2Number]; l1Origin != nil {
		l1Number = l1Origin.L1BlockHeight.Uint64()
	} else {
		var l2Num = uint64(math.MaxUint64)
		for num, val := range a.l1Origins {
			if num > l2Number {
				l2Num = min(l2Num, num)
				l1Origin = val
			}
		}
		if l1Origin == nil {
			a.Fatalf("cannot find l1 origin for l2_number: %d", l2Number)
			return
		}
		l1Number = l1Origin.L1BlockHeight.Uint64() - 1
	}
	snapshot = a.reorgPoints[l1Number]

	if params == nil {
		params = &ReorgParams{}
	}

	a.Logf("%s: start reorg l1chain, l1_number: %d, l2_number: %d, delay_time: %d, delay_number: %d", a.ClientType(), l1Number, l2Number, params.DelayTime, params.DelayNumber)

	var (
		startTime  int64
		blockCount = params.DelayNumber
		txs        []*types.Transaction
	)
	for l1Num := l1Number + 1; true; l1Num++ {
		block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(l1Num))
		if err != nil {
			break
		}
		if startTime == 0 {
			startTime = int64(block.Time())
		}
		txs = append(txs, block.Transactions()...)
		blockCount++
		a.Logf("%s: reorg l1 chain, l1_number: %d", a.ClientType(), l1Num)
	}
	if blockCount <= 0 || startTime+params.DelayTime <= 0 {
		a.Errorf("%s: no txs to reorg, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number)
		return
	}

	// Revert l1 chain.
	a.RevertSnapshot(snapshot)
	if number, err := client.BlockNumber(ctx); err != nil {
		a.Errorf("%s: failed to get block number, err: %v", a.ClientType(), err)
		return
	} else if number != l1Number {
		a.Errorf("%s: failed to revert l1 chain, expect: %d, actual: %d", a.ClientType(), l1Number, number)
		return
	}

	var (
		curTxs []*types.Transaction
		count  = min((len(txs)+blockCount-1)/blockCount, len(txs))
	)
	for i := 0; i < blockCount; i++ {
		curTxs, txs = txs[:min(count, len(txs))], txs[min(count, len(txs)):]
		for _, tx := range curTxs {
			err := client.SendTransaction(ctx, tx)
			a.FailIfNotNil(err, fmt.Sprintf("failed to send tx %s, err: %v", tx.Hash().Hex(), err))
		}
		a.SetNextBlockTimestamp(uint64(startTime) + uint64(i)*a.SecondsPerSlot + uint64(params.DelayTime))
		a.MineBlock()
		a.Logf("%s: mint a new l1 block, l1_number: %d, tx_count: %d, lest: %d", a.ClientType(), l1Number+1+uint64(i), len(curTxs), len(txs))
	}
}
