package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func l1InitSuite() hivesim.Suite {
	l1geth := hivesim.Suite{
		Name:        "l1geth",
		Description: `l1-geth initialization`,
	}
	l1geth.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "deployL1Contract",
		Description: "Deploy taiko contract on l1 chain and get environment variables",
		Run:         deployL1Contract,
		AlwaysRun:   true,
	})

	return l1geth
}

// Deploy a contract on L1
func deployL1Contract(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createAndConnectNetwork(t, c.Container)

	if err := os.Setenv("L1_NODE_HTTP_ENDPOINT", fmt.Sprintf("http://%v:8545", c.IP)); err != nil {
		t.Fatal(err)
	}
	// deploy l1 contract.
	cmd := exec.Command("sh", "/taiko/deploy_l1_contract.sh")
	if err := runTAP(t, c.Type, cmd); err != nil {
		t.Fatal(err)
	}

	setL1Env(t, c)
}

func setL1Env(t *hivesim.T, c *hivesim.Client) {
	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	//envs := make(map[string]string)
	envs["L1_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L1_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L1_BEACON"] = fmt.Sprintf("http://%v:4000", c.IP)
	envs["L1_PROPOSER_PRIV_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	envs["L1_PROVER_PRIV_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	t.Logf("container: %s, envs: %v", c.Container, envs)

	if err := godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	if err := godotenv.Load(envFile); err != nil {
		t.Fatal("failed to load env file", err)
	}
}
