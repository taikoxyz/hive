package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

func (a *AnvilClient) SetReorgPoint() uint64 {
	client := a.EthClient()

	a.reorgCache = make(map[uint64]*L1BlockInfo)
	a.reorgCh = make(chan struct{})

	number, err := client.BlockNumber(context.Background())
	if err != nil {
		a.Fatalf("failed to get %s latest number, err: %v", a.ClientType(), err)
	}
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-a.reorgCh:
				return
			case <-tick.C:
				header, err := client.HeaderByNumber(context.Background(), nil)
				if err != nil {
					a.Fatalf("failed to get %s latest header, err: %v", a.ClientType(), err)
				}
				if num := header.Number.Uint64(); num > number {
					number = num
					a.reorgCache[number] = &L1BlockInfo{
						Snapshot: a.SetSnapshot(),
						Header:   header,
					}
				}
			}
		}
	}()

	return number
}

func (a *AnvilClient) Reorg(number uint64) {
	close(a.reorgCh)
	info, ok := a.reorgCache[number]
	if !ok {
		return
	}
	a.StopMining()
	defer a.StartMining()

	a.Logf("reorg %s to the number %d", a.ClientType(), info.Header.Number.Uint64())

	var (
		ctx    = context.Background()
		client = a.EthClient()
	)

	curNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		a.Fatalf("failed to get %s latest header, err: %v", a.ClientType(), err)
	}
	blocks := make([]*types.Block, 0)
	for num := info.Header.Number.Uint64() + 1; num <= curNumber; num++ {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(num))
		if err != nil {
			a.Fatalf("failed to get block %d, err: %v", num, err)
		}
		blocks = append(blocks, block)
	}

	a.RevertSnapshot(info.Snapshot)

	for i, block := range blocks {
		for _, tx := range block.Transactions() {
			err = client.SendTransaction(ctx, tx)
			if err != nil {
				a.Fatalf("failed to send tx %s, err: %v", tx.Hash().Hex(), err)
			}
		}
		if i == 0 {
			a.SetNextBlockTimestamp(block.Time() + 1)
		} else {
			a.SetNextBlockTimestamp(block.Time())
		}
		a.MineBlock()
	}

	a.reorgCh = nil
	a.reorgCache = nil
}
