package preconf

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"github.com/urfave/cli/v2"
	"math/big"
	_ "taiko/params"
)

func NewProposer() (*proposer.Proposer, error) {
	propose := &proposer.Proposer{}

	app := cli.NewApp()
	app.Flags = flags.ProposerFlags
	app.Commands = []*cli.Command{
		{
			Name:        "proposer",
			Flags:       flags.ProposerFlags,
			Usage:       "Starts the proposer software",
			Description: "Taiko proposer software",
			Action: func(c *cli.Context) error {
				return propose.InitFromCli(context.Background(), c)
			},
		},
	}

	if err := app.Run([]string{"taiko-client", "proposer"}); err != nil {
		return nil, err
	}

	return propose, nil
}

func ProposeTxLists(
	ctx context.Context,
	propose *proposer.Proposer,
	l2cli *ethclient.Client,
) error {
	canonicalL1OriginCh, err := l2cli.HeadL1Origin(ctx)
	if err != nil {
		return err
	}

	l2Number, err := l2cli.BlockNumber(ctx)
	if err != nil {
		return err
	}

	// Collect all the soft transactions.
	var txs []types.Transactions
	for number := canonicalL1OriginCh.BlockID.Uint64(); /*+ 1*/ number <= l2Number; number++ {
		l2Block, err := l2cli.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			return err
		}
		txs = append(txs, l2Block.Transactions())
	}

	return propose.ProposeTxLists(ctx, txs)
}
