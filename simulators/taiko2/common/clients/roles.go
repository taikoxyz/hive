package clients

import "github.com/ethereum/hive/hivesim"

type ClientDefinitionsByRole struct {
	// L1 chain clients
	Eth1      []*hivesim.ClientDefinition `json:"eth1"`
	Beacon    []*hivesim.ClientDefinition `json:"beacon"`
	Validator []*hivesim.ClientDefinition `json:"validator"`
	// L2 chain clients
	Eth2     []*hivesim.ClientDefinition `json:"taiko-geth"`
	Driver   []*hivesim.ClientDefinition `json:"driver"`
	Proposer []*hivesim.ClientDefinition `json:"proposer"`
	Prover   []*hivesim.ClientDefinition `json:"prover"`
	// other clients.
	Other []*hivesim.ClientDefinition `json:"Other"`
}

func ClientsByRole(
	available []*hivesim.ClientDefinition,
) *ClientDefinitionsByRole {
	var out ClientDefinitionsByRole
	for _, client := range available {
		switch true {
		case client.HasRole("eth1"):
			out.Eth1 = append(out.Eth1, client)
		case client.HasRole("beacon"):
			out.Beacon = append(out.Beacon, client)
		case client.HasRole("validator"):
			out.Validator = append(out.Validator, client)
		case client.HasRole("taiko-geth"):
			out.Eth2 = append(out.Eth2, client)
		case client.HasRole("driver"):
			out.Driver = append(out.Driver, client)
		case client.HasRole("proposer"):
			out.Proposer = append(out.Proposer, client)
		case client.HasRole("prover"):
			out.Prover = append(out.Prover, client)
		}
	}
	return &out
}

func (c *ClientDefinitionsByRole) ClientByNameAndRole(
	name, role string,
) *hivesim.ClientDefinition {
	switch role {
	case "beacon":
		return byName(c.Beacon, name)
	case "validator":
		return byName(c.Validator, name)
	case "eth1":
		return byName(c.Eth1, name)
	case "taiko-geth":
		return byName(c.Eth2, name)
	case "driver":
		return byName(c.Driver, name)
	case "proposer":
		return byName(c.Proposer, name)
	case "prover":
		return byName(c.Prover, name)
	}
	return nil
}

func byName(
	clients []*hivesim.ClientDefinition,
	name string,
) *hivesim.ClientDefinition {
	for _, client := range clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}

func (c *ClientDefinitionsByRole) Combinations() NodeDefinitions {
	var (
		nodes     NodeDefinitions
		eth1      = c.Eth1
		beacon    = c.Beacon
		validator = c.Validator
		eth2      = c.Eth2
		driver    = c.Driver
		proposer  = c.Proposer
		prover    = c.Prover
		flag      = true
	)

	for flag {
		flag = false
		node := NodeDefinition{}
		if len(eth1) > 0 {
			node.L1EthClient = eth1[0].Name
			eth1 = eth1[1:]
			flag = true
		}
		if len(beacon) > 0 {
			node.ConsensusClient = beacon[0].Name
			beacon = beacon[1:]
			flag = true
		}
		if len(validator) > 0 {
			node.ValidatorClient = validator[0].Name
			validator = validator[1:]
			flag = true
		}
		if len(eth2) > 0 {
			node.L2EthClient = eth2[0].Name
			eth2 = eth2[1:]
			flag = true
		}
		if len(driver) > 0 {
			node.DriverClient = driver[0].Name
			driver = driver[1:]
			flag = true
		}
		if len(proposer) > 0 {
			node.ProposerClient = proposer[0].Name
			proposer = proposer[1:]
			flag = true
		}
		if len(prover) > 0 {
			node.ProverClient = prover[0].Name
			prover = prover[1:]
			flag = true
		}
		if flag {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
