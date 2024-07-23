package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/taiko"
	"github.com/ethereum/hive/taiko/bindings/taikol1"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func proposerTests() hivesim.Suite {
	suite := taiko.ProposerSuite(network, envFile)
	suite.Add(hivesim.TestSpec{
		Name:        "proposer test",
		Description: "make sure proposer can work normally",
		Run:         testNewProposeEvent,
	})

	return suite
}

func testNewProposeEvent(t *hivesim.T) {
	if err := godotenv.Load(envFile); err != nil {
		t.Fatal(err)
	}
	l1Client, err := ethclient.Dial(os.Getenv("L1_WS"))
	if err != nil {
		t.Fatal(err)
	}
	token, err := taikol1.NewTaikoL1(common.HexToAddress(os.Getenv("TAIKO_L1")), l1Client)

	sink := make(chan *taikol1.TaikoL1BlockProposed, 3)
	sub, err := token.WatchBlockProposed(nil, sink, nil, nil)
	if err != nil {
		t.Fatal("failed to watch blockProposed event", err)
	}
	defer sub.Unsubscribe()

	// wait 100 seconds until proposer done one.
	select {
	case <-time.After(time.Second * 100):
	case event := <-sink:
		t.Logf("blockID: %d, meta: %v", event.BlockId.Uint64(), event.Meta)
	}
}
