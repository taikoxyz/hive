package taiko

import (
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func DriverSuite(network, envfile string) hivesim.Suite {
	client := hivesim.Suite{
		Name:        "driver",
		Description: `driver connection test`,
	}
	params, _ := godotenv.Read(envfile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "driver",
		Name:        "driver",
		Description: "driver join network and init taiko contracts",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			CreateOrConnectNetwork(t, c.Container, network)
			if err := InitTaikoContract(envfile); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	return client
}

func ProposerSuite(network, envfile string) hivesim.Suite {
	client := hivesim.Suite{
		Name:        "proposer",
		Description: `proposer connection test`,
	}
	params, _ := godotenv.Read(envfile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "proposer",
		Name:        "proposer",
		Description: "proposer join network and init taiko contract",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			CreateOrConnectNetwork(t, c.Container, network)
			if err := InitTaikoContract(envfile); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	return client
}

func ProverSuite(network, envfile string) hivesim.Suite {
	client := hivesim.Suite{
		Name:        "prover",
		Description: `prover connection test`,
	}
	params, _ := godotenv.Read(envfile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "prover",
		Name:        "prover",
		Description: "prover join network and init taiko contract",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			CreateOrConnectNetwork(t, c.Container, network)
			if err := InitTaikoContract(envfile); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	return client
}
