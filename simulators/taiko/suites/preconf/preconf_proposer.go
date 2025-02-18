package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"math/big"
	"os"
	"taiko/bindings/pacaya"
	"taiko/common/clients"
)

func preconferProposer(rpccli *rpc.Client, preconfURL string) error {
	// get txs from l2 node tx mempool.
	var (
		ctx          = context.Background()
		l1cli, l2cli = rpccli.L1, rpccli.L2
		txsLst       = make([]types.Transactions, 0)
	)

	anchorL1Header, err := l1cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Println("preconferProposer anchor_id", anchorL1Header.Number.Uint64())

	for times := 5; times > 0; times-- {
		l2BlockID, err := l2cli.BlockNumber(ctx)
		if err != nil {
			return err
		}

		_, txs, err := clients.BuildPreconfBlock(ctx, rpccli, preconfURL, anchorL1Header, l2BlockID+1, nil)
		if err != nil {
			return err
		}
		txsLst = append(txsLst, txs)
	}

	// propose txs.
	_, _, err = proposeTxLists(ctx, rpccli, anchorL1Header)
	if err != nil {
		return err
	}

	return nil
}

func proposeTxLists(ctx context.Context, rpccli *rpc.Client, anchorheader *types.Header) (*rawdb.L1Origin, []types.Transactions, error) {
	l2cli := rpccli.L2
	canonicalL1Origin, err := l2cli.HeadL1Origin(ctx)
	if err != nil {
		return nil, nil, err
	}

	l2Number, err := l2cli.BlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}

	// no soft blocks need to be proposed
	if canonicalL1Origin.BlockID.Uint64() == l2Number {
		return nil, nil, nil
	}

	// Collect all the soft transactions.
	var (
		total int
		txs   []types.Transactions
	)
	for number := canonicalL1Origin.BlockID.Uint64() + 1; number <= l2Number; number++ {
		l2Block, err := l2cli.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, l2Block.Transactions()[1:])
		total += l2Block.Transactions().Len()
	}

	_ = os.Setenv("L2_HTTP", "http://localhost:6045")
	proposerClient := &proposer.Proposer{}
	if err := clients.NewTaikoClient(proposerClient, flags.ProposerFlags); err != nil {
		return nil, nil, err
	}

	builder := NewCalldataTransactionBuilder(
		rpccli,
		proposerClient.Config.ProposeBlockTxGasLimit,
		config.NewChainConfig(
			l2cli.ChainID,
			0,
			pacaya.PacayaForkNumber(l2cli),
		),
		proposerClient.Config.RevertProtectionEnabled,
	)

	txCandidate, err := builder.BuildPacaya(ctx, txs, anchorheader)
	if err != nil {
		return nil, nil, err
	}

	return canonicalL1Origin, txs, proposerClient.SendTx(ctx, txCandidate)
}
