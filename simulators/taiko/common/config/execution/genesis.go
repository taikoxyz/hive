package execution_config

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/runtime/interop"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
	"github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
)

var (
	log = logrus.WithField("prefix", "genesis")
)

type GenesisState struct {
	ForkName          string
	NumValidators     uint64
	GenesisTime       uint64
	GenesisTimeDelay  uint64
	OverrideEth1Data  bool
	ExecutionEndpoint string
	BeaconConfig      *params.BeaconChainConfig
	Genesis           *core.Genesis
}

// Represents a json object of hex string and uint64 values for
// validators on Ethereum. This file can be generated using the official staking-deposit-cli.
type depositDataJSON struct {
	PubKey                string `json:"pubkey"`
	Amount                uint64 `json:"amount"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	DepositDataRoot       string `json:"deposit_data_root"`
	Signature             string `json:"signature"`
}

func GenerateGenesis(generateGenesisStateFlags *GenesisState) (int, *core.Genesis, error) {
	f := generateGenesisStateFlags
	if err := params.SetActive(f.BeaconConfig.Copy()); err != nil {
		return 0, nil, err
	}

	v, err := version.FromString(f.ForkName)
	if err != nil {
		return 0, nil, err
	}
	gen := f.Genesis
	if gen != nil {
		// set timestamps for genesis and shanghai fork
		gen.Timestamp = generateGenesisStateFlags.GenesisTime
		gen.Config.ShanghaiTime = interop.GethShanghaiTime(generateGenesisStateFlags.GenesisTime, params.BeaconConfig())
		//gen.Config.CancunTime = interop.GethCancunTime(f.GenesisTime, params.BeaconConfig())
		gen.Config.CancunTime = interop.GethCancunTime(generateGenesisStateFlags.GenesisTime, params.BeaconConfig())
		log.
			WithField("shanghai", fmt.Sprintf("%d", *gen.Config.ShanghaiTime)).
			WithField("cancun", fmt.Sprintf("%d", *gen.Config.CancunTime)).
			Info("setting fork geth times")
		if v > version.Altair {
			// set ttd to zero so EL goes post-merge immediately
			gen.Config.TerminalTotalDifficulty = big.NewInt(0)
			gen.Config.TerminalTotalDifficultyPassed = true
		}
	}

	return v, gen, nil
}

func generateBeaconState(ctx context.Context, generateGenesisStateFlags *GenesisState) (state.BeaconState, *core.Genesis, error) {
	f := generateGenesisStateFlags
	v, gen, err := GenerateGenesis(f)
	if err != nil {
		return nil, nil, err
	}

	if f.GenesisTime == 0 {
		f.GenesisTime = uint64(time.Now().Unix())
		log.Info("No genesis time specified, defaulting to now()")
	}
	log.Infof("Delaying genesis %v by %v seconds", f.GenesisTime, f.GenesisTimeDelay)
	f.GenesisTime += f.GenesisTimeDelay
	log.Infof("Genesis is now %v", f.GenesisTime)

	opts := make([]interop.PremineGenesisOpt, 0)
	nv := f.NumValidators
	if nv == 0 {
		return nil, nil, fmt.Errorf("expected --num-validators > 0 or --deposit-json-file to have been provided")
	}

	gb := gen.ToBlock()

	// TODO: expose the PregenesisCreds option with a cli flag - for now defaulting to no withdrawal credentials at genesis
	genesisState, err := interop.NewPreminedGenesis(ctx, f.GenesisTime, nv, 0, v, gb, opts...)
	if err != nil {
		return nil, nil, err
	}

	if f.OverrideEth1Data {
		log.Print("Overriding Eth1Data with data from execution client")
		conn, err := rpc.Dial(f.ExecutionEndpoint)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"could not dial %s please make sure you are running your execution client",
				f.ExecutionEndpoint)
		}
		client := ethclient.NewClient(conn)
		header, err := client.HeaderByNumber(ctx, big.NewInt(0))
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get header by number")
		}
		t, err := trie.NewTrie(params.BeaconConfig().DepositContractTreeDepth)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not create deposit tree")
		}
		depositRoot, err := t.HashTreeRoot()
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get hash tree root")
		}
		e1d := &ethpb.Eth1Data{
			DepositRoot:  depositRoot[:],
			DepositCount: 0,
			BlockHash:    header.Hash().Bytes(),
		}
		if err := genesisState.SetEth1Data(e1d); err != nil {
			return nil, nil, err
		}
		if err := genesisState.SetEth1DepositIndex(0); err != nil {
			return nil, nil, err
		}
	}

	return genesisState, gen, err
}

func depositEntriesFromJSON(enc []byte) ([][]byte, []*ethpb.Deposit_Data, error) {
	var depositJSON []*depositDataJSON
	if err := json.Unmarshal(enc, &depositJSON); err != nil {
		return nil, nil, err
	}
	dds := make([]*ethpb.Deposit_Data, len(depositJSON))
	roots := make([][]byte, len(depositJSON))
	for i, val := range depositJSON {
		root, data, err := depositJSONToDepositData(val)
		if err != nil {
			return nil, nil, err
		}
		dds[i] = data
		roots[i] = root
	}
	return roots, dds, nil
}

func depositJSONToDepositData(input *depositDataJSON) ([]byte, *ethpb.Deposit_Data, error) {
	root, err := hex.DecodeString(strings.TrimPrefix(input.DepositDataRoot, "0x"))
	if err != nil {
		return nil, nil, err
	}
	pk, err := hex.DecodeString(strings.TrimPrefix(input.PubKey, "0x"))
	if err != nil {
		return nil, nil, err
	}
	creds, err := hex.DecodeString(strings.TrimPrefix(input.WithdrawalCredentials, "0x"))
	if err != nil {
		return nil, nil, err
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(input.Signature, "0x"))
	if err != nil {
		return nil, nil, err
	}
	return root, &ethpb.Deposit_Data{
		PublicKey:             pk,
		WithdrawalCredentials: creds,
		Amount:                input.Amount,
		Signature:             sig,
	}, nil
}
