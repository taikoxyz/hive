package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func l1Suite() []hivesim.ClientTestSpec {
	return []hivesim.ClientTestSpec{
		{
			Role:        "geth",
			Name:        "contract",
			Description: "Deploy taiko contract on l1 chain",
			Run:         deployL1Contract,
			AlwaysRun:   true,
		},
		{
			Role:        "geth",
			Name:        "setL1Env",
			Description: "",
			Run:         setL1Env,
			AlwaysRun:   true,
		},
	}
}

func setL1Env(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createAndConnectNetwork(t, c.Container)

	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L1_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L1_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L1_BEACON"] = fmt.Sprintf("http://%v:4000", c.IP)
	envs["L1_PROPOSER_PRIV_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	envs["L1_PROVER_PRIV_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	t.Logf("===================, %v", envs)

	if err = godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	t.Logf("TAIKO_L1=%s", os.Getenv("TAIKO_L1"))
}

// Deploy a contract on L1
func deployL1Contract(t *hivesim.T, c *hivesim.Client) {
	url := fmt.Sprintf("http://%v:8545", c.IP)
	if err := os.Setenv("L1_NODE_HTTP_ENDPOINT", url); err != nil {
		t.Fatal(err)
	}
	// deploy l1 contract.
	cmd := exec.Command("sh", "/taiko/deploy_l1_contract.sh")
	if err := runTAP(t, c.Type, cmd); err != nil {
		t.Fatal(err)
	}
}
