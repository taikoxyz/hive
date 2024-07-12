package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
	"os"
)

func clientSuite() hivesim.Suite {
	client := hivesim.Suite{
		Name:        "taikoClient",
		Description: `taikoClient connection test`,
	}
	params, _ := godotenv.Read(envFile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "driver",
		Name:        "driver",
		Description: "driver test connection",
		Parameters:  params,
		Run:         testConnection,
		AlwaysRun:   true,
	})
	client.Add(hivesim.ClientTestSpec{
		Role:        "proposer",
		Name:        "proposer",
		Description: "proposer test connection",
		Parameters:  params,
		Run:         testConnection,
		AlwaysRun:   true,
	})
	return client
}

func testConnection(t *hivesim.T, c *hivesim.Client) {
	createAndConnectNetwork(t, c.Container)
	l1Url := os.Getenv("L1_WS")
	t.Logf("l1 l1Url: %s", l1Url)
	l2Url := os.Getenv("L2_WS")
	t.Logf("l2 l1Url: %s", l2Url)

	l1Cli, err := ethclient.Dial(l1Url)
	if err != nil {
		t.Fatal(err)
	}
	chainID, err := l1Cli.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testConnection test, l1 chainID: %d", chainID.Uint64())
}
