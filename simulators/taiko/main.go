package main

import "time"

func main() {
	taiko := NewCommonSpec()
	taiko.RunSuite(l1InitSuite(), l2InitSuite())
	taiko.RunSuite(clientSuite())

	time.Sleep(100 * time.Second)

	taiko.Release()
}
