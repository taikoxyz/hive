package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

type EthExposeAPI interface {
	NetworkIP() string
	HttpURL() string
	WSURL() string
}

type EthNode struct {
	*HiveManagedClient

	HttpPort   int64
	WSPort     int64
	EnginePort int64
	EthClient  *ethclient.Client
}

func (ec *EthNode) Start() (err error) {
	if !ec.IsRunning() {
		if err := ec.HiveManagedClient.Start(); err != nil {
			return err
		}
	}

	// Connect to eth client
	ec.EthClient, err = ethclient.Dial(ec.HttpURL())
	if err != nil {
		return err
	}

	return ec.EthIsReady(context.Background(), time.Second*20)
}

func (ec *EthNode) HttpURL() string {
	if ec.HttpPort == 0 {
		ec.HttpPort = EthHttpPort
	}
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), ec.HttpPort)
}

func (ec *EthNode) WSURL() string {
	if ec.WSPort == 0 {
		ec.WSPort = EthWSPort
	}
	return fmt.Sprintf("ws://%v:%d", ec.NetworkIP(), ec.WSPort)
}

func (ec *EthNode) EngineURL() string {
	if ec.EnginePort == 0 {
		ec.EnginePort = EthEngineRPC
	}
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), ec.EnginePort)
}

func (ec *EthNode) EthIsReady(ctx context.Context, timeout time.Duration) error {
	for ; ; <-time.Tick(time.Second) {
		ec.Logf("waiting for %s to be ready", ec.ClientType())
		select {
		case <-time.After(timeout):
			return fmt.Errorf("reach timeout but %s is not ready", ec.ClientType())
		default:
			_, err := ec.EthClient.ChainID(ctx)
			if err != nil {
				continue
			}
			return nil
		}
	}
}

func (ec *EthNode) WaitLatestNumber(ctx context.Context, timeout time.Duration, number uint64) error {
	ec.Logf("%s: wait latest number %d", ec.ClientType(), number)
	current, times := uint64(0), timeout/time.Second
	for times > 0 && number >= current {
		select {
		case <-time.Tick(time.Second):
			number, err := ec.EthClient.BlockNumber(ctx)
			if err != nil {
				ec.Logf("failed to get block number from %s, err: %v", ec.ClientType(), err)
				continue
			}
			if number >= current {
				current = number + 1
				times = timeout / time.Second
				break
			} else {
				times--
			}
		}
	}

	if number >= current {
		return fmt.Errorf("%s failed to reach current number %d, current number: %d", ec.ClientType(), number, current)
	}
	return nil
}

func (ec *EthNode) WaitTargetNumber(ctx context.Context, timeout time.Duration, number uint64) error {
	defer ec.Logf("%s: wait target number %d", ec.ClientType(), number)
	for ; ; <-time.Tick(time.Second) {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("reach timeout but %s is not ready", ec.ClientType())
		default:
			_, err := ec.EthClient.HeaderByNumber(ctx, new(big.Int).SetUint64(number))
			if err == nil {
				return nil
			}
		}
	}
}
