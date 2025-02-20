package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"taiko/params"
	"testing"
)

var (
	rpccli       *rpc.Client
	l1Cli, l2Cli *rpc.EthClient
	l1Ontake     *ontake.OntakeL1Clients
	l1Pacaya     *pacaya.PacayaL1Clients
)

func init() {
	var err error
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

	l1Cli, l2Cli = rpccli.L1, rpccli.L2
	l1Ontake, err = ontake.NewOntakeL1Clients(l1Cli)
	if err != nil {
		panic(err)
	}

	l1Pacaya, err = pacaya.NewPacayaL1Clients(l2Cli)
	if err != nil {
		panic(err)
	}
}

func TestVerify(t *testing.T) {
	l2Header, err := l2Cli.HeaderByNumber(context.Background(), nil)
	assert.NoError(t, err)

	tx, err := l1Ontake.TaikoL1.VerifyBlocks(params.L1Auths[0], 32)
	assert.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), l1Cli, tx)
	assert.NoError(t, err)
	t.Log(receipt.Status)

	if l2Header.Number.Uint64() < 10 {
		tx, err := l1Ontake.TaikoL1.VerifyBlocks(params.L1Auths[0], 32)
		assert.NoError(t, err)
		receipt, err := bind.WaitMined(context.Background(), l1Cli, tx)
		assert.NoError(t, err)
		t.Log(receipt.Status)
	} else {
		tx, err := l1Pacaya.TaikoInbox.VerifyBatches(params.L1Auths[0], 32)
		assert.NoError(t, err)
		receipt, err := bind.WaitMined(context.Background(), l1Cli, tx)
		assert.NoError(t, err)
		t.Log(receipt.Status)
	}
}

func TestCC(t *testing.T) {
	st, err := l1Ontake.TaikoL1.State(nil)
	assert.NoError(t, err)
	t.Log(st.SlotB.LastVerifiedBlockId)
}

func TestPreconfAPI(t *testing.T) {
	l2Block, err := l2Cli.BlockByNumber(context.Background(), nil)
	assert.NoError(t, err)

	anchorL1Header, err := l1Cli.HeaderByNumber(context.Background(), nil)
	assert.NoError(t, err)

	header, txs, err := BuildPreconfBlock(
		context.Background(),
		rpccli,
		fmt.Sprintf("http://localhost:%d", PreconfServerPort),
		anchorL1Header,
		l2Block.NumberU64(),
		nil, //[]*types.Transaction{l2Block.Transactions()[0]},
	)
	assert.NoError(t, err)
	t.Log(header.Hash().String())

	l2Block, err = l2Cli.BlockByHash(context.Background(), header.Hash())
	assert.NoError(t, err)

	assert.Equalf(t, l2Block.Transactions().Len(), len(txs), "transaction length not match")

	txsMap := make(map[common.Hash]*types.Transaction)
	for _, tx := range txs {
		txsMap[tx.Hash()] = tx
	}

	for _, tx := range l2Block.Transactions() {
		assert.Equalf(t, true, txsMap[tx.Hash()] != nil, fmt.Sprintf("tx %s not found", tx.Hash().String()))
	}
}

func TestCC1(t *testing.T) {
	t.Log(os.Getenv("L1_PROPOSER_PRIV_KEY"))
	preconferPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIV_KEY")))
	assert.NoError(t, err)

	payload, err := rlp.EncodeToBytes(&preconfblocks.BuildPreconfBlockRequestBody{})
	assert.NoError(t, err)

	hash := crypto.Keccak256(payload)

	sign, err := crypto.Sign(hash, preconferPrivKey)
	assert.NoError(t, err)

	pubKey, err := crypto.Ecrecover(hash, sign)
	assert.NoError(t, err)
	t.Log(common.Bytes2Hex(crypto.FromECDSAPub(&preconferPrivKey.PublicKey)))
	t.Log(common.Bytes2Hex(pubKey))

	isValid := crypto.VerifySignature(pubKey, hash, sign[:64])
	t.Log(isValid)

	pk, err := crypto.UnmarshalPubkey(pubKey)
	assert.NoError(t, err)

	t.Log(crypto.PubkeyToAddress(*pk).Hex())
}
