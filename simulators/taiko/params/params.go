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

//go:embed .env
var envContent []byte

//go:embed l1contract_txs.txt
var txsContent []byte

//go:embed config.yml
var ConfigContent []byte

//go:embed genesis.json
var GenesisContent []byte

type Auths struct {
	Address    common.Address
	SecretKey  string
	PrivateKey *ecdsa.PrivateKey
}

var (
	envParams   = hivesim.Params{}
	ContractTxs = make([]*types.Transaction, 0)
	L1Auths     []*bind.TransactOpts
	ChainAuths  []*Auths
)

func init() {
	// Load env params.
	vals, err := godotenv.UnmarshalBytes(envContent)
	if err != nil {
		panic(err)
	}
	for k, v := range vals {
		envParams[k] = v
		fmt.Printf("%s=%s\n", k, v)
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
		"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
		"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", // 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
		"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a", // 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
		"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", // 0x90F79bf6EB2c4f870365E785982E1f101E93b906
		"0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a", // 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65
		"0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba", // 0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc
		"0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e", // 0x976EA74026E726554dB657fA54763abd0C3a0aa9
		"0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356", // 0x14dC79964da2C08b23698B3D3cc7Ca32193d9955
		"0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97", // 0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f
		"0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6", // 0xa0Ee7A142d267C1f36714E4a8F75612F20a79720
	} {
		priv, err := crypto.ToECDSA(common.FromHex(privateKeyHex))
		if err != nil {
			panic(err)
		}
		ChainAuths = append(ChainAuths, &Auths{
			Address:    crypto.PubkeyToAddress(priv.PublicKey),
			SecretKey:  privateKeyHex,
			PrivateKey: priv,
		})

		auth, err := bind.NewKeyedTransactorWithChainID(priv, big.NewInt(32382))
		if err != nil {
			panic(err)
		}
		L1Auths = append(L1Auths, auth)
	}
}

func ParamByKey(key string) string {
	return envParams[key]
}

func ParamToPriv(key string) *ecdsa.PrivateKey {
	priv, err := crypto.ToECDSA(common.FromHex(envParams[key]))
	if err != nil {
		return nil
	}
	return priv
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
	}
	// set envs
	_ = os.Setenv(key, value)
}
