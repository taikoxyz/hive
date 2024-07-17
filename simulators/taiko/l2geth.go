package main

import (
	"fmt"

	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func l2InitSuite() hivesim.Suite {
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
		Run:       setL2Env,
		AlwaysRun: true,
	})

	return l2Geth
}

func setL2Env(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createOrConnectNetwork(t, c.Container)

	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L2_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L2_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L2_AUTH"] = fmt.Sprintf("http://%v:8551", c.IP)

	t.Logf("container: %s, envs: %v", c.Container, envs)

	if err := godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	if err = godotenv.Load(envFile); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}
}
