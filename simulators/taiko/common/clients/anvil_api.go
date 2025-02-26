package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"time"
)

func (a *AnvilClient) MineBlock() {
	err := a.EthClient.CallContext(context.Background(), nil, "evm_mine")
	a.FailIfNotNil(err, fmt.Sprintf("failed to mine block, err: %v", err))
}

func (a *AnvilClient) StartMining() {
	err := a.EthClient.CallContext(context.Background(), nil, "evm_setIntervalMining", a.SecondsPerSlot)
	if err != nil {
		a.Fatalf("failed to start mining, err: %v", err)
	}
}

func (a *AnvilClient) StopMining() {
	err := a.EthClient.CallContext(context.Background(), nil, "evm_setIntervalMining", 0)
	if err != nil {
		a.Fatalf("failed to stop mining, err: %v", err)
	}
}

func (a *AnvilClient) SetNextBlockTimestamp(timestamp uint64) {
	err := a.EthClient.CallContext(context.Background(), nil, "evm_setNextBlockTimestamp", timestamp)
	if err != nil {
		a.Fatalf("failed to set next block timestamp, err: %v", err)
	}
}

func (a *AnvilClient) IncreaseTime(timestamp uint64) {
	a.Logf("%s: increase time: %ds", a.ClientType(), timestamp)
	err := a.EthClient.CallContext(context.Background(), nil, "evm_increaseTime", timestamp)
	if err != nil {
		a.Fatalf("failed to increase time, err: %v", err)
	}
}

func (a *AnvilClient) SetSnapshot() (snapshotID string) {
	err := a.EthClient.CallContext(context.Background(), &snapshotID, "evm_snapshot")
	if err != nil {
		a.Fatalf("failed to take snapshot, err: %v", err)
	}
	return snapshotID
}

func (a *AnvilClient) RevertSnapshot(snapshotID string) {
	err := a.EthClient.CallContext(context.Background(), nil, "evm_revert", snapshotID)
	if err != nil {
		a.Fatalf("failed to revert snapshot, err: %v", err)
	}
}

func (a *AnvilClient) GetLastVerifiedBlockId(ctx context.Context, l1cli *rpc.EthClient) (id uint64) {
	pacayal1, err := pacaya.NewPacayaL1Clients(l1cli)
	if err != nil {
		a.Fatal(err)
	}
	ontakel1, err := ontake.NewOntakeL1Clients(l1cli)
	if err != nil {
		a.Fatal(err)
	}
	state2, err := pacayal1.TaikoInbox.GetStats2(&bind.CallOpts{Context: ctx})
	if err != nil {
		_, slotB, err := ontakel1.TaikoL1.GetStateVariables(&bind.CallOpts{Context: ctx})
		if err != nil {
			a.Fatal(err)
		}
		return slotB.LastVerifiedBlockId
	}

	return state2.LastVerifiedBatchId
}

func (a *AnvilClient) WaitLatestVerifiedNumber(ctx context.Context, timeout time.Duration, verifiedNumber uint64) error {
	a.Logf("%s: wait latest verified number %d", a.ClientType(), verifiedNumber)
	current, times := uint64(0), timeout/time.Second
	for times > 0 && verifiedNumber >= current {
		select {
		case <-time.Tick(time.Second):
			if number := a.GetLastVerifiedBlockId(ctx, a.EthClient); number >= current {
				current = number + 1
				times = timeout / time.Second
				break
			} else {
				times--
			}
		}
	}

	if verifiedNumber >= current {
		return fmt.Errorf("%s failed to reach current number %d, current number: %d", a.ClientType(), verifiedNumber, current)
	}
	return nil
}
