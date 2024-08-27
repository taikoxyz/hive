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

	Roles = map[Role]bool{
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

type ClientGroups []map[Role]*hivesim.ClientDefinition

func GetClientsByRole(
	available []*hivesim.ClientDefinition,
) ClientGroups {
	var (
		clientGroups ClientGroups
		group        = make(map[Role]*hivesim.ClientDefinition)
	)

	for _, client := range available {
		for _, role := range client.Meta.Roles {
			if _, ok := Roles[Role(role)]; !ok {
				continue
			}
			if _, ok := group[Role(role)]; ok {
				clientGroups = append(clientGroups, group)
				group = map[Role]*hivesim.ClientDefinition{Role(role): client}
			} else {
				group[Role(role)] = client
			}
		}
	}
	if len(group) > 0 {
		clientGroups = append(clientGroups, group)
	}

	return clientGroups
}
