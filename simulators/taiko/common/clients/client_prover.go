package clients

import (
	"context"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

type ProverClient struct {
	*HiveManagedClient

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
