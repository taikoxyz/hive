package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	l1Cli *ethclient.Client
	l2Cli *ethclient.Client
)

func init() {
	var err error
	l1Cli, err = ethclient.Dial("ws://localhost:8545")
	if err != nil {
		panic(err)
	}
	l2Cli, err = ethclient.Dial("ws://localhost:6046")
	if err != nil {
		panic(err)
	}
}

func setIntervalMining(url string, interval int) error {
	client, err := rpc.Dial(url)
	if err != nil {
		return err
	}
	return client.CallContext(context.Background(), nil, "evm_setIntervalMining", interval)
}

func setL1Automine(url string, automine bool) error {
	client, err := rpc.Dial(url)
	if err != nil {
		return err
	}
	return client.CallContext(context.Background(), nil, "evm_setAutomine", automine)
}

func increaseTime(timestamp uint64) error {
	cli := l1Cli.Client()
	err := cli.CallContext(context.Background(), nil, "evm_increaseTime", timestamp)
	return err
}
