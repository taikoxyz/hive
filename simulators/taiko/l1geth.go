package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
	"math/big"
	"os"
	"os/exec"
)

func l1InitSuite() hivesim.Suite {
	l1geth := hivesim.Suite{
		Name:        "l1geth",
		Description: `l1-geth initialization`,
	}
	l1geth.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "deployContractByTxs",
		Description: "Deploy taiko contract on l1 chain and get environment variables",
		Run:         deployContractByTxs,
		AlwaysRun:   true,
	})

	return l1geth
}

func deployContractByTxs(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createAndConnectNetwork(t, c.Container)

	fp, err := os.Open("./l1contract_txs.txt")
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

	setL1Env(t, c)
}

// Deploy a contract on L1
func deployL1Contract(t *hivesim.T, c *hivesim.Client) {
	// Create and connect network.
	createAndConnectNetwork(t, c.Container)

	if err := os.Setenv("L1_NODE_HTTP_ENDPOINT", fmt.Sprintf("http://%v:8545", c.IP)); err != nil {
		t.Fatal(err)
	}
	// deploy l1 contract.
	cmd := exec.Command("sh", "/taiko/deploy_l1_contract.sh")
	if err := runTAP(t, c.Type, cmd); err != nil {
		t.Fatal(err)
	}

	setL1Env(t, c)

	cli := ethclient.NewClient(c.RPC())
	number, err := cli.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// get txs
	txs := make([]*types.Transaction, 0)
	for i := 0; i <= int(number); i++ {
		block, err := cli.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			t.Fatal("failed to get block", err)
		}
		txs = append(txs, block.Transactions()...)
	}

	for _, tx := range txs {
		dt, err := json.Marshal(tx)
		if err != nil {
			t.Error(err)
		}

		t.Log(",,,,,,,,,,,,,,,, tx content: ", string(dt))
	}
}

func setL1Env(t *hivesim.T, c *hivesim.Client) {
	envs, err := godotenv.Read(envFile)
	if err != nil {
		t.Fatal("failed to load env file", err)
	}

	envs["L1_HTTP"] = fmt.Sprintf("http://%v:8545", c.IP)
	envs["L1_WS"] = fmt.Sprintf("ws://%v:8546", c.IP)
	envs["L1_BEACON"] = fmt.Sprintf("http://%v:3500", c.IP)
	envs["TAIKO_L1"] = "0x9E545E3C0baAB3E08CdfD552C960A1050f373042"
	envs["TAIKO_TOKEN"] = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	envs["L1_PROPOSER_PRIV_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	t.Logf("container: %s, envs: %v", c.Container, envs)

	if err := godotenv.Write(envs, envFile); err != nil {
		t.Fatal("failed to write env file", err)
	}
}
