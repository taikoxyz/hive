package utils

import (
	"bytes"
	"compress/zlib"
	"context"
	crand "crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"math/rand/v2"
	"os"
	"taiko/bindings/ontake/proverset"
	"taiko/bindings/pacaya/taikotoken"
	"taiko/params"
)

type ERC20API interface {
	Transfer(auth *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error)
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)
}

func DeployContracts(envs hivesim.Params, l1cli *rpc.EthClient) error {
	chainID := l1cli.ChainID
	sk, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return fmt.Errorf("failed to get private key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	if err != nil {
		return fmt.Errorf("chainID: %s, err: %v", chainID.String(), err)
	}

	signedTxs := make([]*types.Transaction, 0, len(params.ContractTxs))
	for _, tx := range params.ContractTxs {
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			return fmt.Errorf("failed to sign tx: %v", err)
		}

		if err := l1cli.SendTransaction(context.Background(), signedTx); err != nil {
			return fmt.Errorf("failed to send tx, %s: %v", signedTx.Hash().String(), err)
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
				return fmt.Errorf("failed to call contract, hash: %s, status: %d", tx.Hash().String(), receipt.Status)
			}
		}
	}

	return initOntakeContracts(envs, l1cli)
}

// InitTaikoContract init taiko contracts.
func initOntakeContracts(envs hivesim.Params, l1cli *rpc.EthClient) error {
	l1ChainID := l1cli.ChainID

	ownerAuth, err := getAuth(envs["L1_CONTRACT_OWNER_PRIVATE_KEY"], l1ChainID)
	if err != nil {
		return fmt.Errorf("failed to get owner auth: %v", err)
	}

	taikoToken, err := taikotoken.NewTaikoToken(common.HexToAddress(envs["TAIKO_TOKEN"]), l1cli)
	if err != nil {
		return err
	}

	// Transfer some tokens to proposers.
	balance, err := taikoToken.BalanceOf(nil, ownerAuth.From)
	if err != nil {
		return err
	}
	bls := new(big.Int).Div(balance, common.Big256)

	if os.Getenv("IS_GUARDIAN") == "true" {
		if err = erc20Transfer(l1cli, taikoToken, ownerAuth, common.HexToAddress(envs["GUARDIAN_PROVER_MINORITY"]), bls); err != nil {
			return err
		}
		if err = erc20Transfer(l1cli, taikoToken, ownerAuth, common.HexToAddress(envs["GUARDIAN_PROVER_CONTRACT"]), bls); err != nil {
			return err
		}
	}

	proverSetAddr := common.HexToAddress(envs["PROVER_SET"])
	taikoInboxAddr := common.HexToAddress(envs["TAIKO_INBOX"])

	if proverSetAddr != (common.Address{}) {
		proverSet, err := proverset.NewProverSet(proverSetAddr, l1cli)
		if err != nil {
			return err
		}
		for _, auth := range params.L1Auths[1:5] {
			// Enable prover set.
			if err = enableProverSet(l1cli, proverSet, ownerAuth, auth.From); err != nil {
				return err
			}
		}
		// Transfer some tokens to proposers.
		if err = erc20Transfer(l1cli, taikoToken, ownerAuth, proverSetAddr, bls); err != nil {
			return err
		}
		if err = erc20Approve(l1cli, taikoToken, ownerAuth, taikoInboxAddr, bls); err != nil {
			return err
		}
	} else {
		for _, auth := range params.L1Auths[1:5] {
			// Transfer some tokens to proposers.
			if err = erc20Transfer(l1cli, taikoToken, ownerAuth, auth.From, bls); err != nil {
				return err
			}
			if err = erc20Approve(l1cli, taikoToken, auth, taikoInboxAddr, bls); err != nil {
				return err
			}
		}
	}

	return nil
}

func erc20Transfer(l1cli *rpc.EthClient, token ERC20API, auth *bind.TransactOpts, to common.Address, balance *big.Int) error {
	if to == (common.Address{}) {
		return nil
	}
	tx, err := token.Transfer(auth, to, balance)
	if _, err = bind.WaitMined(context.Background(), l1cli, tx); err != nil {
		return err
	}

	bls, err := token.BalanceOf(nil, to)
	if err != nil {
		return err
	}
	if bls.Cmp(balance) < 0 {
		return fmt.Errorf("failed to transfer, to: %s, expect: %s, actual: %s", to.String(), balance.String(), bls.String())
	}

	return err
}

func erc20Approve(l1cli *rpc.EthClient, token ERC20API, auth *bind.TransactOpts, spender common.Address, amount *big.Int) error {
	if spender == (common.Address{}) {
		return nil
	}
	tx, err := token.Approve(auth, spender, amount)
	_, err = bind.WaitMined(context.Background(), l1cli, tx)

	bls, err := token.Allowance(nil, auth.From, spender)
	if err != nil {
		return err
	}
	if bls.Cmp(amount) < 0 {
		return fmt.Errorf("failed to approve, spender: %s, expect: %s, actual: %s", spender.String(), amount.String(), bls.String())
	}

	return err
}

func enableProverSet(l1cli *rpc.EthClient, proverSet *proverset.ProverSet, auth *bind.TransactOpts, prover common.Address) error {
	tx, err := proverSet.EnableProver(auth, prover, true)
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(context.Background(), l1cli, tx)
	return err
}

func getAuth(val string, chainID *big.Int) (*bind.TransactOpts, error) {
	ownerPrivKey, err := crypto.ToECDSA(common.FromHex(val))
	if err != nil {
		return nil, err
	}
	return bind.NewKeyedTransactorWithChainID(ownerPrivKey, chainID)
}

// compress compresses the given txList bytes using zlib.
func compress(txListBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	if _, err := w.Write(txListBytes); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
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
	for _, atr := range params.ChainAuths[5:] {
		auth, err := bind.NewKeyedTransactorWithChainID(atr.PrivateKey, l2cli.ChainID)
		if err != nil {
			return nil, err
		}
		nonce, err := l2cli.PendingNonceAt(context.Background(), auth.From)
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

// RandomHash generates a random blob of data and returns it as a hash.
func RandomHash() common.Hash {
	var hash common.Hash
	if n, err := crand.Read(hash[:]); n != common.HashLength || err != nil {
		panic(err)
	}
	return hash
}

// RandomBytes generates a random bytes.
func RandomBytes(size int) (b []byte) {
	b = make([]byte, size)
	if _, err := crand.Read(b); err != nil {
		log.Crit("Generate random bytes error", "error", err)
	}
	return
}
