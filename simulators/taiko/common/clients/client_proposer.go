package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
)

type ProposerClient struct {
	Index int
	*HiveManagedClient
	Envs hivesim.Params
	*State
}

func (p *ProposerClient) Start() (err error) {
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

func (p *ProposerClient) Shutdown() error {
	if err := p.HiveManagedClient.Shutdown(); err != nil {
		return err
	}
	p.State.Close()

	return nil
}

func (p *ProposerClient) ProposeTxLists(
	ctx context.Context,
	t *hivesim.T,
) (*rawdb.L1Origin, []types.Transactions, error) {
	canonicalL1Origin, err := p.L2.HeadL1Origin(ctx)
	if err != nil {
		return nil, nil, err
	}

	l2Number, err := p.L2.BlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}

	// no soft blocks need to be proposed
	if canonicalL1Origin.BlockID.Uint64() == l2Number {
		return nil, nil, nil
	}

	// Collect all the soft transactions.
	var (
		total int
		txs   []types.Transactions
	)
	for number := canonicalL1Origin.BlockID.Uint64() + 1; number <= l2Number; number++ {
		l2Block, err := p.L2.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, l2Block.Transactions())
		total += l2Block.Transactions().Len()
	}
	t.Logf("propose tx list, batch count: %d, txs count: %d", len(txs), total)

	return canonicalL1Origin, txs, nil //p.Proposer.ProposeTxLists(ctx, txs)
}

func (p *ProposerClient) GetProposeEvents(ctx context.Context, start uint64) []*pacaya.TaikoInboxClientBatchProposed {
	iter, err := p.PacayaClients.TaikoInbox.FilterBatchProposed(&bind.FilterOpts{Context: ctx, Start: start})
	p.FailIfNotNil(err)

	var events []*pacaya.TaikoInboxClientBatchProposed
	if iter.Next() {
		events = append(events, iter.Event)
	}

	return events
}
