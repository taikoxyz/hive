package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

	if ec.BlockNumber(ctx).Uint64() >= number {
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			ec.FailIfNotNil(ctx.Err())
		case <-time.After(timeout):
			ec.Fatalf("%s: wait latest number %d", ec.ClientType(), number)
		case <-ticker.C:
			l2Number := ec.BlockNumber(ctx)
			if l2Number.Uint64() >= number {
				return
			}
		}
	}
}

func WaitPreconfStatus(ctx context.Context, rpccli *rpc.Client, isPreconf bool, timeout time.Duration, number uint64) error {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for preconf status, isPreconf: %v, number: %d", isPreconf, number)
		case <-tick.C:
			l1Origin, err := rpccli.L2.L1OriginByID(ctx, big.NewInt(int64(number)))
			if err != nil {
				continue
			}
			if isPreconf && (l1Origin.L1BlockHeight == nil && l1Origin.L1BlockHash == (common.Hash{})) {
				return nil
			}
			if !isPreconf && (l1Origin.L1BlockHeight != nil && l1Origin.L1BlockHash != (common.Hash{})) {
				return nil
			}
		}
	}
}
