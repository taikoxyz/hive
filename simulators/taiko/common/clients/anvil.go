package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	AnvilPort int64 = 8545
)

type AnvilClient struct {
	*BaseNode
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
