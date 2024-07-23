package taiko

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
)

func L1InitSuite(network string, envfile string, txsfile string) hivesim.Suite {
	l1geth := hivesim.Suite{
		Name:        "l1geth",
		Description: `l1geth initialization`,
	}
	l1geth.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "deployContractByTxs",
		Description: "Deploy taiko contract on l1 chain and get environment variables",
		Parameters: map[string]string{
			"HIVE_CHECK_LIVE_PORT": "8545",
		},
		Run: func(t *hivesim.T, c *hivesim.Client) {
			// Create and connect network.
			CreateOrConnectNetwork(t, c.Container, network)
			deployContractByTxs(t, c, txsfile)
			setL1Env(t, c, envfile)
		},
		AlwaysRun: true,
	})

	return l1geth
}

func deployContractByTxs(t *hivesim.T, c *hivesim.Client, txsfile string) {
	fp, err := os.Open(txsfile)
	if err != nil {
		t.Fatalf("failed to open txs file: %v", err)
	}

	client := ethclient.NewClient(c.RPC())

	fscan := bufio.NewScanner(fp)
	fscan.Split(bufio.ScanLines)

	txs := make([]*types.Transaction, 0)
	for fscan.Scan() {
		var tx types.Transaction
		if err = json.Unmarshal([]byte(fscan.Text()), &tx); err != nil {
			t.Fatalf("failed to unmarshal tx: %v", err)
		}
		txs = append(txs, &tx)
	}

	for _, tx := range txs {
		if err = client.SendTransaction(context.Background(), tx); err != nil {
			t.Fatalf("failed to send tx: %v", err)
		}
	}

	// check results.
	for _, tx := range txs {
		if tx.To() == nil {
			_, err := bind.WaitDeployed(context.Background(), client, tx)
			if err != nil {
				t.Fatalf("failed to wait deployed: %v", err)
			}
		} else {
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				t.Fatalf("failed to wait mined, hash: %s, err: %v", tx.Hash().String(), err)
			}
			if receipt.Status != types.ReceiptStatusSuccessful {
				t.Fatalf("failed to call contract, hash: %s", tx.Hash().String())
			}
		}
	}
}

func setL1Env(t *hivesim.T, c *hivesim.Client, envfile string) {
	envs, err := godotenv.Read(envfile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L1_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L1_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L1_BEACON"] = fmt.Sprintf("http://%v:3500", c.IP)

	t.Logf("container: %s, envs: %v", c.Container, envs)

	if err := godotenv.Write(envs, envfile); err != nil {
		t.Fatal("failed to write env file", err)
	}

	if err = godotenv.Load(envfile); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}
}
