package utils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"taiko/bindings/proverset"
	"taiko/bindings/taikotoken"
	"taiko/params"
)

func DeployContracts(ctx context.Context, url string) error {
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	sk, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	if err != nil {
		return err
	}

	signedTxs := make([]*types.Transaction, 0, len(params.ContractTxs))
	for _, tx := range params.ContractTxs {
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			return err
		}

		if err := client.SendTransaction(ctx, signedTx); err != nil {
			return err
		}
		fmt.Println("successfully send tx, hash: ", signedTx.Hash().String())
		signedTxs = append(signedTxs, signedTx)
	}

	// Wait the latest tx mined.
	for _, tx := range signedTxs {
		if tx.To() == nil {
			_, err := bind.WaitDeployed(context.Background(), client, tx)
			if err != nil {
				return fmt.Errorf("failed to wait deployed: %v", err)
			}
		} else {
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				return fmt.Errorf("failed to wait mined, hash: %s, err: %v", tx.Hash().String(), err)
			}
			if receipt.Status != types.ReceiptStatusSuccessful {
				return fmt.Errorf("failed to call contract, hash: %s", tx.Hash().String())
			}
		}
	}

	// init contracts.
	envs := params.EnvParams()
	envs["L1_HTTP"] = url
	return initTaikoContract(envs)
}

// InitTaikoContract init taiko contracts.
func initTaikoContract(params map[string]string) error {
	for k, v := range params {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
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

	l1ProposerPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIV_KEY")))
	if err != nil {
		return err
	}

	decimal, err := taikoToken.Decimals(nil)
	if err != nil {
		return err
	}
	allow := new(big.Int).Exp(big.NewInt(1_000_000_100), new(big.Int).SetUint64(uint64(decimal)), nil)
	fmt.Println(decimal, allow.String())

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

	proverBalance := new(big.Int).Div(balance, common.Big32)
	if proverBalance.Cmp(common.Big0) <= 0 {
		return errors.New("prover balance is less than or equal to 0")
	}

	if os.Getenv("IS_GUARDIAN") == "true" {
		_, err = taikoToken.Transfer(opts, crypto.PubkeyToAddress(l1ProposerPrivKey.PublicKey), proverBalance)
		if err != nil {
			return err
		}
		if err = transferTaikoToken(taikoToken, opts, "GUARDIAN_PROVER_MINORITY", proverBalance); err != nil {
			return err
		}
		if err = transferTaikoToken(taikoToken, opts, "GUARDIAN_PROVER_CONTRACT", proverBalance); err != nil {
			return err
		}
	} else {
		if err = transferTaikoToken(taikoToken, opts, "PROVER_SET", proverBalance); err != nil {
			return err
		}
		if err = enableProver(l1client, opts, crypto.PubkeyToAddress(l1ProposerPrivKey.PublicKey), true); err != nil {
			return err
		}
		if err = enableProver(l1client, opts, crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey), true); err != nil {
			return err
		}
	}

	if err = setAllowance(l1client, l1ProverPrivKey, taikoToken); err != nil {
		return err
	}
	if err = setAllowance(l1client, ownerPrivKey, taikoToken); err != nil {
		return err
	}

	return nil
}

func transferTaikoToken(taikoToken *taikotoken.TaikoToken, auth *bind.TransactOpts, env string, balance *big.Int) error {
	if os.Getenv(env) == "" {
		return fmt.Errorf("%s varibale is empty", env)
	}
	_, err := taikoToken.Transfer(
		auth,
		common.HexToAddress(os.Getenv(env)),
		balance,
	)
	return err
}

func enableProver(client *ethclient.Client, auth *bind.TransactOpts, _prover common.Address, _isProver bool) error {
	proverSet := os.Getenv("PROVER_SET")
	if proverSet == "" {
		return fmt.Errorf("PROVER_SET variable is empty")
	}
	prover, err := proverset.NewProverSet(common.HexToAddress(proverSet), client)
	if err != nil {
		return err
	}
	_, err = prover.EnableProver(auth, _prover, _isProver)
	return err
}

func setAllowance(client *ethclient.Client, key *ecdsa.PrivateKey, taikoToken *taikotoken.TaikoToken) error {
	taikoL1 := os.Getenv("TAIKO_L1")
	if taikoL1 == "" {
		return fmt.Errorf("TAIKO_L1 variable is empty")
	}
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

	tx, err := taikoToken.Approve(auth, common.HexToAddress(taikoL1), bigInt)
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

// StringToBytes32 converts the given string to [32]byte.
func StringToBytes32(str string) [32]byte {
	var b [32]byte
	copy(b[:], []byte(str))

	return b
}
