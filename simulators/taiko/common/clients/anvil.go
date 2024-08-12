package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"taiko/common/utils"
	"time"
)

var (
	AnvilPort int64 = 8545
)

type AnvilClient struct {
	*HiveManagedClient
	Logger    utils.Logging
	Network   string
	networkIP string
}

func (ec *AnvilClient) Start() error {
	ec.Logf("Starting taiko geth client")
	if !ec.IsRunning() {
		return ec.HiveManagedClient.Start()
	}

	return ec.Init(context.Background())
}

func (ec *AnvilClient) Init(ctx context.Context) error {
	// Wait until taiko geth is ready.
	return ec.EthIsReady(ctx, time.Second*20)
}

func (ec *AnvilClient) Shutdown() error {
	return ec.HiveManagedClient.Shutdown()
}

func (ec *AnvilClient) Logf(format string, values ...interface{}) {
	if l := ec.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (ec *AnvilClient) HttpURL() string {
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), AnvilPort)
}

func (ec *AnvilClient) WSURL() string {
	return fmt.Sprintf("ws://%v:%d", ec.NetworkIP(), AnvilPort)
}

func (ec *AnvilClient) NetworkIP() string {
	if ec.networkIP != "" {
		return ec.networkIP
	}

	t := ec.T
	var err error
	ec.networkIP, err = t.Sim.ContainerNetworkIP(t.SuiteID, ec.Network, ec.Client.Container)
	if err != nil {
		t.Logf("Error getting network IP: %v", err)
		return ec.HiveManagedClient.GetHost()
	}
	t.Logf("execution network IP: %s", ec.networkIP)

	return ec.networkIP
}

func (ec *AnvilClient) EthIsReady(ctx context.Context, timeout time.Duration) error {
	client, err := ethclient.DialContext(ctx, ec.HttpURL())
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

func (ec *AnvilClient) VerifyNumber(ctx context.Context, timeout time.Duration, targetNumber uint64) error {
	client, err := ethclient.DialContext(ctx, ec.HttpURL())
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
				ec.Logf("failed to get block number, err: %v", err)
				continue
			}
			if number >= targetNumber {
				return nil
			}
		}
	}
}
