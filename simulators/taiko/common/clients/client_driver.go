package clients

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-node/p2p"
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

	p2pNode   *p2p.NodeP2P
	p2pSigner p2p.Signer

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

	d.p2pNode, d.p2pSigner, err = GetP2PNode(d.ctx, d.GetPreconfP2PNode())
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

func BuildPreconfRequestBody(
	ctx context.Context,
	rpccli *rpc.Client,
	privateKey *ecdsa.PrivateKey,
	l1Number, l2Number *big.Int,
) (*preconfblocks.BuildPreconfBlockRequestBody, error) {
	l1cli, l2cli := rpccli.L1, rpccli.L2

	l1Header, err := l1cli.HeaderByNumber(ctx, l1Number)
	if err != nil {
		return nil, fmt.Errorf("cannot get l1 header: %w", err)
	}

	var l2BlockID uint64
	if l2Number != nil {
		l2BlockID = l2Number.Uint64()
	} else {
		l2Header, err := l2cli.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot get l2 header: %w", err)
		}
		l2BlockID = l2Header.Number.Uint64()
	}

	signedTxs, err := utils.CreateL2Txs(context.Background(), l2cli, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed txs: %w", err)
	}

	parent, err := l2cli.HeaderByNumber(ctx, big.NewInt(0).SetUint64(l2BlockID-1))
	if err != nil {
		return nil, fmt.Errorf("cannot get parent block number, expect_number: %d: %v", l2BlockID-1, err)
	}

	preconfCfg, err := rpccli.GetProtocolConfigs(nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get protocol configs: %w", err)
	}

	baseFee, err := rpccli.CalculateBaseFee(
		ctx,
		parent,
		true,
		preconfCfg.BaseFeeConfig(),
		l1Header.Time,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate base fee: %w", err)
	}

	constructor, _ := anchorTxConstructor.New(rpccli)
	// Assemble a TaikoAnchor.anchorV3 transaction
	anchorTx, err := constructor.AssembleAnchorV3Tx(
		ctx,
		l1Header.Number,
		l1Header.Root,
		parent.GasUsed,
		preconfCfg.BaseFeeConfig(),
		[][32]byte{},
		new(big.Int).Add(parent.Number, common.Big1),
		baseFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct anchored tx: %w", err)
	}

	txBytes, err := utils.EncodeAndCompressTxList(append([]*types.Transaction{anchorTx}, signedTxs...))
	if err != nil {
		return nil, fmt.Errorf("failed to encode and compress anchor tx list: %w", err)
	}

	extraData := encoding.EncodeBaseFeeConfig(preconfCfg.BaseFeeConfig())
	reqBody := &preconfblocks.BuildPreconfBlockRequestBody{
		ExecutableData: &preconfblocks.ExecutableData{
			ParentHash:    parent.Hash(),
			FeeRecipient:  crypto.PubkeyToAddress(privateKey.PublicKey),
			Number:        l2BlockID,
			GasLimit:      uint64(preconfCfg.BlockMaxGasLimit()) + taiko.AnchorV3GasLimit,
			Timestamp:     l1Header.Time,
			Transactions:  txBytes,
			BaseFeePerGas: baseFee.Uint64(),
			ExtraData:     hexutil.Bytes(extraData[:]),
		},
	}

	return reqBody, nil
}

func BuildPreconfBlock(
	preconfURL string,
	requestBody *preconfblocks.BuildPreconfBlockRequestBody,
) (*types.Header, error) {

	// Try to propose a soft block with batch ID 0
	res, err := resty.New().
		R().
		SetBody(requestBody).
		Post(preconfURL + "/preconfBlocks")
	if err != nil {
		return nil, fmt.Errorf("failed to build preconf blocks: %w", err)
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to build preconf blocks: %s", res.String())
	}

	var body *preconfblocks.BuildPreconfBlockResponseBody
	if err = json.Unmarshal(res.Body(), &body); err != nil {
		return nil, err
	}

	return body.BlockHeader, nil
}
