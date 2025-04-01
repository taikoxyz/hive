package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"time"
)

type ProverClient struct {
	Index int
	*HiveManagedClient
	Envs hivesim.Params
	*State
}

func (p *ProverClient) Start() (err error) {
	if err := p.HiveManagedClient.Start(); err != nil {
		return err
	}

	client, err := rpc.NewClient(context.Background(), GetClientConfig(p.Envs))
	if err != nil {
		return err
	}

	p.State, err = NewState(client)
	if err != nil {
		return err
	}

	return err
}

func (p *ProverClient) Shutdown() error {
	if err := p.HiveManagedClient.Shutdown(); err != nil {
		return err
	}
	p.State.Close()

	return nil
}

func (p *ProverClient) VerifyBlocks(opts *bind.TransactOpts) {
	latestVerifyId, err := GetLastVerifiedBlockId(context.Background(), p.Client)
	p.FailIfNotNil(err, "failed to get last verified block id")

	var tx *types.Transaction
	if latestVerifyId < p.PacayaClients.ForkHeight {
		tx, err = p.OntakeClients.TaikoL1.VerifyBlocks(opts, 32)
	} else {
		tx, err = p.PacayaClients.TaikoInbox.VerifyBatches(opts, 32)
	}
	p.FailIfNotNil(err, "failed to verify blocks")

	_, err = bind.WaitMined(context.Background(), p.L1, tx)
	p.FailIfNotNil(err, "failed to wait mined")

	p.Logf("%s: the latest verified id: %d", p.ClientType(), latestVerifyId)
}

func (p *ProverClient) WaitLatestBatchesProved(ctx context.Context, timeout time.Duration, number uint64) error {
	p.Logf("%s: wait latest batches proven: %d", p.ClientType(), number)

	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	sink := make(chan *pacaya.TaikoInboxClientBatchesProved, 10)
	sub, err := p.PacayaClients.TaikoInbox.WatchBatchesProved(&bind.WatchOpts{Context: subCtx}, sink)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for event := range sink {
		if event.BatchIds[len(event.BatchIds)-1] >= number {
			return nil
		}
	}

	return fmt.Errorf("%s: failed to wait batches proved: %d", p.ClientType(), number)
}

func (p *ProverClient) GetLastVerifiedBlockId(ctx context.Context) uint64 {
	lastVerifiedBlockID, err := GetLastVerifiedBlockId(ctx, p.Client)
	p.FailIfNotNil(err, "failed to get last verified block id")
	return lastVerifiedBlockID
}

func (p *ProverClient) WaitLatestVerifiedNumber(ctx context.Context, timeout time.Duration, verifiedNumber uint64) {
	p.Logf("%s: wait latest verified number %d", p.ClientType(), verifiedNumber)

	tmTicker := time.NewTicker(time.Second)
	defer tmTicker.Stop()
	var number uint64
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(timeout):
			p.Fatalf("failed to wait latest verified, expect_number: %d, actual_number: %d", verifiedNumber, number)
		case <-tmTicker.C:
			if number = p.GetLastVerifiedBlockId(ctx); number >= verifiedNumber {
				return
			}
		}
	}
}

func GetLastVerifiedBlockId(ctx context.Context, client *rpc.Client) (uint64, error) {
	var lastVerifiedBlockID uint64
	stateVars, err := client.GetProtocolStateVariablesPacaya(&bind.CallOpts{Context: ctx})
	if err != nil {
		slot1, _, err := client.GetProtocolStateVariablesOntake(&bind.CallOpts{Context: ctx})
		if err != nil {
			return 0, err
		}
		lastVerifiedBlockID = slot1.LastSyncedBlockId
	} else {
		lastVerifiedBlockID = stateVars.Stats2.LastVerifiedBatchId
	}

	return lastVerifiedBlockID, nil
}
