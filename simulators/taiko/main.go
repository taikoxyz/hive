package main

import "time"

func main() {
	taiko := NewCommonSpec()
	taiko.RunSuite(l1InitSuite())
	taiko.RunSuite(l2InitSuite())
	taiko.RunSuite(driverSuite())
	taiko.RunSuite(proposerSuite())
	taiko.RunSuite(proverSuite())

	time.Sleep(50 * time.Second)

	taiko.Release()
}
