package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"taiko2/common/utils"
	"time"
)

type TaikoGethClient struct {
	*HiveManagedClient
	Logger    utils.Logging
	Network   string
	networkIP string

	eth *ethclient.Client
}

func (t *TaikoGethClient) Logf(format string, values ...interface{}) {
	if l := t.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (t *TaikoGethClient) Start() error {
	t.Logf("Starting taiko geth client")
	if !t.IsRunning() {
		return t.HiveManagedClient.Start()
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
	return t.HiveManagedClient.Shutdown()
}

func (t *TaikoGethClient) NetworkIP() string {
	if t.networkIP != "" {
		return t.networkIP
	}

	var err error
	t.networkIP, err = t.T.Sim.ContainerNetworkIP(t.T.SuiteID, t.Network, t.Client.Container)
	if err != nil {
		t.Logf("Error getting network IP: %v", err)
		return t.GetHost()
	}
	t.T.Logf("taiko geth network IP %v", t.networkIP)

	return t.networkIP
}

func (t *TaikoGethClient) HttpURL() string {
	return fmt.Sprintf("http://%v:%d", t.GetHost(), PortHttpRPC)
}
