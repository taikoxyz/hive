package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-resty/resty/v2"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"os"
	"taiko/common/utils"
	"taiko/params"
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

func (d *DriverClient) BuildPreconfBlock(
	l2BlockID uint64,
	l1Head *types.Header,
) (*types.Header, types.Transactions, error) {
	d.Logf("%s: build soft block", d.ClientType())

	if l1Head == nil {
		l1Head = d.L1Head.Load()
	}

	// Create and send a batch of txs.
	_, signedTxs, err := buildPreconfBlock(context.Background(), d.Client, d.PreconfServerURL(), l1Head, l2BlockID, nil)
	if err != nil {
		return nil, nil, err
	}
	defer d.Logf("%s: build soft block end, transaction length: %d", d.ClientType(), signedTxs.Len())

	return l1Head, signedTxs, nil
}

func buildPreconfBlock(
	ctx context.Context,
	rpccli *rpc.Client,
	preconfURL string,
	anchoredL1Block *types.Header,
	l2BlockID uint64,
	txs types.Transactions,
) (*types.Header, types.Transactions, error) {
	l2cli := rpccli.L2

	// Create and send a batch of txs.
	/*if txs == nil {
		signedTxs, err := utils.CreateL2Txs(context.Background(), l2cli, true)
		if err != nil {
			return nil, nil, err
		}
		txs = signedTxs
	}*/
	txBytes, err := utils.EncodeAndCompressTxList(txs)
	if err != nil {
		return nil, nil, err
	}

	preconferPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIV_KEY")))
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

	reqBody := &preconfblocks.BuildPreconfBlockRequestBody{
		ExecutableData: &preconfblocks.ExecutableData{
			ParentHash:   parent.Hash(),
			FeeRecipient: params.ParamToAddress("L2_SUGGESTED_FEE_RECIPIENT"),
			Number:       l2BlockID,
			GasLimit:     uint64(preconfCfg.BlockMaxGasLimit()),
			Timestamp:    anchoredL1Block.Time,
			Transactions: txBytes,
		},
		AnchorBlockID:   anchoredL1Block.Number.Uint64(),
		AnchorStateRoot: anchoredL1Block.Root,
		SignalSlots:     [][32]byte{},
		BaseFeeConfig:   preconfCfg.BaseFeeConfig(),
	}

	payload, err := rlp.EncodeToBytes(reqBody)
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
