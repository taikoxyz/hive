package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/taiko"
	"math/big"
	"os"
	"time"
)

func proverTestSuite() hivesim.Suite {
	suite := taiko.ProverSuite(network, envFile)
	suite.Add(hivesim.TestSpec{
		Name:        "prover",
		Description: "test prover work normally",
		Run:         testPriver,
		AlwaysRun:   true,
	})
	return suite
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
