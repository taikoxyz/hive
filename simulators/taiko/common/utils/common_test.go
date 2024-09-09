package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"taiko/bindings/guardianprover"
	"taiko/bindings/libproving"
	"taiko/bindings/taikol1"
	"taiko/params"
	"testing"
)

var client *ethclient.Client

func init() {
	url := "ws://localhost:8545"
	var err error
	client, err = ethclient.Dial(url)
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
	cli := client.Client()
	err := cli.CallContext(context.Background(), nil, "evm_increaseTime", timestamp)
	return err
}

func TestAA(t *testing.T) {
	assert.NoError(t, increaseTime(1200))
}

func mineBlock(client *rpc.Client) error {
	return client.CallContext(context.Background(), nil, "evm_mine")
}

func TestDeployContracts(t *testing.T) {
	url := "http://localhost:8545"
	assert.NoError(t, DeployContracts(context.Background(), url))
	assert.NoError(t, setIntervalMining(url, 3))
}

func TestCC1(t *testing.T) {
	taikoL1, err := taikol1.NewTaikoL1(common.HexToAddress(params.EnvParams()["TAIKO_L1"]), client)
	assert.NoError(t, err)

	id, err := taikoL1.GetTransitionId(nil)
	assert.NoError(t, err)

	t.Log(id)
}

func TestCC(t *testing.T) {
	proving, err := libproving.NewLibProving(common.HexToAddress(params.EnvParams()["TAIKO_L1"]), client)
	assert.NoError(t, err)

	iter2, err := proving.FilterTransitionProvedV2(nil, nil)
	assert.NoError(t, err)

	t.Log(iter2.Next())
}

func TestTaikoL1(t *testing.T) {
	tl1, err := taikol1.NewTaikoL1(common.HexToAddress(params.EnvParams()["PROVER_SET"]), client)
	assert.NoError(t, err)
	out, err := tl1.GetTransitionId(nil)
	assert.NoError(t, err)
	t.Log(out)

	iter, err := tl1.FilterTransitionProved(nil, nil)
	assert.NoError(t, err)
	t.Log(iter.Next())

	iter2, err := tl1.FilterTransitionProvedV2(nil, nil)
	assert.NoError(t, err)
	t.Log(iter2.Next())

	prover, err := guardianprover.NewGuardianProver(common.HexToAddress(params.EnvParams()["GUARDIAN_PROVER_MINORITY"]), client)
	assert.NoError(t, err)

	minGuardians, err := prover.MinGuardians(nil)
	assert.NoError(t, err)
	t.Log("minGuardians: ", minGuardians)

	sink := make(chan *guardianprover.GuardianProverApproved, 3)
	sub, err := prover.WatchApproved(nil, sink, nil)
	assert.NoError(t, err)

	for {
		select {
		case result := <-sink:
			t.Log("approvalBits: ", result.ApprovalBits, result.MinGuardiansReached)
		case <-sub.Err():
			break
		}
	}
}
