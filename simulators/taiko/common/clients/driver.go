package clients

import (
	"fmt"
	"github.com/marioevz/eth-clients/clients"
	"taiko/common/utils"
)

type DriverClient struct {
	Client
	Logger      utils.Logging
	ClientIndex int
}

func (d *DriverClient) Logf(format string, values ...interface{}) {
	if l := d.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (d *DriverClient) Start() error {
	d.Logf("Starting driver client %d", d.ClientIndex)
	if !d.Client.IsRunning() {
		if managedClient, ok := d.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			return managedClient.Start()
		}
	}
	return nil
}

func (d *DriverClient) Shutdown() error {
	if managedClient, ok := d.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}
