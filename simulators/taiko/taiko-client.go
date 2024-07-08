package main

import "os"

func taikoClientEnv() map[string]string {
	return map[string]string{
		"L1_WS":                      os.Getenv("L1_WS"),
		"L1_HTTP":                    os.Getenv("L1_HTTP"),
		"L1_BEACON":                  os.Getenv("L1_BEACON"),
		"L2_HTTP":                    os.Getenv("L2_HTTP"),
		"L2_WS":                      os.Getenv("L2_WS"),
		"L2_AUTH":                    os.Getenv("L2_AUTH"),
		"TAIKO_L1":                   os.Getenv("TAIKO_L1"),
		"TAIKO_L2":                   os.Getenv("TAIKO_L2"),
		"TAIKO_TOKEN":                os.Getenv("TAIKO_TOKEN"),
		"L1_PROPOSER_PRIV_KEY":       os.Getenv("L1_PROPOSER_PRIV_KEY"),
		"L2_SUGGESTED_FEE_RECIPIENT": os.Getenv("L2_SUGGESTED_FEE_RECIPIENT"),
		"PROVER_SET":                 os.Getenv("PROVER_SET"),
		"L1_PROVER_PRIV_KEY":         os.Getenv("L1_PROVER_PRIV_KEY"),
	}
}
