package clients

import (
	"fmt"
	"github.com/marioevz/eth-clients/clients"
	"taiko2/common/utils"
)

type ProverClient struct {
	Client
	Logger      utils.Logging
	ClientIndex int
}

func (t *ProverClient) Logf(format string, values ...interface{}) {
	if l := t.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (t *ProverClient) Start() error {
	t.Logf("Starting prover client %d", t.ClientIndex)
	if !t.Client.IsRunning() {
		if managedClient, ok := t.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			return managedClient.Start()
		}
	}

	return nil
}

func (t *ProverClient) Shutdown() error {
	if managedClient, ok := t.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}
