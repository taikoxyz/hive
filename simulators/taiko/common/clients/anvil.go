package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	AnvilPort int64 = 8545
)

type AnvilClient struct {
	*EthNode
}

func (a *AnvilClient) SetL1Snapshot() string {
	client, err := rpc.Dial(a.HttpURL())
	if err != nil {
		a.Fatalf("failed to dial anvil client, err: %v", err)
	}
	var snapshotID string
	err = client.CallContext(context.Background(), &snapshotID, "evm_snapshot")
	if err != nil {
		a.Fatalf("failed to take snapshot, err: %v", err)
	}
	return snapshotID
}

func (a *AnvilClient) RevertL1Snapshot(snapshotID string) {
	client, err := rpc.Dial(a.HttpURL())
	if err != nil {
		a.Fatalf("failed to dial anvil client, err: %v", err)
	}
	err = client.CallContext(context.Background(), nil, "evm_revert", snapshotID)
	if err != nil {
		a.Fatalf("failed to revert snapshot, err: %v", err)
	}
}

func (a *AnvilClient) SetL1Automine(automine bool) {
	client, err := rpc.Dial(a.HttpURL())
	if err != nil {
		a.Fatalf("failed to dial anvil client, err: %v", err)
	}
	err = client.CallContext(context.Background(), nil, "evm_setAutomine", automine)
	if err != nil {
		a.Fatalf("failed to set automine, err: %v", err)
	}
}

func (a *AnvilClient) SetIntervalMining(interval uint64) {
	client, err := rpc.Dial(a.HttpURL())
	if err != nil {
		a.Fatalf("failed to dial anvil client, err: %v", err)
	}
	err = client.CallContext(context.Background(), nil, "evm_setIntervalMining", interval)
	if err != nil {
		a.Fatalf("failed to set interval mining, err: %v", err)
	}
}
