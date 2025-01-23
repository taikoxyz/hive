package clients

import (
	"context"
	tkutils "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/utils"
	"github.com/urfave/cli/v2"
	"os"
)

type DriverClient struct {
	*HiveManagedClient
}

func (d *DriverClient) Start() error {
	if err := d.HiveManagedClient.Start(); err != nil {
		return err
	}

	d.Logf("driver client, L1_BEACON: %s", os.Getenv("L1_BEACON"))

	return nil
}

func (d *DriverClient) Shutdown() error {
	if err := d.HiveManagedClient.Shutdown(); err != nil {
		return err
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
