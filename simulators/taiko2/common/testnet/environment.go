package testnet

import (
	"taiko2/common/clients"
)

type Environment struct {
	Clients        *clients.ClientDefinitionsByRole
	LogEngineCalls bool
}
