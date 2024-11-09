package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/hive/hivesim"
	"github.com/go-resty/resty/v2"
	softblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/soft_blocks"
	"taiko/common/utils"
	"taiko/params"
	"time"
)

type DriverClient struct {
	SoftBlockServerPort uint64
	*HiveManagedClient

	softServerURL string
}

func (d *DriverClient) SoftServerURL() string {
	if d.softServerURL == "" {
		d.softServerURL = fmt.Sprintf("http://%s:%v", d.NetworkIP(), d.SoftBlockServerPort)
	}
	return d.softServerURL
}

func (d *DriverClient) BuildSoftBlock(
	t *hivesim.T,
	l1cli, l2cli *ethclient.Client,
	l2BlockID uint64,
	batchID uint64,
	endOfBlock bool,
	endOfPreconf bool,
) (int, error) {
	// Create and send a batch of txs.
	signedTxs := utils.CreateL2Txs(context.Background(), t, l2cli, true)
	b, err := utils.EncodeAndCompressTxList(signedTxs)
	if err != nil {
		return 0, err
	}

	var marker softblocks.TransactionBatchMarker
	if endOfBlock {
		marker = softblocks.BatchMarkerEOB
	} else if endOfPreconf {
		marker = softblocks.BatchMarkerEOP
	} else {
		marker = softblocks.BatchMarkerEmpty
	}

	l1Head, err := l1cli.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	var txBatch = &softblocks.TransactionBatch{
		BlockID:          l2BlockID,
		ID:               batchID,
		TransactionsList: b,
		BatchMarker:      marker,
		Signature:        "",
		BlockParams: &softblocks.SoftBlockParams{
			AnchorBlockID:   l1Head.Number.Uint64(),
			AnchorStateRoot: l1Head.Root,
			Timestamp:       l1Head.Time + 12,
			Coinbase:        params.L2Auths[0].From,
		},
	}
	payload, err := rlp.EncodeToBytes(txBatch)
	if err != nil {
		return 0, err
	}

	sig, err := crypto.Sign(crypto.Keccak256(payload), params.PrivateKeys[0])
	if err != nil {
		return 0, err
	}
	txBatch.Signature = common.Bytes2Hex(sig)

	// Try to propose a soft block with batch ID 0
	res, err := resty.New().
		R().
		SetBody(&softblocks.BuildSoftBlockRequestBody{
			TransactionBatch: txBatch,
		}).
		Post(d.SoftServerURL() + "/softBlocks")
	if err != nil {
		return 0, err
	}
	if !res.IsSuccess() {
		return 0, errors.New(res.String())
	}
	time.Sleep(time.Millisecond * 100)

	t.Log("Build soft block response", "body", res.String())

	return signedTxs.Len(), nil
}

func (d *DriverClient) RemoveSoftBlocks(t *hivesim.T, newLastBlockID uint64) error {
	// Remove soft blocks
	res, err := resty.New().
		R().
		SetBody(&softblocks.RemoveSoftBlocksRequestBody{
			NewLastBlockID: newLastBlockID,
		}).
		Delete(d.SoftServerURL() + "/softBlocks")
	if err != nil {
		return err
	}
	if !res.IsSuccess() {
		return errors.New(res.String())
	}
	time.Sleep(time.Millisecond * 100)

	t.Log("Remove soft blocks response", "body", res.String())

	return nil
}
