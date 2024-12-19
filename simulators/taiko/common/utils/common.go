package utils

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"math/rand/v2"
	"os"
	"taiko/bindings/proverset"
	"taiko/bindings/taikol1"
	"taiko/bindings/taikotoken"
	"taiko/params"
)

func DeployContracts(ctx context.Context, l1Url, l2Url string) error {
	l1cli, err := ethclient.DialContext(ctx, l1Url)
	if err != nil {
		return err
	}
	chainID, err := l1cli.ChainID(ctx)
	if err != nil {
		return err
	}

	l2cli, err := ethclient.DialContext(ctx, l2Url)
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

		if err := l1cli.SendTransaction(ctx, signedTx); err != nil {
			return err
		}
		fmt.Println("successfully send tx, hash: ", signedTx.Hash().String())
		signedTxs = append(signedTxs, signedTx)
	}

	// Wait the latest tx mined.
	for _, tx := range signedTxs {
		if tx.To() == nil {
			_, err := bind.WaitDeployed(context.Background(), l1cli, tx)
			if err != nil {
				return fmt.Errorf("failed to wait deployed: %v", err)
			}
		} else {
			receipt, err := bind.WaitMined(context.Background(), l1cli, tx)
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
	envs["L1_HTTP"] = l1Url
	envs["L2_HTTP"] = l2Url
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return initTaikoContract(l1cli, l2cli)
}

func initL2Genesis(l1cli, l2cli *ethclient.Client, ownerAuth *bind.TransactOpts) error {
	taikoL1, err := taikol1.NewTaikoL1(common.HexToAddress(os.Getenv("TAIKO_L1")), l1cli)
	if err != nil {
		return err
	}
	genesisHeader, err := l2cli.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		return err
	}

	fmt.Println("l2genesis hash: ", genesisHeader.Hash().String())

	_, err = taikoL1.InitL2Genesis(ownerAuth, genesisHeader.Hash())
	return err
}

// InitTaikoContract init taiko contracts.
func initTaikoContract(l1cli, l2cli *ethclient.Client) error {
	l1ChainID, err := l1cli.ChainID(context.Background())
	if err != nil {
		return err
	}

	taikoToken, err := taikotoken.NewTaikoToken(common.HexToAddress(os.Getenv("TAIKO_TOKEN")), l1cli)
	if err != nil {
		return err
	}

	proverAuth, err := getAuth("L1_PROVER_PRIV_KEY", l1ChainID)
	if err != nil {
		return err
	}

	proposerAuth, err := getAuth("L1_PROPOSER_PRIV_KEY", l1ChainID)
	if err != nil {
		return err
	}

	ownerAuth, err := getAuth("L1_CONTRACT_OWNER_PRIVATE_KEY", l1ChainID)
	if err != nil {
		return err
	}

	if err = initL2Genesis(l1cli, l2cli, ownerAuth); err != nil {
		return err
	}

	decimal, err := taikoToken.Decimals(nil)
	if err != nil {
		return err
	}
	allow := new(big.Int).Exp(big.NewInt(1_000_000_100), new(big.Int).SetUint64(uint64(decimal)), nil)
	fmt.Println(decimal, allow.String())

	// Transfer some tokens to provers.
	balance, err := taikoToken.BalanceOf(nil, ownerAuth.From)
	if err != nil {
		return err
	}
	if balance.Cmp(common.Big0) <= 0 {
		return errors.New("balance is less than or equal to 0")
	}

	proverBalance := new(big.Int).Div(balance, common.Big32)
	if proverBalance.Cmp(common.Big0) <= 0 {
		return errors.New("prover balance is less than or equal to 0")
	}

	if os.Getenv("IS_GUARDIAN") == "true" {
		_, err = taikoToken.Transfer(ownerAuth, proposerAuth.From, proverBalance)
		if err != nil {
			return err
		}
		if err = transferTaikoToken(taikoToken, ownerAuth, "GUARDIAN_PROVER_MINORITY", proverBalance); err != nil {
			return err
		}
		if err = transferTaikoToken(taikoToken, ownerAuth, "GUARDIAN_PROVER_CONTRACT", proverBalance); err != nil {
			return err
		}
	} else {
		if err = transferTaikoToken(taikoToken, ownerAuth, "PROVER_SET", proverBalance); err != nil {
			return err
		}
		if err = enableProver(l1cli, ownerAuth, proposerAuth.From, true); err != nil {
			return err
		}
		if err = enableProver(l1cli, ownerAuth, proverAuth.From, true); err != nil {
			return err
		}
	}

	if err = setAllowance(l1cli, proverAuth, taikoToken); err != nil {
		return err
	}
	if err = setAllowance(l1cli, ownerAuth, taikoToken); err != nil {
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

func setAllowance(client *ethclient.Client, auth *bind.TransactOpts, taikoToken *taikotoken.TaikoToken) error {
	taikoL1 := os.Getenv("TAIKO_L1")
	if taikoL1 == "" {
		return fmt.Errorf("TAIKO_L1 variable is empty")
	}
	decimal, err := taikoToken.Decimals(nil)
	if err != nil {
		return err
	}

	var bigInt = new(big.Int).Exp(big.NewInt(1_000_000_000), new(big.Int).SetUint64(uint64(decimal)), nil)

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

func getAuth(key string, chainID *big.Int) (*bind.TransactOpts, error) {
	ownerPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv(key)))
	if err != nil {
		return nil, err
	}
	return bind.NewKeyedTransactorWithChainID(ownerPrivKey, chainID)
}

// StringToBytes32 converts the given string to [32]byte.
func StringToBytes32(str string) [32]byte {
	var b [32]byte
	copy(b[:], []byte(str))

	return b
}

// compress compresses the given txList bytes using zlib.
func compress(txListBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	if _, err := w.Write(txListBytes); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// EncodeAndCompressTxList encodes and compresses the given transactions list.
func EncodeAndCompressTxList(txs types.Transactions) ([]byte, error) {
	b, err := rlp.EncodeToBytes(txs)
	if err != nil {
		return nil, err
	}

	return compress(b)
}

func CreateL2Txs(
	ctx context.Context,
	l2cli *rpc.EthClient,
	send bool,
) (types.Transactions, error) {
	var txs types.Transactions
	for _, auth := range params.L2Auths {
		nonce, err := l2cli.PendingNonceAt(ctx, auth.From)
		if err != nil {
			return nil, fmt.Errorf("cannot get nonce: %v", err)
		}
		to := common.BigToAddress(big.NewInt(rand.Int64()))
		signedTx, err := auth.Signer(auth.From, types.NewTx(&types.DynamicFeeTx{
			To:        &to,
			Nonce:     nonce,
			Value:     big.NewInt(0),
			GasTipCap: new(big.Int).SetUint64(10 * 1e9),
			GasFeeCap: new(big.Int).SetUint64(20 * 1e9),
			Gas:       2_100_000,
			Data:      nil,
		}))

		if send {
			if err = l2cli.SendTransaction(ctx, signedTx); err != nil {
				return nil, fmt.Errorf("cannot send transaction:, address: %s, err: %v", auth.From, err)
			}
		}

		txs = append(txs, signedTx)
	}
	return txs, nil
}
