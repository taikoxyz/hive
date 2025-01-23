package params

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/hive/hivesim"
	"github.com/joho/godotenv"
	"math/big"
	"os"
	"strconv"
	"time"
)

//go:embed config.yml
var ConfigContent []byte

//go:embed genesis.json
var GenesisContent []byte

type TaikoVersion string

var (
	envParams = hivesim.Params{}

	OntakeVersion  = TaikoVersion("ontake")
	PacayaVersion  = TaikoVersion("pacaya")
	CurrentVersion = TaikoVersion(os.Getenv("TAIKO_VERSION"))

	ContractTxs = make([]*types.Transaction, 0)
	PrivateKeys []*ecdsa.PrivateKey
	L1Auths     []*bind.TransactOpts
	L2Auths     []*bind.TransactOpts

	ZeroAddress = common.Address{}
)

func init() {
	var (
		txsContent []byte
		vals       map[string]string
		err        error
	)

	// Load env params.
	vals, err = godotenv.Read(fmt.Sprintf("%s/.env", CurrentVersion))
	if err != nil {
		panic(err)
	}
	txsContent, err = os.ReadFile(fmt.Sprintf("%s/contract_txs.json", CurrentVersion))
	if err != nil {
		panic(err)
	}

	for k, v := range vals {
		envParams[k] = v
		// set envs
		_ = os.Setenv(k, v)
	}
	err = os.WriteFile("/tmp/jwt.hex", []byte("c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7"), 0644)
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("JWT_SECRET", "/tmp/jwt.hex")

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

	for _, privateKeyHex := range []string{
		"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", // 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
		"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a", // 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
		"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", // 0x90F79bf6EB2c4f870365E785982E1f101E93b906
		"0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a", // 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65
		"0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba", // 0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc
	} {
		priv, err := crypto.ToECDSA(common.FromHex(privateKeyHex))
		if err != nil {
			panic(err)
		}
		PrivateKeys = append(PrivateKeys, priv)

		auth, err := bind.NewKeyedTransactorWithChainID(priv, big.NewInt(32382))
		if err != nil {
			panic(err)
		}
		L1Auths = append(L1Auths, auth)

		auth, err = bind.NewKeyedTransactorWithChainID(priv, big.NewInt(167001))
		if err != nil {
			panic(err)
		}
		L2Auths = append(L2Auths, auth)
	}
}

func ParamByKey(key string) string {
	return envParams[key]
}

func ParamToBool(key string) bool {
	parsed, err := strconv.ParseBool(envParams[key])
	if err != nil {
		return false
	}
	return parsed
}

func ParamToDuration(key string) time.Duration {
	parsed, err := time.ParseDuration(envParams[key])
	if err != nil {
		return 0
	}
	return parsed
}

func ParamToAddress(key string) common.Address {
	return common.HexToAddress(envParams[key])
}

func ParamToBytes(key string) []byte {
	return common.FromHex(envParams[key])
}

func ParamToUint64(key string) uint64 {
	parsed, err := strconv.ParseUint(envParams[key], 0, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func EnvParams() hivesim.Params {
	return envParams.Copy()
}

func SetEnvParams(key, value string) {
	if value == "" {
		delete(envParams, key)
	} else {
		envParams[key] = value
		// set envs
		_ = os.Setenv(key, value)
	}
}
