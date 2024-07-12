package main

import "time"

func main() {
	taiko := NewCommonSpec()
	taiko.RunSuite(l1InitSuite())
	taiko.RunSuite(l2InitSuite())
	taiko.RunSuite(clientSuite())

	time.Sleep(50 * time.Second)

	taiko.Release()
}
