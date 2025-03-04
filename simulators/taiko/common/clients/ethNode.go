package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"time"
)

type EthExposeAPI interface {
	ClientType() string
	HttpURL() string
	WSURL() string
	EngineURL() string
	HTTPClient() *rpc.EthClient
}

type EthNode struct {
	*HiveManagedClient

	HttpPort   int64
	WSPort     int64
	EnginePort int64
	EthClient  *rpc.EthClient
}

func (ec *EthNode) Start() (err error) {
	if !ec.IsRunning() {
		if err := ec.HiveManagedClient.Start(); err != nil {
			return err
		}
	}

	// Try 20 times until the eth client is ready to connect.
	for times := 0; times < 20; times++ {
		ec.EthClient, err = rpc.NewEthClient(context.Background(), ec.WSURL(), time.Second)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	return err
}

func (ec *EthNode) HTTPClient() *rpc.EthClient {
	return ec.EthClient
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

func (ec *EthNode) WaitLatestNumber(ctx context.Context, timeout time.Duration, number uint64) {
	ec.Logf("%s: wait latest number %d", ec.ClientType(), number)

	cli := ethclient.NewClient(ec.EthClient.Client)

	high, err := cli.BlockNumber(ctx)
	ec.FailIfNotNil(err, "failed to get latest number")

	if high >= number {
		return
	}

	headerCh := make(chan *types.Header, 10)
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	sub, err := cli.SubscribeNewHead(subCtx, headerCh)
	ec.FailIfNotNil(err, "failed to subscribe to new head")
	defer sub.Unsubscribe()

	for header := range headerCh {
		if header.Number.Uint64() >= number {
			return
		}
	}

	ec.Fatalf("failed to wait latest number %d", number)
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
