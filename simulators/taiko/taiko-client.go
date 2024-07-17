package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
	"math/big"
	"os"
	"time"
)

func driverSuite() hivesim.Suite {
	client := hivesim.Suite{
		Name:        "driver",
		Description: `driver connection test`,
	}
	params, _ := godotenv.Read(envFile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "driver",
		Name:        "driver",
		Description: "driver join network and init taiko contracts",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			createOrConnectNetwork(t, c.Container)
			if err := initTaikoContract(); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	return client
}

func proposerSuite() hivesim.Suite {
	client := hivesim.Suite{
		Name:        "proposer",
		Description: `proposer connection test`,
	}
	params, _ := godotenv.Read(envFile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "proposer",
		Name:        "proposer",
		Description: "proposer join network and init taiko contract",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			createOrConnectNetwork(t, c.Container)
			if err := initTaikoContract(); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	return client
}

func proverSuite() hivesim.Suite {
	client := hivesim.Suite{
		Name:        "prover",
		Description: `prover connection test`,
	}
	params, _ := godotenv.Read(envFile)

	client.Add(hivesim.ClientTestSpec{
		Role:        "prover",
		Name:        "prover",
		Description: "prover join network and init taiko contract",
		Parameters:  params,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			createOrConnectNetwork(t, c.Container)
			if err := initTaikoContract(); err != nil {
				t.Fatalf("failed to init taiko contract: %v", err)
			}
		},
		AlwaysRun: true,
	})
	client.Add(hivesim.TestSpec{
		Name:        "prover",
		Description: "test prover work normally",
		Run:         testPriver,
		AlwaysRun:   true,
	})
	return client
}

func testPriver(t *hivesim.T) {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	for {
		select {
		case <-time.After(time.Second * 100):
			t.Errorf("prover work timeout but finalized number not changed")
			return
		case <-tick.C:
			l2Client, err := ethclient.Dial(os.Getenv("L2_HTTP"))
			if err != nil {
				continue
			}
			header, err := l2Client.HeaderByNumber(context.Background(), big.NewInt(-3))
			if err != nil {
				t.Errorf("failed to get header: %v", err)
			} else if header.Number.Uint64() > 0 {
				t.Logf("finalized number changed, prover work normal, number: %d", header.Number.Uint64())
				break
			}
		}
	}
}
