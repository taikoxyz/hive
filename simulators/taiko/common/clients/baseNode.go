package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

type BaseNode struct {
	*HiveManagedClient

	HttpPort int64
	WSPort   int64
}

func (ec *BaseNode) Start() error {
	if !ec.IsRunning() {
		return ec.HiveManagedClient.Start()
	}

	return ec.Init(context.Background())
}

func (ec *BaseNode) Init(ctx context.Context) error {
	// Wait until taiko geth is ready.
	return ec.EthIsReady(ctx, time.Second*20)
}

func (ec *BaseNode) HttpURL() string {
	if ec.HttpPort == 0 {
		ec.HttpPort = EthHttpPort
	}
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), ec.HttpPort)
}

func (ec *BaseNode) WSURL() string {
	if ec.WSPort == 0 {
		ec.WSPort = EthWSPort
	}
	return fmt.Sprintf("ws://%v:%d", ec.NetworkIP(), ec.WSPort)
}

func (ec *BaseNode) EthIsReady(ctx context.Context, timeout time.Duration) error {
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

func (ec *BaseNode) VerifyNumber(ctx context.Context, timeout time.Duration, targetNumber uint64) error {
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
