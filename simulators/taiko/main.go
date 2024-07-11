package main

func main() {
	taiko := NewCommonSpec()
	if err := taiko.AddSuite(l1InitSuite(), l2InitSuite() /*, clientSuite()*/); err != nil {
		panic(err)
	}
	taiko.StartTest()
}
