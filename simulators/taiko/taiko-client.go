package main

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/taiko"
	"github.com/joho/godotenv"
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
			taiko.CreateOrConnectNetwork(t, c.Container, network)
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
			taiko.CreateOrConnectNetwork(t, c.Container, network)
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
			taiko.CreateOrConnectNetwork(t, c.Container, network)
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

// Prover test case to check if prover work normally or not.
func testPriver(t *hivesim.T) {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	for {
		select {
		case <-time.After(time.Second * 200):
			t.Errorf("prover work timeout but finalized number not changed")
			return
		case <-tick.C:
			l2Client, err := ethclient.Dial(os.Getenv("L2_HTTP"))
			if err != nil {
				continue
			}
			header, err := l2Client.HeaderByNumber(context.Background(), big.NewInt(-3))
			if err != nil {
				continue
			} else if header.Number.Uint64() > 0 {
				t.Logf("finalized number changed, prover work normal, number: %d", header.Number.Uint64())
				return
			}
		}
	}
}
