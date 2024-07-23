package taiko

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func L2InitSuite(network string, envfile string) hivesim.Suite {
	l2Geth := hivesim.Suite{
		Name:        "l2geth",
		Description: `l2geth initialization`,
	}
	l2Geth.Add(hivesim.ClientTestSpec{
		Role:        "taiko-geth",
		Name:        "getL2Env",
		Description: "Get environment variables from taiko-geth",
		Parameters: map[string]string{
			"HIVE_CHECK_LIVE_PORT": "8545",
		},
		Run: func(t *hivesim.T, c *hivesim.Client) {
			// Create and connect network.
			CreateOrConnectNetwork(t, c.Container, network)
			setL2Env(t, c, envfile)
		},
		AlwaysRun: true,
	})

	return l2Geth
}

func setL2Env(t *hivesim.T, c *hivesim.Client, envfile string) {
	envs, err := godotenv.Read(envfile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L2_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L2_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L2_AUTH"] = fmt.Sprintf("http://%v:8551", c.IP)

	t.Logf("container: %s, envs: %v", c.Container, envs)

	if err := godotenv.Write(envs, envfile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	if err = godotenv.Load(envfile); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}
}
