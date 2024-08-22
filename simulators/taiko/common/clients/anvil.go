package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

var (
	AnvilPort int64 = 8545
)

type L1BlockInfo struct {
	Header   *types.Header
	Snapshot string
}

type AnvilClient struct {
	SecondsPerSlot uint64
	*EthNode

	reorgCache map[string]*L1BlockInfo
}

func (a *AnvilClient) SetReorgPoint() (string, uint64) {
	client := a.EthClient()
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		a.Fatalf("failed to get %s latest header, err: %v", a.ClientType(), err)
	}
	if a.reorgCache == nil {
		a.reorgCache = make(map[string]*L1BlockInfo)
	}
	snapshot := a.setSnapshot()
	a.reorgCache[snapshot] = &L1BlockInfo{
		Snapshot: snapshot,
		Header:   header,
	}
	return snapshot, header.Number.Uint64()
}

func (a *AnvilClient) Reorg(snapshot string) {
	info, ok := a.reorgCache[snapshot]
	if !ok {
		return
	}
	a.StopMining()
	defer a.StartMining()

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

	a.revertSnapshot(info.Snapshot)
	//a.SetNextBlockTimestamp(info.Header.Time + a.SecondsPerSlot + 1)

	for _, block := range blocks {
		for _, tx := range block.Transactions() {
			err = client.SendTransaction(ctx, tx)
			if err != nil {
				a.Fatalf("failed to send tx %s, err: %v", tx.Hash().Hex(), err)
			}
		}
		a.MineBlock()
		time.Sleep(3 * time.Second)
	}

	delete(a.reorgCache, snapshot)
}

func (a *AnvilClient) MineBlock() {
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_mine")
	if err != nil {
		a.Fatalf("failed to mine block, err: %v", err)
	}
}

func (a *AnvilClient) StartMining() {
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_setIntervalMining", a.SecondsPerSlot)
	if err != nil {
		a.Fatalf("failed to start mining, err: %v", err)
	}
}

func (a *AnvilClient) StopMining() {
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_setIntervalMining", 0)
	if err != nil {
		a.Fatalf("failed to stop mining, err: %v", err)
	}
}

func (a *AnvilClient) SetNextBlockTimestamp(timestamp uint64) {
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_setNextBlockTimestamp", timestamp)
	if err != nil {
		a.Fatalf("failed to set next block timestamp, err: %v", err)
	}
}

func (a *AnvilClient) setSnapshot() (snapshotID string) {
	client := a.RPClient()
	err := client.CallContext(context.Background(), &snapshotID, "evm_snapshot")
	if err != nil {
		a.Fatalf("failed to take snapshot, err: %v", err)
	}
	return snapshotID
}

func (a *AnvilClient) revertSnapshot(snapshotID string) {
	//a.Logf("revert %s snapshot %s", a.ClientType(), snapshotID)
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_revert", snapshotID)
	if err != nil {
		a.Fatalf("failed to revert snapshot, err: %v", err)
	}
}
