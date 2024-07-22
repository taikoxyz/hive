package main

import taiko2 "github.com/ethereum/hive/taiko"

func main() {
	taiko := taiko2.NewCommonSpec()
	taiko.RunSuite(l1InitSuite())
	taiko.RunSuite(l2InitSuite())
	taiko.RunSuite(driverSuite())
	taiko.RunSuite(proposerSuite())
	taiko.RunSuite(proverSuite())

	taiko.Release()
}
