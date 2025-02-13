package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"taiko/params"
	"testing"
)

var (
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
}

func TestVerify(t *testing.T) {
	tx, err := l1Ontake.TaikoL1.VerifyBlocks(params.L1Auths[0], 16)
	assert.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), l1Cli, tx)
	assert.NoError(t, err)
	t.Log(receipt.Status)
}

func TestGetBlock(t *testing.T) {
	st, err := l1Ontake.TaikoL1.State(nil)
	assert.NoError(t, err)
	t.Log(st.SlotB.LastVerifiedBlockId)

	blockId := st.SlotB.LastVerifiedBlockId

	block, err := l1Ontake.TaikoL1.GetBlockV2(nil, blockId)
	assert.NoError(t, err)
	t.Log(block)

	tid := block.VerifiedTransitionId.Uint64()

	ts, err := l1Ontake.TaikoL1.GetTransition(nil, blockId, uint32(tid))
	assert.NoError(t, err)

	tid1, err := l1Ontake.TaikoL1.GetTransitionById(nil, blockId, ts.BlockHash)
	assert.NoError(t, err)

	t.Log(tid1.Uint64())
}

func TestCC(t *testing.T) {
	st, err := l1Ontake.TaikoL1.State(nil)
	assert.NoError(t, err)
	t.Log(st.SlotB.LastVerifiedBlockId)

	parent, err := l2Cli.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(st.SlotB.LastVerifiedBlockId))
	assert.NoError(t, err)

	res, err := l1Ontake.TaikoL1.GetTransition0(&bind.CallOpts{Context: context.Background()}, st.SlotB.LastVerifiedBlockId+1, parent.Hash())
	assert.NoError(t, err)
	t.Log(res)
}

func TestPacaya(t *testing.T) {
	st, err := l1Pacaya.TaikoInbox.State(nil)
	assert.NoError(t, err)
	t.Log(st.Stats2.LastVerifiedBatchId)
}
