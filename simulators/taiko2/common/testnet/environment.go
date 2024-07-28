package testnet

import (
	"taiko2/common/clients"
	consensus_config "taiko2/common/config/consensus"
)

type Environment struct {
	Clients        *clients.ClientDefinitionsByRole
	Validators     consensus_config.ValidatorsSetupDetails
	LogEngineCalls bool
}
