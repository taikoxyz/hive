package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

type ProverClient struct {
	*HiveManagedClient

	L1Auth *bind.TransactOpts
	*State
}

func (p *ProverClient) Start() (err error) {
	if err := p.HiveManagedClient.Start(); err != nil {
		return err
	}

	client, err := rpc.NewClient(context.Background(), GetClientConfig())
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
	if opts == nil {
		opts = p.L1Auth
	}
	tx, err := p.OntakeClients.TaikoL1.VerifyBlocks(opts, 32)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), p.L1, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("failed to verify blocks")
	}

	return nil
}
