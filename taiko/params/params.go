package params

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

//go:embed .env
var envContent []byte

//go:embed l1contract_txs.txt
var txsContent []byte

//go:embed config.yml
var ConfigContent []byte

//go:embed genesis.json
var GenesisContent []byte

var (
	EnvParams   = hivesim.Params{}
	ContractTxs = make([]*types.Transaction, 0)
)

func init() {
	// Load env params.
	envParams, err := godotenv.UnmarshalBytes(envContent)
	if err != nil {
		panic(err)
	}
	for k, v := range envParams {
		EnvParams[k] = v
	}

	// Load taiko contract txs.
	fscan := bufio.NewScanner(bytes.NewReader(txsContent))
	fscan.Split(bufio.ScanLines)
	for fscan.Scan() {
		var tx types.Transaction
		if err = json.Unmarshal([]byte(fscan.Text()), &tx); err != nil {
			panic(err)
		}
		ContractTxs = append(ContractTxs, &tx)
	}
}
