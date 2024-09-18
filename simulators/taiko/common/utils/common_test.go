package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"taiko/bindings/libproving"
	"taiko/bindings/taikol1"
	"taiko/params"
	"testing"
)

var (
	client   *ethclient.Client
	taikoL1  *taikol1.TaikoL1
	proverL1 *libproving.LibProving
)

func init() {
	url := "ws://localhost:8545"
	var err error
	client, err = ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	taikoL1, err = taikol1.NewTaikoL1(common.HexToAddress(params.EnvParams()["TAIKO_L1"]), client)
	if err != nil {
		panic(err)
	}
	proverL1, err = libproving.NewLibProving(common.HexToAddress(params.EnvParams()["TAIKO_L1"]), client)
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

func TestCC1(t *testing.T) {
	result, _ := taikoL1.State(nil)
	_ = result

	cfg, _ := taikoL1.GetConfig(nil)
	_ = cfg
	t.Log("OntakeForkHeight: ", cfg.OntakeForkHeight)

	for num := result.SlotB.LastVerifiedBlockId; num < result.SlotB.NumBlocks; num++ {
		blk, err := taikoL1.GetBlock(nil, num)
		if err != nil {
			t.Log(err.Error())
			break
		} else {
			t.Log(blk.BlockId)
		}
	}
}

func TestCC(t *testing.T) {
	tIter, _ := taikoL1.FilterBlockProposed(nil, nil, nil)
	tIterV2, _ := taikoL1.FilterBlockProposedV2(nil, nil)
	for tIter.Next() {
		t.Logf("v1 blockId: %d", tIter.Event.BlockId.Uint64())
		//data, _ := json.Marshal(tIter.Event.Meta)
		//t.Log("v1 content: ", string(data))
		t.Log("v1 L1Hash", common.BytesToHash(tIter.Event.Meta.L1Hash[:]))
		t.Log("---------------------------------------------------------------------------------------------------")
	}
	for tIterV2.Next() {
		t.Logf("v2 blockId: %d", tIterV2.Event.BlockId.Uint64())
		//data, _ := json.Marshal(tIterV2.Event.Meta)
		//t.Log("v2 content: ", string(data))
		t.Log("v2 AnchorBlockHash: ", common.BytesToHash(tIterV2.Event.Meta.AnchorBlockHash[:]))
		t.Log("---------------------------------------------------------------------------------------------------")
	}
}

func TestCC22(t *testing.T) {
	pIter, _ := proverL1.FilterLibProvingData(nil, nil)
	pIterV2, _ := proverL1.FilterLibProvingDataV2(nil, nil)
	for pIter.Next() {
		//metaHash1 := common.BytesToHash(pIter.Event.MetaHash1[:])
		//metaHash2 := common.BytesToHash(pIter.Event.MetaHash2[:])
		//if metaHash1 != metaHash2 {
		t.Logf("v1 blockId: %d", pIter.Event.BlockId.Uint64())
		//data, _ := json.Marshal(pIter.Event.Meta)
		//t.Logf("v1 content: %s", string(data))
		t.Log("v1 L1Hash", common.BytesToHash(pIter.Event.Meta.L1Hash[:]))
		t.Log("---------------------------------------------------------------------------------------------------")
		//}
	}
	for pIterV2.Next() {
		//metaHash1 := common.BytesToHash(pIterV2.Event.MetaHash1[:])
		//metaHash2 := common.BytesToHash(pIterV2.Event.MetaHash2[:])
		//if metaHash1 != metaHash2 {
		t.Logf("v2 blockId: %d", pIterV2.Event.BlockId.Uint64())
		//data, _ := json.Marshal(pIterV2.Event.Meta)
		//t.Logf("v2 content: %s", string(data))
		t.Log("v2 AnchorBlockHash: ", common.BytesToHash(pIterV2.Event.Meta.AnchorBlockHash[:]))
		t.Log("---------------------------------------------------------------------------------------------------")
		//}
	}
}
