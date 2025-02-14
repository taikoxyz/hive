package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"taiko/params"
	"testing"
)

var (
	rpccli   *rpc.Client
	l1Cli    *ethclient.Client
	l2Cli    *ethclient.Client
	l1Ontake *ontake.OntakeL1Clients
	l1Pacaya *pacaya.PacayaL1Clients
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

	l1Ontake, err = ontake.NewOntakeL1Clients(l1Cli)
	if err != nil {
		panic(err)
	}

	l1Pacaya, err = pacaya.NewPacayaL1Clients(l1Cli)
	if err != nil {
		panic(err)
	}

	rpccli, err = rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:        "ws://localhost:8545",
		L2Endpoint:        "ws://localhost:6046",
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:  "http://localhost:6051",
		JwtSecret:         "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
	})
	if err != nil {
		panic(err)
	}
}

func TestVerify(t *testing.T) {
	tx, err := l1Ontake.TaikoL1.VerifyBlocks(params.L1Auths[0], 16)
	assert.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), l1Cli, tx)
	assert.NoError(t, err)
	t.Log(receipt.Status)
}

func TestCC(t *testing.T) {
	st, err := l1Ontake.TaikoL1.State(nil)
	assert.NoError(t, err)
	t.Log(st.SlotB.LastVerifiedBlockId)
}

func TestPreconfAPI(t *testing.T) {
	l2Number, err := l2Cli.BlockNumber(context.Background())
	assert.NoError(t, err)

	anchorL1Header, err := l1Cli.HeaderByNumber(context.Background(), nil)
	assert.NoError(t, err)

	header, txs, err := buildPreconfBlock(context.Background(), rpccli, fmt.Sprintf(""), anchorL1Header, l2Number)
	assert.NoError(t, err)

	l1Block, err := l1Cli.BlockByHash(context.Background(), header.Hash())
	assert.NoError(t, err)

	assert.Equalf(t, l1Block.Transactions().Len(), len(txs), "transaction length not match")

	txsMap := make(map[common.Hash]*types.Transaction)
	for _, tx := range txs {
		txsMap[tx.Hash()] = tx
	}

	for _, tx := range l1Block.Transactions() {
		assert.Equalf(t, true, txsMap[tx.Hash()] != nil, fmt.Sprintf("tx %s not found", tx.Hash().String()))
	}
}
