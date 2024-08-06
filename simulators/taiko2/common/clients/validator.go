package clients

import (
	"fmt"
	"github.com/marioevz/eth-clients/clients"
	"taiko2/common/utils"
)

type ValidatorClient struct {
	Client
	Logger      utils.Logging
	ClientIndex int

	BeaconClient *BeaconClient
}

func (vc *ValidatorClient) Logf(format string, values ...interface{}) {
	if l := vc.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (vc *ValidatorClient) Start() error {
	if !vc.Client.IsRunning() {
		if managedClient, ok := vc.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			return managedClient.Start()
		}
	}
	return nil
}

func (vc *ValidatorClient) Shutdown() error {
	if managedClient, ok := vc.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}

type ValidatorClients []*ValidatorClient

// Return subset of clients that are currently running
func (all ValidatorClients) Running() ValidatorClients {
	res := make(ValidatorClients, 0)
	for _, vc := range all {
		if vc.IsRunning() {
			res = append(res, vc)
		}
	}
	return res
}
