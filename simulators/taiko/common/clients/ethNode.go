package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	client     *rpc.Client
}

func (ec *EthNode) Start() error {
	if !ec.IsRunning() {
		if err := ec.HiveManagedClient.Start(); err != nil {
			return err
		}
	}

	return ec.EthIsReady(context.Background(), time.Second*20)
}

func (ec *EthNode) RPClient() *rpc.Client {
	if ec.client == nil {
		var err error
		ec.client, err = rpc.Dial(ec.HttpURL())
		if err != nil {
			ec.Fatalf("failed to dial %s node", ec.ClientType())
		}
	}
	return ec.client
}

func (ec *EthNode) EthClient() *ethclient.Client {
	return ethclient.NewClient(ec.RPClient())
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
	ethClient := ec.EthClient()
	for ; ; <-time.Tick(time.Second) {
		ec.Logf("waiting for %s to be ready", ec.ClientType())
		select {
		case <-time.After(timeout):
			return fmt.Errorf("reach timeout but l1geth is not ready")
		default:
			_, err := ethClient.ChainID(ctx)
			if err != nil {
				continue
			}
			return nil
		}
	}
}

func (ec *EthNode) WaitNumber(ctx context.Context, timeout time.Duration, targetNumber uint64) error {
	ethClient := ec.EthClient()

	current, times := uint64(0), timeout/time.Second
	for times > 0 && targetNumber > current {
		select {
		case <-time.Tick(time.Second):
			number, err := ethClient.BlockNumber(ctx)
			if err != nil {
				ec.Logf("failed to get block number, err: %v", err)
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

	if targetNumber > current {
		return fmt.Errorf("failed to reach current number %d, current number: %d", targetNumber, current)
	}
	return nil
}
