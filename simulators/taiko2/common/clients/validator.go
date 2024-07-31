package clients

import (
	"fmt"
	"github.com/marioevz/eth-clients/clients"
	"github.com/marioevz/eth-clients/clients/validator"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"taiko2/common/utils"
)

type ValidatorClient struct {
	Client
	Logger      utils.Logging
	ClientIndex int

	Keys         map[common.ValidatorIndex]*validator.ValidatorKeys
	BeaconClient *BeaconClient
}

func (vc *ValidatorClient) Logf(format string, values ...interface{}) {
	if l := vc.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (vc *ValidatorClient) Start() error {
	if !vc.Client.IsRunning() {
		if len(vc.Keys) == 0 {
			vc.Logf("Skipping validator because it has 0 validator keys")
			return nil
		}
		if managedClient, ok := vc.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			return managedClient.Start()
		}
	}
	return nil
}

func (v *ValidatorClient) ContainsKey(pk [48]byte) bool {
	for _, k := range v.Keys {
		if k.ValidatorPubkey == pk {
			return true
		}
	}
	return false
}

func (vc *ValidatorClient) Shutdown() error {
	if managedClient, ok := vc.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}

func (v *ValidatorClient) ContainsValidatorIndex(
	index common.ValidatorIndex,
) bool {
	_, ok := v.Keys[index]
	return ok
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
