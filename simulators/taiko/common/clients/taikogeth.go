package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"taiko/common/utils"
	"time"
)

type TaikoGethClient struct {
	*HiveManagedClient
	Logger    utils.Logging
	Network   string
	networkIP string
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

func (t *TaikoGethClient) EthIsReady(ctx context.Context, timeout time.Duration) error {
	client, err := ethclient.DialContext(ctx, t.HttpURL())
	if err != nil {
		return err
	}
	for ; ; <-time.Tick(time.Second) {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("reach timeout but l1geth is not ready")
		default:
			_, err = client.ChainID(ctx)
			if err != nil {
				continue
			}
			return nil
		}
	}
}

func (t *TaikoGethClient) VerifyNumber(ctx context.Context, timeout time.Duration, targetNumber uint64) error {
	client, err := ethclient.DialContext(ctx, t.HttpURL())
	if err != nil {
		return err
	}
	for ; ; <-time.Tick(time.Second) {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("reach timeout but l2geth is not ready")
		default:
			number, err := client.BlockNumber(ctx)
			if err != nil {
				continue
			}
			if number >= targetNumber {
				return nil
			}
		}
	}
}

func (t *TaikoGethClient) Init(ctx context.Context) error {
	// Wait until taiko geth is ready.
	return t.EthIsReady(ctx, time.Second*20)
}

func (t *TaikoGethClient) Shutdown() error {
	return t.HiveManagedClient.Shutdown()
}

func (t *TaikoGethClient) HttpURL() string {
	return fmt.Sprintf("http://%v:%d", t.NetworkIP(), EthHttpPort)
}

func (t *TaikoGethClient) WSURL() string {
	return fmt.Sprintf("ws://%v:%d", t.NetworkIP(), EthWSPort)
}

func (t *TaikoGethClient) EngineURL() string {
	return fmt.Sprintf("http://%v:%d", t.NetworkIP(), EthEngineRPC)
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
