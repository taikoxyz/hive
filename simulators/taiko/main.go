package main

import (
	"github.com/ethereum/hive/hivesim"
)

var (
	l1geth = hivesim.Suite{
		Name:        "geth",
		Description: `l1-geth initialization`,
	}
	l2Geth = hivesim.Suite{
		Name:        "taiko-geth",
		Description: `taiko-geth initialization`,
	}
	taikoClient = hivesim.Suite{
		Name:        "taikoClient",
		Description: `taikoClient connection test`,
	}
)

func main() {
	// l1 geth initialization and test cases.
	l1geth.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "contract",
		Description: "Deploy taiko contract on l1 chain",
		Run:         deployL1Contract,
		AlwaysRun:   false,
	})
	l1geth.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "set env",
		Description: "",
		Run:         setL1Env,
		AlwaysRun:   false,
	})

	// taiko geth initialization and test cases.
	l2Geth.Add(hivesim.ClientTestSpec{
		Role:        "taiko-geth",
		Name:        "taiko-geth",
		Description: "Set environment variables for taiko-geth",
		Run:         setL2Env,
		AlwaysRun:   false,
	})
	l2Geth.Add(hivesim.ClientTestSpec{
		Role: "taiko-geth", Name: "taiko-geth",
		Description: "test taiko-geth connection",
		Run:         testGeth,
		AlwaysRun:   false,
	})

	// taikoClient connection and test cases.
	params := taikoClientEnv()
	taikoClient.Add(hivesim.ClientTestSpec{
		Role: "taikoClient", Name: "taikoClient",
		Description: "test taikoClient connection",
		Parameters:  params, AlwaysRun: false,
		Run: nil,
	})
	taikoClient.Add(hivesim.ClientTestSpec{
		Role: "proposer", Name: "proposer",
		Description: "test proposer connection",
		Parameters:  params, AlwaysRun: false,
		Run: nil,
	})
	taikoClient.Add(hivesim.ClientTestSpec{
		Role: "prover", Name: "prover",
		Description: "test prover connection",
		Parameters:  params, AlwaysRun: false,
		Run: nil,
	})

	// Run the simulations.
	suites := []hivesim.Suite{
		//l1geth,
		l2Geth,
		//taikoClient,
	}
	hivesim.MustRun(hivesim.New(), suites...)
}
