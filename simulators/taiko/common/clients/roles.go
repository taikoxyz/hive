package clients

import "github.com/ethereum/hive/hivesim"

type Role string

var (
	Anvil       Role = "anvil"
	Eth1        Role = "eth1"
	Beacon      Role = "beacon"
	Validator   Role = "validator"
	Eth2        Role = "taiko-geth"
	Reth        Role = "taiko-reth"
	Driver      Role = "driver"
	Proposer    Role = "proposer"
	Prover      Role = "prover"
	Redis       Role = "redis"
	Postgres    Role = "postgres"
	BlobApi     Role = "blobscan-api"
	BlobIndexer Role = "blobscan-indexer"

	Roles = map[Role]bool{
		Anvil:       true,
		Eth1:        true,
		Beacon:      true,
		Validator:   true,
		Eth2:        true,
		Reth:        true,
		Driver:      true,
		Proposer:    true,
		Prover:      true,
		Redis:       true,
		Postgres:    true,
		BlobApi:     true,
		BlobIndexer: true,
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
		for _, r := range client.Meta.Roles {
			if _, ok := Roles[Role(r)]; !ok {
				continue
			}
			role := Role(r)

			if group[role] != nil || (role == Eth2 && group[Reth] != nil) || (role == Reth && group[Eth2] != nil) {
				clientGroups = append(clientGroups, group)
				group = map[Role]*hivesim.ClientDefinition{role: client}
			} else {
				group[role] = client
			}
		}
	}
	if len(group) > 0 {
		clientGroups = append(clientGroups, group)
	}

	return clientGroups
}
