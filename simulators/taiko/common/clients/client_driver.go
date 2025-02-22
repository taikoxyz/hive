package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/taiko"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-resty/resty/v2"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	anchorTxConstructor "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/anchor_tx_constructor"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"os"
	"taiko/common/utils"
)

type DriverClient struct {
	*HiveManagedClient

	*State
}

func (d *DriverClient) Start() (err error) {
	if err := d.HiveManagedClient.Start(); err != nil {
		return err
	}

	d.Logf("driver client, L1_BEACON: %s", os.Getenv("L1_BEACON"))

	client, err := rpc.NewClient(context.Background(), GetClientConfig())
	if err != nil {
		return err
	}

	d.State, err = NewState(client)
	if err != nil {
		return err
	}

	return err
}

func (d *DriverClient) Shutdown() error {
	if err := d.HiveManagedClient.Shutdown(); err != nil {
		return err
	}
	if d.State != nil {
		d.State.Close()
	}

	return nil
}

func (d *DriverClient) PreconfServerURL() string {
	return fmt.Sprintf("http://%s:%v", d.NetworkIP(), PreconfServerPort)
}

func BuildPreconfBlock(
	ctx context.Context,
	rpccli *rpc.Client,
	preconfURL string,
	anchoredL1Block *types.Header,
	l2BlockID uint64,
	txs types.Transactions,
) (*types.Header, types.Transactions, error) {
	l2cli := rpccli.L2

	parent, err := l2cli.HeaderByNumber(ctx, big.NewInt(0).SetUint64(l2BlockID-1))
	if err != nil {
		return nil, nil, err
	}

	preconfCfg, err := rpccli.GetProtocolConfigs(nil)
	if err != nil {
		return nil, nil, err
	}

	// Create and send a batch of txs.
	if txs == nil {

		signedTxs, err := utils.CreateL2Txs(context.Background(), l2cli, true)
		if err != nil {
			return nil, nil, err
		}
		txs = signedTxs
	}

	baseFee, err := rpccli.CalculateBaseFee(
		ctx,
		parent,
		true,
		preconfCfg.BaseFeeConfig(),
		anchoredL1Block.Time,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate base fee: %w", err)
	}

	constructor, _ := anchorTxConstructor.New(rpccli)
	// Assemble a TaikoAnchor.anchorV3 transaction
	anchorTx, err := constructor.AssembleAnchorV3Tx(
		ctx,
		anchoredL1Block.Number,
		anchoredL1Block.Root,
		parent.GasUsed,
		preconfCfg.BaseFeeConfig(),
		[][32]byte{},
		new(big.Int).Add(parent.Number, common.Big1),
		baseFee,
	)
	if err != nil {
		return nil, nil, err
	}

	txBytes, err := utils.EncodeAndCompressTxList(append([]*types.Transaction{anchorTx}, txs...))
	if err != nil {
		return nil, nil, err
	}

	preconferPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIV_KEY")))
	if err != nil {
		return nil, nil, err
	}

	extraData := encoding.EncodeBaseFeeConfig(preconfCfg.BaseFeeConfig())
	reqBody := &preconfblocks.BuildPreconfBlockRequestBody{
		ExecutableData: &preconfblocks.ExecutableData{
			ParentHash:    parent.Hash(),
			FeeRecipient:  crypto.PubkeyToAddress(preconferPrivKey.PublicKey),
			Number:        l2BlockID,
			GasLimit:      uint64(preconfCfg.BlockMaxGasLimit()) + taiko.AnchorV3GasLimit,
			Timestamp:     anchoredL1Block.Time,
			Transactions:  txBytes,
			BaseFeePerGas: baseFee.Uint64(),
			ExtraData:     hexutil.Bytes(extraData[:]),
		},
	}

	payload, err := rlp.EncodeToBytes(reqBody.ExecutableData)
	if err != nil {
		return nil, nil, err
	}

	hash := crypto.Keccak256(payload)
	sig, err := crypto.Sign(hash, preconferPrivKey)
	if err != nil {
		return nil, nil, err
	}
	reqBody.Signature = common.Bytes2Hex(sig)

	// Try to propose a soft block with batch ID 0
	res, err := resty.New().
		R().
		SetBody(reqBody).
		Post(preconfURL + "/preconfBlocks")
	if err != nil {
		return nil, nil, err
	}
	if !res.IsSuccess() {
		return nil, nil, fmt.Errorf("failed to build preconf block: %v", res.Error())
	}

	var body *preconfblocks.BuildPreconfBlockResponseBody
	if err = json.Unmarshal(res.Body(), &body); err != nil {
		return nil, nil, err
	}

	return body.BlockHeader, txs, nil
}
