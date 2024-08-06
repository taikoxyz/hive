package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/marioevz/eth-clients/clients"
	"taiko2/common/utils"
	"time"
)

type TaikoGethClient struct {
	Client
	Logger utils.Logging
	eth    *ethclient.Client
}

func (t *TaikoGethClient) Logf(format string, values ...interface{}) {
	if l := t.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (t *TaikoGethClient) Start() error {
	if !t.Client.IsRunning() {
		if managedClient, ok := t.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			return managedClient.Start()
		}
	}

	return t.Init(context.Background())
}

func (t *TaikoGethClient) Init(ctx context.Context) error {
	var err error
	t.eth, err = ethclient.Dial(t.HttpURL())
	if err != nil {
		return err
	}
	// Wait until taiko geth is ready.
	for {
		if _, err = t.eth.ChainID(ctx); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 2):
		}
	}
}

func (t *TaikoGethClient) Shutdown() error {
	if managedClient, ok := t.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}

func (t *TaikoGethClient) HttpURL() string {
	return fmt.Sprintf("http://%v:%d", t.GetHost(), PortHttpRPC)
}
