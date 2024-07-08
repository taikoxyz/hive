package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func setL1Env(t *hivesim.T, c *hivesim.Client) {
	// Show contract addresses
	if err := godotenv.Load("/taiko/.env"); err != nil {
		t.Fatal(err)
	}
	setEnv(t, "L1_WS", fmt.Sprintf("ws://%v:8546", c.IP))
	setEnv(t, "L1_BEACON", fmt.Sprintf("http://%v:4000", c.IP))
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
