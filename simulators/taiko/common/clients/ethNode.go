package clients

import (
	"context"
	"fmt"
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

	//OntakeL1 *ontake.OntakeL1Clients
	//PacayaL1 *pacaya.PacayaL1Clients
	//OntakeL2 *ontake.OntakeL2Clients
	//PacayaL2 *pacaya.PacayaL2Clients
}

func (ec *EthNode) Start() (err error) {
	if !ec.IsRunning() {
		if err := ec.HiveManagedClient.Start(); err != nil {
			return err
		}
	}

	// try 20 times until the eth client is connected.
	for times := 0; times < 20; times++ {
		ec.EthClient, err = rpc.NewEthClient(context.Background(), ec.HttpURL(), time.Second)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	return err
}

/*
func (ec *EthNode) InitL1Clients() (err error) {
	ec.OntakeL1, err = ontake.NewOntakeL1Clients(ec.EthClient)
	if err != nil {
		return err
	}
	ec.PacayaL1, err = pacaya.NewPacayaL1Clients(ec.EthClient)
	return err
}

func (ec *EthNode) InitL2Clients() (err error) {
	ec.OntakeL2, err = ontake.NewOntakeL2Clients(ec.EthClient)
	if err != nil {
		return err
	}
	ec.PacayaL2, err = pacaya.NewPacayaL2Clients(ec.EthClient)
	return err
}*/

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
