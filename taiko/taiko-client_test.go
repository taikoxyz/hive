package taiko

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"
)

func TestDeployTaikoContracts(t *testing.T) {
	cli, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)

	chainID, err := cli.ChainID(context.Background())
	assert.NoError(t, err)
	t.Log(chainID.Uint64())

	fp, err := os.Open("./l1contract_txs.txt")
	assert.NoError(t, err)

	fscan := bufio.NewScanner(fp)
	fscan.Split(bufio.ScanLines)

	txs := make([]*types.Transaction, 0)
	for fscan.Scan() {
		var tx types.Transaction
		err = json.Unmarshal([]byte(fscan.Text()), &tx)
		assert.NoError(t, err)
		txs = append(txs, &tx)
	}
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Nonce() < txs[j].Nonce()
	})

	for _, tx := range txs {
		if err := cli.SendTransaction(context.Background(), tx); err != nil {
			t.Error(err)
		}
		t.Logf("send tx, hash: %s, nonce: %d", tx.Hash().String(), tx.Nonce())
	}

	for _, signedTx := range txs {
		if signedTx.To() == nil {
			_, err := bind.WaitDeployed(context.Background(), cli, signedTx)
			assert.NoError(t, err)
			t.Logf("hash: %s, nonce: %d", signedTx.Hash().Hex(), signedTx.Nonce())
		} else {
			receipt, err := bind.WaitMined(context.Background(), cli, signedTx)
			assert.NoError(t, err)
			t.Logf(" hash: %s, nonce: %d, status: %d", signedTx.Hash().Hex(), signedTx.Nonce(), receipt.Status)
		}
	}
}

func TestSetupTest(t *testing.T) {
	assert.NoError(t, godotenv.Load(".env"))

	t.Log(os.Getenv("L1_PROPOSER_PRIV_KEY"))
	t.Log(os.Getenv("L1_HTTP"))

	err := initTaikoContract(".env")
	assert.NoError(t, err)
}

func TestL2Reorg(t *testing.T) {
	ctx := context.Background()

	err := RevertL1Geth(ctx, 10)
	assert.NoError(t, err)
	//var number = uint64(3)
	//err := Revert(ctx, os.Getenv("L2_"), "0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345", number)
	//if err != nil {
	//	t.Fatalf("failed to revert taiko-geth to %d number, err: %s", number, err.Error())
	//}

}
