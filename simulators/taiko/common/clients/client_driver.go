package clients

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/taiko"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/hive/hivesim"
	"github.com/go-resty/resty/v2"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	anchorTxConstructor "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/anchor_tx_constructor"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"taiko/common/utils"
)

type DriverClient struct {
	Index int
	*HiveManagedClient
	Envs hivesim.Params
	*State
}

func (d *DriverClient) Start() (err error) {
	if err := d.HiveManagedClient.Start(); err != nil {
		return err
	}

	d.Logf("driver client, L1_BEACON: %s", d.Envs["L1_BEACON"])

	client, err := rpc.NewClient(context.Background(), GetClientConfig(d.Envs))
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
	err := d.HiveManagedClient.Shutdown()
	d.FailIfNotNil(err, fmt.Sprintf("failed to shutdown %s", d.ClientType()))
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
	privateKey *ecdsa.PrivateKey,
	preconfURL string,
	anchoredL1Block *types.Header,
	l2BlockID uint64,
) (*types.Header, types.Transactions, error) {
	l2cli := rpccli.L2

	signedTxs, err := utils.CreateL2Txs(context.Background(), l2cli, true)
	if err != nil {
		return nil, nil, err
	}

	parent, err := l2cli.HeaderByNumber(ctx, big.NewInt(0).SetUint64(l2BlockID-1))
	if err != nil {
		return nil, nil, err
	}

	preconfCfg, err := rpccli.GetProtocolConfigs(nil)
	if err != nil {
		return nil, nil, err
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

	txBytes, err := utils.EncodeAndCompressTxList(append([]*types.Transaction{anchorTx}, signedTxs...))
	if err != nil {
		return nil, nil, err
	}

	extraData := encoding.EncodeBaseFeeConfig(preconfCfg.BaseFeeConfig())
	reqBody := &preconfblocks.BuildPreconfBlockRequestBody{
		ExecutableData: &preconfblocks.ExecutableData{
			ParentHash:    parent.Hash(),
			FeeRecipient:  crypto.PubkeyToAddress(privateKey.PublicKey),
			Number:        l2BlockID,
			GasLimit:      uint64(preconfCfg.BlockMaxGasLimit()) + taiko.AnchorV3GasLimit,
			Timestamp:     anchoredL1Block.Time,
			Transactions:  txBytes,
			BaseFeePerGas: baseFee.Uint64(),
			ExtraData:     hexutil.Bytes(extraData[:]),
		},
	}

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

	return body.BlockHeader, signedTxs, nil
}
