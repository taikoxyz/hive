package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestDeployContracts(t *testing.T) {
	url := "http://localhost:8545"
	assert.NoError(t, DeployContracts(context.Background(), url))
	assert.NoError(t, setIntervalMining(url, 3))
}
