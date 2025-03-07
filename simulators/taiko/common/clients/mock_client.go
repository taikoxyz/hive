package clients

import (
	"context"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	txmgrMetrics "github.com/ethereum-optimism/optimism/op-service/txmgr/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/hivesim"
	tkutils "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/utils"
	pkgFlags "github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	"github.com/urfave/cli/v2"
)

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

type MockClient struct {
	*proposer.Config
	Envs hivesim.Params
	txmgr.TxManager
}

func (t *MockClient) InitFromCli(ctx context.Context, c *cli.Context) (err error) {
	cfg, err := proposer.NewConfigFromCliContext(c)
	if err != nil {
		return err
	}
	t.Config = cfg

	cfg.L1ProposerPrivKey, err = crypto.ToECDSA(common.FromHex(t.Envs["L1_PROPOSER_PRIV_KEY"]))
	if err != nil {
		return nil
	}

	txMgrCfg := pkgFlags.InitTxmgrConfigsFromCli(
		t.Envs["L1_WS"],
		cfg.L1ProposerPrivKey,
		c,
	)

	t.TxManager, err = txmgr.NewSimpleTxManager(
		"txMgr",
		log.Root(),
		new(txmgrMetrics.NoopTxMetrics),
		*txMgrCfg,
	)

	return err
}

func (t *MockClient) Name() string {
	return "mock_client"
}

func (t *MockClient) Start() error { return nil }

func (t *MockClient) Close(ctx context.Context) {}
