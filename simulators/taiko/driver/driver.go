package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/taiko"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func driverSuite() hivesim.Suite {
	suite := taiko.DriverSuite(network, envFile)

	suite.Add(&hivesim.TestSpec{
		Name:        "driver test",
		Description: "make sure driver can work normally",
		Run:         testDriver,
	})

	return suite
}

func testDriver(t *hivesim.T) {
	if err := godotenv.Load(envFile); err != nil {
		t.Fatal(err)
	}

	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()

	for {
		select {
		case <-time.After(time.Second * 100):
			t.Errorf("drvier work timeout but taiko-geth number not changed")
		case <-tick.C:
			l2Client, err := ethclient.Dial(os.Getenv("L2_HTTP"))
			if err != nil {
				t.Fatal(err)
			}
			l2Number, err := l2Client.BlockNumber(context.Background())
			if err != nil {
				t.Error("failed to get l2 latest number", err)
			}
			if l2Number > 0 {
				t.Logf("taiko-geth's latest number changed")
				return
			}
		}
	}
}
