package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/taiko/bindings/taikotoken"
)

const (
	envFile = "/taiko/.env"
	network = "hive_taiko_network"
)

var networkCreated = make(map[hivesim.SuiteID]bool)

// createNetwork ensures there is a separate network to be able to send the client traffic
// from two separate IP addrs.
func createOrConnectNetwork(t *hivesim.T, container string) {
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true
	}

	if err := t.Sim.ConnectContainer(t.SuiteID, network, container); err != nil {
		t.Fatal("can't connect container to network:", err)
	}
}

func initTaikoContract() error {
	l1client, err := ethclient.Dial(os.Getenv("L1_HTTP"))
	if err != nil {
		return err
	}

	l1ChainID, err := l1client.ChainID(context.Background())
	if err != nil {
		return err
	}

	taikoToken, err := taikotoken.NewTaikoToken(common.HexToAddress(os.Getenv("TAIKO_TOKEN")), l1client)
	if err != nil {
		return err
	}

	l1ProverPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROVER_PRIV_KEY")))
	if err != nil {
		return err
	}

	decimal, err := taikoToken.Decimals(nil)
	if err != nil {
		return err
	}
	allow := new(big.Int).Exp(big.NewInt(1_000_000_100), new(big.Int).SetUint64(uint64(decimal)), nil)
	fmt.Println(decimal, allow.String())

	allowance, err := taikoToken.Allowance(
		nil,
		crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey),
		common.HexToAddress(os.Getenv("TAIKO_L1")),
	)
	if err != nil {
		return err
	}

	if allowance.Cmp(common.Big0) == 0 {
		ownerPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_CONTRACT_OWNER_PRIVATE_KEY")))
		if err != nil {
			return err
		}

		// Transfer some tokens to provers.
		balance, err := taikoToken.BalanceOf(nil, crypto.PubkeyToAddress(ownerPrivKey.PublicKey))
		if err != nil {
			return err
		}
		if balance.Cmp(common.Big0) <= 0 {
			return errors.New("balance is less than or equal to 0")
		}

		opts, err := bind.NewKeyedTransactorWithChainID(ownerPrivKey, l1ChainID)
		if err != nil {
			return err
		}

		proverBalance := new(big.Int).Div(balance, common.Big3)
		if proverBalance.Cmp(common.Big0) <= 0 {
			return errors.New("prover balance is less than or equal to 0")
		}

		_, err = taikoToken.Transfer(opts, crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey), proverBalance)
		if err != nil {
			return err
		}

		_, err = taikoToken.Transfer(
			opts,
			common.HexToAddress(os.Getenv("GUARDIAN_PROVER_MINORITY")),
			new(big.Int).Div(proverBalance, common.Big2),
		)
		if err != nil {
			return err
		}

		_, err = taikoToken.Transfer(
			opts,
			common.HexToAddress(os.Getenv("GUARDIAN_PROVER_CONTRACT")),
			new(big.Int).Div(proverBalance, common.Big2),
		)
		if err != nil {
			return err
		}

		if err = setAllowance(l1client, l1ProverPrivKey, taikoToken); err != nil {
			return err
		}
		if err = setAllowance(l1client, ownerPrivKey, taikoToken); err != nil {
			return err
		}
	}
	return nil
}

func setAllowance(client *ethclient.Client, key *ecdsa.PrivateKey, taikoToken *taikotoken.TaikoToken) error {
	decimal, err := taikoToken.Decimals(nil)
	if err != nil {
		return err
	}

	var bigInt = new(big.Int).Exp(big.NewInt(1_000_000_000), new(big.Int).SetUint64(uint64(decimal)), nil)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return err
	}

	tx, err := taikoToken.Approve(auth, common.HexToAddress(os.Getenv("TAIKO_L1")), bigInt)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("approve failed, tx hash: %s", tx.Hash().String())
	}
	return nil
}
