package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-resty/resty/v2"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/common/utils"
	"taiko/params"
)

type DriverClient struct {
	*HiveManagedClient

	*rpc.Client
	*State
}

func (d *DriverClient) Start() (err error) {
	if err := d.HiveManagedClient.Start(); err != nil {
		return err
	}

	d.Logf("driver client, L1_BEACON: %s", os.Getenv("L1_BEACON"))

	d.Client, err = rpc.NewClient(context.Background(), GetClientConfig())
	if err != nil {
		return err
	}

	d.State, err = NewState(d.Client)
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
	signedTxs, err := buildPreconfBlock(context.Background(), d.Client, d.PreconfServerURL(), l2BlockID, l1Head)
	if err != nil {
		return nil, nil, err
	}
	defer d.Logf("%s: build soft block end, transaction length: %d", d.ClientType(), signedTxs.Len())

	return l1Head, signedTxs, d.StateError()
}

func (d *DriverClient) RemovePreconfBlocks(newLastBlockID uint64) error {
	d.Logf("%s: remove soft block, target height: %d", d.ClientType(), newLastBlockID)
	defer d.Logf("%s: remove soft block end", d.ClientType())

	return removePreconfBlocks(d.PreconfServerURL(), newLastBlockID)
}

func buildPreconfBlock(
	ctx context.Context,
	rpcCli *rpc.Client,
	preconfServerURL string,
	l2BlockID uint64,
	l1Head *types.Header,
) (types.Transactions, error) {
	// Create and send a batch of txs.
	signedTxs, err := utils.CreateL2Txs(context.Background(), rpcCli.L2, true)
	if err != nil {
		return nil, err
	}
	txBytes, err := utils.EncodeAndCompressTxList(signedTxs)
	if err != nil {
		return nil, err
	}

	l2Block, err := rpcCli.L2.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	executableData := engine.BlockToExecutableData(l2Block, nil, nil)
	executableData.ExecutionPayload.Transactions = [][]byte{txBytes}

	var txBatch = &preconfblocks.BuildPreconfBlockRequestBody{
		// TODO
		//ExecutableData:  executableData.ExecutionPayload,
		AnchorBlockID:   l1Head.Number.Uint64(),
		AnchorStateRoot: l1Head.Root,
		SignalSlots:     [][32]byte{},
	}
	payload, err := rlp.EncodeToBytes(txBatch)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(crypto.Keccak256(payload), params.ChainAuths[0].PrivateKey)
	if err != nil {
		return nil, err
	}
	txBatch.Signature = common.Bytes2Hex(sig)

	// Try to propose a soft block with batch ID 0
	res, err := resty.New().
		R().
		SetBody(&preconfblocks.BuildPreconfBlockRequestBody{}).
		Post(preconfServerURL + "/preconfBlocks")
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return append(l2Block.Transactions(), signedTxs...), errors.New(res.String())
	}

	return signedTxs, nil
}

func removePreconfBlocks(softURL string, newLastBlockID uint64) error {
	// Remove soft blocks
	res, err := resty.New().
		R().
		SetBody(&preconfblocks.RemovePreconfBlocksRequestBody{
			NewLastBlockID: newLastBlockID,
		}).
		Delete(softURL + "/preconfBlocks")
	if err != nil {
		return err
	}
	if !res.IsSuccess() {
		return errors.New(res.String())
	}
	return nil
}
