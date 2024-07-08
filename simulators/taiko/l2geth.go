package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
)

func setL2Env(t *hivesim.T, c *hivesim.Client) {
	setEnv(t, "L2_AUTH", fmt.Sprintf("http://%v:8551", c.IP))
	setEnv(t, "L2_WS", fmt.Sprintf("ws://%v:8546", c.IP))
	setEnv(t, "TAIKO_L2", "0x1670010000000000000000000000000000010001")
}

func testGeth(t *hivesim.T, c *hivesim.Client) {
	url := fmt.Sprintf("http://%v:8545", c.IP)

	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("http endpoint: %s", url)
	t.Logf("geth message, chainid: %d", chainID.Uint64())
}
