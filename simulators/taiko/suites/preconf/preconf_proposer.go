package preconf

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"taiko/bindings/pacaya"
	"taiko/common/clients"
	"time"
)

func preconferBlock(privateKey *ecdsa.PrivateKey, rpccli *rpc.Client, preconfURL string, preconfs int, l1Number *big.Int) (*types.Header, *types.Header, error) {
	// get txs from l2 node tx mempool.
	var (
		ctx   = context.Background()
		l1cli = rpccli.L1
	)

	anchorL1Header, err := l1cli.HeaderByNumber(ctx, l1Number)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get anchor header: %v", err)
	}

	var latestL2Header *types.Header
	for range preconfs {
		requestBody, err := clients.BuildPreconfRequestBody(ctx, rpccli, privateKey, anchorL1Header.Number, nil)
		if err != nil {
			return nil, nil, err
		}
		latestL2Header, err = clients.SendPreconfBlock(preconfURL, requestBody)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build preconf block: %v", err)
		}
		time.Sleep(time.Second)
	}

	return latestL2Header, anchorL1Header, nil
}

func proposeBlock(ctx context.Context, envs hivesim.Params, rpccli *rpc.Client, anchorheader *types.Header) (*rawdb.L1Origin, error) {
	l2cli := rpccli.L2

	mockClient := &clients.MockClient{Envs: envs}
	if err := clients.NewTaikoClient(mockClient, flags.ProposerFlags); err != nil {
		return nil, err
	}

	canonicalL1Origin, err := l2cli.HeadL1Origin(ctx)
	if err != nil {
		return nil, err
	}

	// Collect all the soft transactions.
	var allTxs []types.Transactions
	for number := canonicalL1Origin.BlockID.Uint64() + 1; true; number++ {
		l2Block, err := l2cli.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			break
		}
		allTxs = append(allTxs, l2Block.Transactions()[1:])
	}

	builder := NewCalldataTransactionBuilder(
		envs,
		rpccli,
		mockClient.ProposeBlockTxGasLimit,
		config.NewChainConfig(
			l2cli.ChainID,
			0,
			pacaya.PacayaForkNumber(l2cli),
		),
		mockClient.RevertProtectionEnabled,
	)

	txCandidate, err := builder.BuildPacaya(ctx, allTxs, anchorheader)
	if err != nil {
		return nil, err
	}

	_, err = mockClient.Send(ctx, *txCandidate)
	if err != nil {
		return nil, err
	}

	return canonicalL1Origin, nil
}
