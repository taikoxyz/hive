package params

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:embed l1contract_txs.txt
var txsContent []byte

var ContractTxs = make([]*types.Transaction, 0)

func initTxs() {
	// Load taiko contract txs.
	fscan := bufio.NewScanner(bytes.NewReader(txsContent))
	fscan.Split(bufio.ScanLines)
	for fscan.Scan() {
		var tx types.Transaction
		if err := json.Unmarshal([]byte(fscan.Text()), &tx); err != nil {
			panic(err)
		}
		ContractTxs = append(ContractTxs, &tx)
	}
}
