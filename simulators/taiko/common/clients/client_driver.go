package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-resty/resty/v2"
	tkflags "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	tkutils "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/utils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/driver"
	softblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/soft_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/urfave/cli/v2"
	"os"
	"taiko/common/utils"
	"taiko/params"
)

type DriverClient struct {
	*HiveManagedClient

	*rpc.Client
	*driver.Driver

	*State
}

func (d *DriverClient) Start() error {
	if err := d.HiveManagedClient.Start(); err != nil {
		return err
	}

	d.Logf("driver client, L1_BEACON: %s", os.Getenv("L1_BEACON"))

	d.Driver = &driver.Driver{}
	err := NewTaikoClient(d.Driver, tkflags.DriverFlags)
	if err != nil {
		return err
	}

	data, err := json.Marshal(d.Config.ClientConfig)
	if err != nil {
		return err
	}
	d.Logf("driver's client config: %s", string(data))

	d.Client, err = rpc.NewClient(context.Background(), d.Config.ClientConfig)
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
	d.State.Close()

	return nil
}

func (d *DriverClient) SoftServerURL() string {
	return fmt.Sprintf("http://%s:%v", d.NetworkIP(), d.Config.SoftBlockServerPort)
}

func (d *DriverClient) BuildSoftBlock(
	l2BlockID uint64,
	batchID uint64,
	endOfBlock bool,
	endOfPreconf bool,
	l1Head *types.Header,
) (*types.Header, types.Transactions, error) {
	d.Logf("%s: build soft block", d.ClientType())

	if l1Head == nil {
		l1Head = d.L1Head.Load()
	}

	// Create and send a batch of txs.
	signedTxs, err := buildSoftBlock(d.Client, d.SoftServerURL(), l2BlockID, batchID, endOfBlock, endOfPreconf, l1Head)
	if err != nil {
		return nil, nil, err
	}
	defer d.Logf("%s: build soft block end, transaction length: %d", d.ClientType(), signedTxs.Len())

	return l1Head, signedTxs, d.StateError()
}

func (d *DriverClient) RemoveSoftBlocks(newLastBlockID uint64) error {
	d.Logf("%s: remove soft block, target height: %d", d.ClientType(), newLastBlockID)
	defer d.Logf("%s: remove soft block end", d.ClientType())

	return removeSoftBlocks(d.SoftServerURL(), newLastBlockID)
}

func buildSoftBlock(
	rpcCli *rpc.Client,
	softURL string,
	l2BlockID uint64,
	batchID uint64,
	endOfBlock bool,
	endOfPreconf bool,
	l1Head *types.Header,
) (types.Transactions, error) {
	// Create and send a batch of txs.
	signedTxs, err := utils.CreateL2Txs(context.Background(), rpcCli.L2, true)
	if err != nil {
		return nil, err
	}
	b, err := utils.EncodeAndCompressTxList(signedTxs)
	if err != nil {
		return nil, err
	}

	var marker softblocks.TransactionBatchMarker
	if endOfBlock {
		marker = softblocks.BatchMarkerEOB
	} else if endOfPreconf {
		marker = softblocks.BatchMarkerEOP
	} else {
		marker = softblocks.BatchMarkerEmpty
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
		return nil, err
	}

	sig, err := crypto.Sign(crypto.Keccak256(payload), params.PrivateKeys[0])
	if err != nil {
		return nil, err
	}
	txBatch.Signature = common.Bytes2Hex(sig)

	l2Block, err := rpcCli.L2.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Try to propose a soft block with batch ID 0
	res, err := resty.New().
		R().
		SetBody(&softblocks.BuildSoftBlockRequestBody{
			TransactionBatch: txBatch,
		}).
		Post(softURL + "/softBlocks")
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return append(l2Block.Transactions(), signedTxs...), errors.New(res.String())
	}

	return signedTxs, nil
}

func removeSoftBlocks(softURL string, newLastBlockID uint64) error {
	// Remove soft blocks
	res, err := resty.New().
		R().
		SetBody(&softblocks.RemoveSoftBlocksRequestBody{
			NewLastBlockID: newLastBlockID,
		}).
		Delete(softURL + "/softBlocks")
	if err != nil {
		return err
	}
	if !res.IsSuccess() {
		return errors.New(res.String())
	}
	return nil
}

func NewTaikoClient[T tkutils.SubcommandApplication](client T, flags []cli.Flag) error {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:  "client",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return client.InitFromCli(context.Background(), c)
			},
		},
	}
	return app.Run([]string{"taiko-client", "client"})
}
