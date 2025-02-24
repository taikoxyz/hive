package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"math/big"
	"taiko/bindings/pacaya"
	"taiko/common/clients"
	"taiko/common/utils"
	"time"
)

func preconferBlock(envs hivesim.Params, rpccli *rpc.Client, preconfURL string, preconfs int) (*types.Header, *types.Header, error) {
	// get txs from l2 node tx mempool.
	var (
		ctx          = context.Background()
		l1cli, l2cli = rpccli.L1, rpccli.L2
	)

	anchorL1Header, err := l1cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get anchor header: %v", err)
	}

	l2BlockID, err := l2cli.BlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}

	var latestL2Header *types.Header
	for idx := 1; idx <= preconfs; idx++ {

		signedTxs, err := utils.CreateL2Txs(context.Background(), l2cli, true)
		if err != nil {
			return nil, nil, err
		}

		latestL2Header, _, err = clients.BuildPreconfBlock(ctx, rpccli, preconfURL, anchorL1Header, l2BlockID+uint64(idx), signedTxs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build preconf block: %v", err)
		}
		time.Sleep(time.Second)
	}

	return latestL2Header, anchorL1Header, nil
}

func proposeBlock(ctx context.Context, envs hivesim.Params, rpccli *rpc.Client, anchorheader *types.Header) (*rawdb.L1Origin, []types.Transactions, error) {
	l2cli := rpccli.L2
	proposerClient := &proposer.Proposer{}
	if err := clients.NewTaikoClient(proposerClient, flags.ProposerFlags); err != nil {
		return nil, nil, err
	}

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
		txs []types.Transactions
	)
	for number := canonicalL1Origin.BlockID.Uint64() + 1; number <= l2Number; number++ {
		l2Block, err := l2cli.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, l2Block.Transactions()[1:])
	}

	builder := NewCalldataTransactionBuilder(
		envs,
		rpccli,
		proposerClient.ProposeBlockTxGasLimit,
		config.NewChainConfig(
			l2cli.ChainID,
			0,
			pacaya.PacayaForkNumber(l2cli),
		),
		proposerClient.RevertProtectionEnabled,
	)

	txCandidate, err := builder.BuildPacaya(ctx, txs, anchorheader)
	if err != nil {
		return nil, nil, err
	}

	return canonicalL1Origin, txs, proposerClient.SendTx(ctx, txCandidate)
}
