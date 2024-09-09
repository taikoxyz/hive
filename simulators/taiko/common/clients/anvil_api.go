package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"taiko/bindings/taikol1"
	"time"
)

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

func (a *AnvilClient) IncreaseTime(timestamp uint64) {
	a.Logf("%s: increase time: %ds", a.ClientType(), timestamp)
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_increaseTime", timestamp)
	if err != nil {
		a.Fatalf("failed to increase time, err: %v", err)
	}
}

func (a *AnvilClient) SetSnapshot() (snapshotID string) {
	client := a.RPClient()
	err := client.CallContext(context.Background(), &snapshotID, "evm_snapshot")
	if err != nil {
		a.Fatalf("failed to take snapshot, err: %v", err)
	}
	return snapshotID
}

func (a *AnvilClient) RevertSnapshot(snapshotID string) {
	//a.Logf("revert %s snapshot %s", a.ClientType(), snapshotID)
	client := a.RPClient()
	err := client.CallContext(context.Background(), nil, "evm_revert", snapshotID)
	if err != nil {
		a.Fatalf("failed to revert snapshot, err: %v", err)
	}
}

func (a *AnvilClient) GetTaikoDataSlotB(ctx context.Context) *taikol1.TaikoDataSlotB {
	_, slotB, err := a.taikoL1.GetStateVariables(&bind.CallOpts{Context: ctx})
	if err != nil {
		a.Fatal(err)
	}
	return &slotB
}

func (a *AnvilClient) WaitVerifiedNumber(ctx context.Context, timeout time.Duration) (uint64, uint64, error) {

	// record reorg point
	snapshotNumber := a.SetReorgPoint()

	stopCh := a.HandleProposedEvent(0)
	defer close(stopCh)

	a.Logf("%s: wait latest l2 blockVerified, l1 snapshot number: %d", a.ClientType(), snapshotNumber)

	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	verifiedCh := make(chan *taikol1.TaikoL1BlockVerified, 10)
	verifiedChV2 := make(chan *taikol1.TaikoL1BlockVerifiedV2, 10)
	sub, err := a.taikoL1.WatchBlockVerified(&bind.WatchOpts{Context: subCtx}, verifiedCh, nil, nil)
	if err != nil {
		return 0, 0, err
	}
	defer sub.Unsubscribe()
	sub2, err := a.taikoL1.WatchBlockVerifiedV2(&bind.WatchOpts{Context: subCtx}, verifiedChV2, nil, nil)
	if err != nil {
		return 0, 0, err
	}
	defer sub2.Unsubscribe()

	tick := time.NewTicker(time.Second)
	var (
		number           uint64
		l1Number         uint64
		latestVerifiedId uint64
	)
	for latestVerifiedId == 0 || len(verifiedCh) != 0 || len(verifiedChV2) != 0 {
		select {
		case <-ctx.Done():
			return 0, 0, ctx.Err()

		case <-tick.C:
			if l1Number < snapshotNumber {
				num, err := a.EthClient().BlockNumber(ctx)
				if err != nil {
					a.Fatalf("%s: failed to get latest number, err: %v", a.ClientType(), err)
				}
				if num >= number {
					number = num
					a.depProposerEvent()
				}
			}

		case result := <-verifiedCh:
			verifiedBlockId := result.BlockId.Uint64()
			a.Logf("%s: get taiko data SlotB, last verifiedBlockId: %d", a.ClientType(), verifiedBlockId)
			if result.Raw.BlockNumber >= snapshotNumber {
				l1Number, latestVerifiedId = result.Raw.BlockNumber, verifiedBlockId
				time.Sleep(time.Millisecond * 200)
				break
			}

		case result := <-verifiedChV2:
			verifiedBlockId := result.BlockId.Uint64()
			a.Logf("%s: get taiko data SlotB, last verifiedBlockId: %d", a.ClientType(), verifiedBlockId)
			if result.Raw.BlockNumber >= snapshotNumber {
				l1Number, latestVerifiedId = result.Raw.BlockNumber, verifiedBlockId
				time.Sleep(time.Millisecond * 200)
				break
			}
		}
	}

	a.Logf("%s: verified channel length: %d, %d", a.ClientType(), len(verifiedCh), len(verifiedChV2))

	return l1Number, latestVerifiedId, nil
}
