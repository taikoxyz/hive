package testnet

import (
	"taiko/common/clients"
)

type Environment struct {
	Clients        *clients.ClientDefinitionsByRole
	LogEngineCalls bool
}
