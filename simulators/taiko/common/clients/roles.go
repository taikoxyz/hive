package clients

import "github.com/ethereum/hive/hivesim"

type Role string

var (
	Anvil     Role = "anvil"
	Eth1      Role = "eth1"
	Beacon    Role = "beacon"
	Validator Role = "validator"
	Eth2      Role = "taiko-geth"
	Driver    Role = "driver"
	Proposer  Role = "proposer"
	Prover    Role = "prover"
	Roles          = map[Role]bool{
		Anvil:     true,
		Eth1:      true,
		Beacon:    true,
		Validator: true,
		Eth2:      true,
		Driver:    true,
		Proposer:  true,
		Prover:    true,
	}
)

type ClientsByRole map[Role]*hivesim.ClientDefinition

func GetClientsByRole(
	available []*hivesim.ClientDefinition,
) ClientsByRole {
	var roleNodes = ClientsByRole{}
	for _, client := range available {
		for _, role := range client.Meta.Roles {
			if _, ok := Roles[Role(role)]; ok {
				roleNodes[Role(role)] = client
			}
		}
	}
	return roleNodes
}
