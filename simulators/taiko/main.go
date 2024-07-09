package main

import (
	"github.com/ethereum/hive/hivesim"
)

func main() {

	suite := hivesim.Suite{
		Name:        "taiko",
		Description: ``,
	}

	var runners []hivesim.ClientTestSpec
	runners = append(runners, l1Suite()...)
	runners = append(runners, l2Suite()...)
	runners = append(runners, clientSuite()...)

	for _, runner := range runners {
		suite.Add(runner)
	}

	hivesim.MustRun(hivesim.New(), suite)
}
