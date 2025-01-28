package preconf

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/taiko"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"math/big"
	"strings"
	"taiko/common/clients"
	"testing"
)

func TestNewProposer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	proposerClient := &proposer.Proposer{}
	err := clients.NewTaikoClient(proposerClient, flags.ProposerFlags)
	assert.Nil(t, err)
	assert.NotNil(t, proposerClient)

	rpcClient, err := rpc.NewClient(ctx, proposerClient.ClientConfig)
	assert.Nil(t, err)

	var num hexutil.Uint
	err = rpcClient.L2.Call(&num, "eth_getBlockTransactionCountByNumber", "pending")
	assert.Nil(t, err)
	t.Log(num)

	pendingCount, err := rpcClient.L2.PendingTransactionCount(ctx)
	assert.Nil(t, err)
	t.Log(pendingCount)
}

func TestCC(t *testing.T) {
	signer := types.NewCancunSigner(big.NewInt(167001))

	signed := `{"type":"0x2","chainId":"0x28c59","nonce":"0x8","to":"0x1670010000000000000000000000000000010001","gas":"0x3d090","gasPrice":null,"maxPriorityFeePerGas":"0x0","maxFeePerGas":"0x882a2e","value":"0x0","input":"0xfd85eb2d000000000000000000000000000000000000000000000000000000000000006f2a408ffc86b00dc24157f7679b45a184d68c567adb03ff6ea944795817277510000000000000000000000000000000000000000000000000000000000005416a0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000004b00000000000000000000000000000000000000000000000000000000004c4b40000000000000000000000000000000000000000000000000000000004fdec7000000000000000000000000000000000000000000000000000000000023c34600","accessList":[],"v":"0x1","r":"0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798","s":"0x67cb0ec335072349c59627a8db16b599da349d6195df213b33d230fc51693eca","yParity":"0x1","hash":"0x45ff24b15e13b9ba3b23c02f43b5b1832e32de903eb0f32bef9e8933747a1bfb"}`

	var tx types.Transaction
	err := json.Unmarshal([]byte(signed), &tx)
	assert.Nil(t, err)

	addr, err := signer.Sender(&tx)
	assert.Nil(t, err)
	t.Log(addr.String() == taiko.GoldenTouchAccount.String())
}

type Taiko struct {
	taikoL2Address common.Address
}

func (t *Taiko) ValidateAnchorTx(tx *types.Transaction, header *types.Header) (bool, error) {
	if tx.Type() != types.DynamicFeeTxType {
		return false, nil
	}

	if tx.To() == nil || *tx.To() != t.taikoL2Address {
		return false, nil
	}

	if !bytes.HasPrefix(tx.Data(), taiko.AnchorSelector) && !bytes.HasPrefix(tx.Data(), taiko.AnchorV2Selector) {
		return false, nil
	}

	if tx.Value().Cmp(common.Big0) != 0 {
		return false, nil
	}

	if tx.Gas() != taiko.AnchorGasLimit {
		return false, nil
	}

	//if tx.GasFeeCap().Cmp(header.BaseFee) != 0 {
	//	return false, nil
	//}

	s := types.NewCancunSigner(big.NewInt(167001))

	addr, err := s.Sender(tx)
	if err != nil {
		return false, err
	}

	return strings.EqualFold(addr.String(), taiko.GoldenTouchAccount.String()), nil
}

func TestBB(t *testing.T) {
	t.Log(common.Bytes2Hex(taiko.AnchorSelector))
	t.Log(common.Bytes2Hex(taiko.AnchorV2Selector))

	taikoL2AddressPrefix := strings.TrimPrefix(big.NewInt(167001).String(), "0")

	tk := &Taiko{
		taikoL2Address: common.HexToAddress(
			"0x" +
				taikoL2AddressPrefix +
				strings.Repeat("0", common.AddressLength*2-len(taikoL2AddressPrefix)-len(taiko.TaikoL2AddressSuffix)) +
				taiko.TaikoL2AddressSuffix,
		),
	}

	signed := `{"type":"0x2","chainId":"0x28c59","nonce":"0x8","to":"0x1670010000000000000000000000000000010001","gas":"0x3d090","gasPrice":null,"maxPriorityFeePerGas":"0x0","maxFeePerGas":"0x882a2e","value":"0x0","input":"0xfd85eb2d000000000000000000000000000000000000000000000000000000000000006fd1ba04a46072faf71135060dc9f756ec9df566e0b3b7cd363e8557aacb7efccc000000000000000000000000000000000000000000000000000000000005416a0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000004b00000000000000000000000000000000000000000000000000000000004c4b40000000000000000000000000000000000000000000000000000000004fdec7000000000000000000000000000000000000000000000000000000000023c34600","accessList":[],"v":"0x1","r":"0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798","s":"0x653ccceb6cf54f04e04bbd6df2afaae52490806336c6c89dbbdc1a4deb0205ca","yParity":"0x1","hash":"0xe8dd0176cfc56c55aa2d9a1a6f5ce4c201bbd1eb191ab33432340f5225212f1c"}`

	var tx types.Transaction
	err := json.Unmarshal([]byte(signed), &tx)
	assert.Nil(t, err)

	ok, err := tk.ValidateAnchorTx(&tx, nil)
	assert.Nil(t, err)

	t.Log(ok)
}

func TestDD(t *testing.T) {
	l2Cli, err := ethclient.Dial("http://localhost:8535")
	assert.Nil(t, err)

	_, err = l2Cli.HeaderByNumber(context.Background(), nil)
	assert.Nil(t, err)

	count, err := l2Cli.PendingTransactionCount(context.Background())
	assert.Nil(t, err)
	t.Log(count)
}
