package main

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func l2Suite() []hivesim.ClientTestSpec {
	return []hivesim.ClientTestSpec{
		{
			Role:        "taiko-geth",
			Name:        "setL2Env",
			Description: "Set environment variables for taiko-geth",
			Run:         setL2Env,
			AlwaysRun:   true,
		},
	}
}

func setL2Env(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createAndConnectNetwork(t, c.Container)

	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L2_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L2_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L2_AUTH"] = fmt.Sprintf("http://%v:8551", c.IP)
	envs["TAIKO_L2"] = "0x1670010000000000000000000000000000010001"
	envs["L2_SUGGESTED_FEE_RECIPIENT"] = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	t.Logf("////////////////////, %v", envs)

	if err = godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}
}
