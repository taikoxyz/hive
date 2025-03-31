package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
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
	a.reorgCh = make(chan struct{})
	a.reorgPoints = make(map[uint64]string)
	a.l1Origins = make(map[uint64]*rawdb.L1Origin)

	var l1Number uint64

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
				if l1Number == 0 {
					continue
				}

				l1Num, err := a.EthClient.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get %s latest number, err: %v", a.ClientType(), err)
				}
				if a.reorgPoints[l1Num] == "" {
					a.reorgPoints[l1Num] = a.SetSnapshot()
				}
			case <-tL1Origin.C:
				l2Num, err := l2cli.BlockNumber(ctx)
				if err != nil {
					a.Fatalf("failed to get l2node latest number, err: %v", err)
				}

				if a.l1Origins[l2Num] == nil {
					l1Origin, err := l2cli.L1OriginByID(ctx, big.NewInt(int64(l2Num)))
					if err != nil && !strings.Contains(err.Error(), "not found") {
						a.Fatalf("failed to get l1origin by id %d, err: %v", l2Num, err)
					}
					if l1Origin == nil {
						continue
					}

					a.l1Origins[l2Num] = l1Origin

					a.Logf("record reorg point, l1_number: %d, l2_number: %d", l1Origin.L1BlockHeight.Uint64(), l2Num)
					l1Number = l1Origin.L1BlockHeight.Uint64()
				}
			}
		}
	}()
}

func (a *AnvilClient) Reorg(l2Number uint64, params *ReorgParams) uint64 {
	a.StopMining()
	defer a.StartMining()
	defer a.StopRecordReorgPoints()

	var (
		ctx      = context.Background()
		client   = a.EthClient
		l1Number = a.l1Origins[l2Number].L1BlockHeight.Uint64()
		snapshot string
	)

	if params == nil {
		params = &ReorgParams{}
	}

	// Get snapshot.
	for ; true; l2Number++ {
		l1Origin, ok := a.l1Origins[l2Number]
		if !ok {
			break
		}
		if l1Origin.L1BlockHeight.Uint64() > l1Number {
			l1Number++
			break
		}
		snapshot = a.reorgPoints[l1Number]
	}

	a.Logf("%s: start reorg l1chain, l1_number: %d, l2_number: %d, delay_time: %d, delay_number: %d", a.ClientType(), l1Number, l2Number, params.DelayTime, params.DelayNumber)

	var (
		startTime  int64
		blockCount = params.DelayNumber
		txs        []*types.Transaction
	)
	for l1Num := l1Number; true; l1Num++ {
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
	if txs == nil || blockCount <= 0 || startTime+params.DelayTime <= 0 {
		a.Errorf("%s: no txs to reorg, l2_number: %d, l1_number: %d", a.ClientType(), l2Number, l1Number)
		return 0
	}

	// Revert l1 chain.
	a.RevertSnapshot(snapshot)
	if number, err := client.BlockNumber(ctx); err != nil {
		a.Errorf("%s: failed to get block number, err: %v", a.ClientType(), err)
		return 0
	} else if number+1 != l1Number {
		a.Errorf("%s: failed to revert l1 chain, expect: %d, actual: %d", a.ClientType(), l1Number-1, number)
		return 0
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
		a.Logf("%s: mint a new l1 block, l2_number: %d, l1_number: %d, tx_count: %d, lest: %d", a.ClientType(), l2Number+uint64(i), l1Number+uint64(i), len(curTxs), len(txs))
	}

	return l2Number
}
