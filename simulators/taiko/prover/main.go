package main

import "github.com/ethereum/hive/taiko"

func main() {
	runner := taiko.NewCommonSpec()
	runner.RunSuite(taiko.L1InitSuite(network, envFile, txsFile))
	runner.RunSuite(taiko.L2InitSuite(network, envFile))
	runner.RunSuite(taiko.DriverSuite(network, envFile))
	runner.RunSuite(taiko.ProposerSuite(network, envFile))
	runner.RunSuite(taiko.ProverSuite(network, envFile))

	runner.Release()
}
