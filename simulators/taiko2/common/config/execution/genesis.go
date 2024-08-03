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
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/container/trie"
	"github.com/prysmaticlabs/prysm/v4/encoding/ssz/detect"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/interop"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"math/big"
	"os"
	"strings"
	"time"
)

var (
	log = logrus.WithField("prefix", "genesis")
)

type GenesisState struct {
	DepositJsonFile    string
	ChainConfigFile    string
	ConfigName         string
	NumValidators      uint64
	GenesisTime        uint64
	GenesisTimeDelay   uint64
	OutputSSZ          string
	OutputJSON         string
	OutputYaml         string
	ForkName           string
	OverrideEth1Data   bool
	ExecutionEndpoint  string
	GethGenesisJsonIn  string
	GethGenesisJsonOut string
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

func CliActionGenerateGenesisState(cliCtx context.Context, generateGenesisStateFlags *GenesisState) error {
	outputJson := generateGenesisStateFlags.OutputJSON
	outputYaml := generateGenesisStateFlags.OutputYaml
	outputSSZ := generateGenesisStateFlags.OutputSSZ
	noOutputFlag := outputSSZ == "" && outputJson == "" && outputYaml == ""
	if noOutputFlag {
		return fmt.Errorf("no outputJson, outputYaml, outputSSZ flag(s) specified. At least one is required")
	}

	st, _, err := generateBeaconState(cliCtx, generateGenesisStateFlags)
	if err != nil {
		return fmt.Errorf("could not generate genesis state: %v", err)
	}

	if outputJson != "" {
		if err := writeToOutputFile(outputJson, st, json.Marshal); err != nil {
			return err
		}
	}
	if outputYaml != "" {
		if err := writeToOutputFile(outputYaml, st, yaml.Marshal); err != nil {
			return err
		}
	}
	if outputSSZ != "" {
		type MinimumSSZMarshal interface {
			MarshalSSZ() ([]byte, error)
		}
		marshalFn := func(o interface{}) ([]byte, error) {
			marshaler, ok := o.(MinimumSSZMarshal)
			if !ok {
				return nil, errors.New("not a marshaler")
			}
			return marshaler.MarshalSSZ()
		}
		if err := writeToOutputFile(outputSSZ, st, marshalFn); err != nil {
			return err
		}
	}
	log.Info("Command completed")
	return nil
}

func setGlobalParams(generateGenesisStateFlags *GenesisState) error {
	chainConfigFile := generateGenesisStateFlags.ChainConfigFile
	if chainConfigFile != "" {
		log.Infof("Specified a chain config file: %s", chainConfigFile)
		return params.LoadChainConfigFile(chainConfigFile, nil)
	}
	cfg, err := params.ByName(generateGenesisStateFlags.ConfigName)
	if err != nil {
		return fmt.Errorf("unable to find config using name %s: %v", generateGenesisStateFlags.ConfigName, err)
	}
	return params.SetActive(cfg.Copy())
}

func GenerateGenesis(generateGenesisStateFlags *GenesisState) (int, *core.Genesis, error) {
	if err := setGlobalParams(generateGenesisStateFlags); err != nil {
		return 0, nil, fmt.Errorf("could not set config params: %v", err)
	}

	f := generateGenesisStateFlags
	v, err := version.FromString(f.ForkName)
	if err != nil {
		return 0, nil, err
	}
	gen := &core.Genesis{}
	if generateGenesisStateFlags.GethGenesisJsonIn != "" {
		gbytes, err := os.ReadFile(generateGenesisStateFlags.GethGenesisJsonIn) // #nosec G304
		if err != nil {
			return 0, nil, errors.Wrapf(err, "failed to read %s", generateGenesisStateFlags.GethGenesisJsonIn)
		}
		if err := json.Unmarshal(gbytes, gen); err != nil {
			return 0, nil, err
		}
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
	} else {
		gen = interop.GethTestnetGenesis(generateGenesisStateFlags.GenesisTime, params.BeaconConfig())
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
	if f.DepositJsonFile != "" {
		expanded, err := file.ExpandPath(f.DepositJsonFile)
		if err != nil {
			return nil, nil, err
		}
		log.Printf("reading deposits from JSON at %s", expanded)
		b, err := os.ReadFile(expanded) // #nosec G304
		if err != nil {
			return nil, nil, err
		}
		roots, dds, err := depositEntriesFromJSON(b)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, interop.WithDepositData(dds, roots))
	} else if nv == 0 {
		return nil, nil, fmt.Errorf("expected --num-validators > 0 or --deposit-json-file to have been provided")
	}

	if f.GethGenesisJsonOut != "" {
		gbytes, err := json.MarshalIndent(gen, "", "\t")
		if err != nil {
			return nil, nil, err
		}
		if err = os.WriteFile(f.GethGenesisJsonOut, gbytes, os.ModePerm); err != nil {
			return nil, nil, errors.Wrapf(err, "failed to write %s", f.GethGenesisJsonOut)
		}
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

func writeToOutputFile(
	fPath string,
	data interface{},
	marshalFn func(o interface{}) ([]byte, error),
) error {
	encoded, err := marshalFn(data)
	if err != nil {
		return err
	}
	if err := file.WriteFile(fPath, encoded); err != nil {
		return err
	}
	log.Printf("Done writing genesis state to %s", fPath)
	return nil
}

func loadGenesis(genesisState *GenesisState) (state.BeaconState, *core.Genesis, error) {
	sb, err := os.ReadFile(genesisState.OutputSSZ)
	if err != nil {
		return nil, nil, err
	}

	if len(sb) < (1 << 10) {
		log.WithField("size", fmt.Sprintf("%d bytes", len(sb))).
			Warn("Genesis state is smaller than one 1Kb. This could be an empty file, git lfs metadata file, or corrupt genesis state.")
	}

	vu, err := detect.FromState(sb)
	if err != nil {
		return nil, nil, err
	}
	gs, err := vu.UnmarshalBeaconState(sb)
	if err != nil {
		return nil, nil, err
	}

	gbytes, err := os.ReadFile(genesisState.OutputSSZ)
	if err != nil {
		return nil, nil, err
	}

	gen := &core.Genesis{}
	if err := json.Unmarshal(gbytes, gen); err != nil {
		return nil, nil, err
	}

	return gs, gen, nil
}
