package params

import (
	_ "embed"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
	"os"
)

//go:embed .env
var envContent []byte

//go:embed config.yml
var ConfigContent []byte

//go:embed genesis.json
var GenesisContent []byte

var (
	envParams   = hivesim.Params{}
	ClusterEnvs = map[int]hivesim.Params{}
)

func init() {
	initEnvs()
	initJWT()
	initTxs()
	initAccounts()
}

func initEnvs() {
	// Load env params.
	vals, err := godotenv.UnmarshalBytes(envContent)
	if err != nil {
		panic(err)
	}
	for k, v := range vals {
		envParams[k] = v
		// set envs
		_ = os.Setenv(k, v)
	}
}

func ParamByKey(key string) string {
	return envParams[key]
}

func ParamToAddress(key string) common.Address {
	return common.HexToAddress(envParams[key])
}

func EnvParams() hivesim.Params {
	return envParams.Copy()
}

func SetEnvParams(key, value string) {
	if value == "" {
		delete(envParams, key)
	} else {
		envParams[key] = value
	}
	// set envs
	_ = os.Setenv(key, value)
}
