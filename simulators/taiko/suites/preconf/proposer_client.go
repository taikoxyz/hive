package preconf

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	cmdutils "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/utils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"github.com/urfave/cli/v2"
	"math/big"
	_ "taiko/params"
)

func (m *MockProposer) InitFromCli(ctx context.Context, c *cli.Context) error {
	return m.Proposer.InitFromCli(ctx, c)
}
func (m *MockProposer) Name() string {
	return "proposer"
}
func (m *MockProposer) Start() error {
	return nil
}
func (m *MockProposer) Close(context.Context) {}

type MockProposer struct {
	*proposer.Proposer
}

func NewProposer() (*MockProposer, error) {
	propose := &MockProposer{Proposer: &proposer.Proposer{}}

	app := cli.NewApp()
	app.Flags = flags.ProposerFlags
	app.Commands = []*cli.Command{
		{
			Name:        "proposer",
			Flags:       flags.ProposerFlags,
			Usage:       "Starts the proposer software",
			Description: "Taiko proposer software",
			Action:      cmdutils.SubcommandAction(propose),
		},
	}

	if err := app.Run([]string{"taiko-client", "proposer"}); err != nil {
		return nil, err
	}

	return propose, nil
}

func (m *MockProposer) ProposeTxLists(ctx context.Context, l2cli *ethclient.Client) error {
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
	for number := canonicalL1OriginCh.L1BlockHeight.Uint64() + 1; number <= l2Number; number++ {
		l2Block, err := l2cli.BlockByNumber(ctx, big.NewInt(int64(number)))
		if err != nil {
			return err
		}
		txs = append(txs, l2Block.Transactions())
	}

	return m.Proposer.ProposeTxLists(ctx, txs)
}
