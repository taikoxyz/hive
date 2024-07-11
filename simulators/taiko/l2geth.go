package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func l2InitSuite() hivesim.Suite {
	l2Geth := hivesim.Suite{
		Name:        "taiko-geth",
		Description: `taiko-geth initialization`,
	}
	l2Geth.Add(hivesim.ClientTestSpec{
		Role:        "taiko-geth",
		Name:        "getL2Env",
		Description: "Get environment variables from taiko-geth",
		Run:         setL2Env,
		AlwaysRun:   true,
	})

	return l2Geth
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

	if err := godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	if err := godotenv.Load(envFile); err != nil {
		t.Fatal("failed to load env file", err)
	}
}

func testGeth(t *hivesim.T) {
	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal(err)
	}

	url := envs["L2_HTTP"]
	t.Log("L2_HTTP endpoint: ", url)

	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Error(err)
	} else {
		t.Log("L2_HTTP chainID: ", chainID.Uint64())
	}

	t.Log("sleep 100 seconds ...")
}
