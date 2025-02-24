package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func (p *ProverClient) VerifyBlocks(opts *bind.TransactOpts) error {
	tx, err := p.OntakeClients.TaikoL1.VerifyBlocks(opts, 32)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), p.L1, tx)
	return err
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
	var lastVerifiedBlockID uint64
	stateVars, err := p.Client.GetProtocolStateVariablesPacaya(&bind.CallOpts{Context: ctx})
	if err != nil {
		slot1, _, err := p.Client.GetProtocolStateVariablesOntake(&bind.CallOpts{Context: ctx})
		p.FailIfNotNil(err, "failed to get protocol state variables")
		lastVerifiedBlockID = slot1.LastSyncedBlockId
	} else {
		lastVerifiedBlockID = stateVars.Stats2.LastVerifiedBatchId
	}

	return lastVerifiedBlockID
}

func (p *ProverClient) WaitLatestVerifiedNumber(ctx context.Context, timeout time.Duration, verifiedNumber uint64) error {
	p.Logf("%s: wait latest verified number %d", p.ClientType(), verifiedNumber)
	current, times := uint64(0), timeout/time.Second
	for times > 0 && verifiedNumber >= current {
		select {
		case <-time.Tick(time.Second):
			if number := p.GetLastVerifiedBlockId(ctx); number >= current {
				current = number + 1
				times = timeout / time.Second
				break
			} else {
				times--
			}
		}
	}

	if verifiedNumber >= current {
		return fmt.Errorf("%s failed to reach current number %d, current number: %d", p.ClientType(), verifiedNumber, current)
	}
	return nil
}
