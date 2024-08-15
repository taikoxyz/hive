package execution_config

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"taiko/common/config"
	taparams "taiko/params"
	"testing"
)

func TestBuildChainConfig(t *testing.T) {
	slotsPerEpoch := uint64(32)
	slotTime := uint64(12)
	beaconChainGenesisTime := uint64(1634025600)
	ttd := big.NewInt(200)
	chainConfig, err := BuildChainConfig(ttd, beaconChainGenesisTime, slotsPerEpoch, slotTime, &config.ForkConfig{
		AltairForkEpoch:    big.NewInt(0),
		BellatrixForkEpoch: big.NewInt(200),
		CapellaForkEpoch:   nil,
		DenebForkEpoch:     nil,
	})
	if err != nil {
		t.Fatalf("Error producing chainConfig: %v", err)
	}
	if chainConfig.ShanghaiTime != nil {
		t.Fatal("ShanghaiTime is not nil")
	}
	if chainConfig.CancunTime != nil {
		t.Fatal("CancunTime is not nil")
	}
	if chainConfig.TerminalTotalDifficulty.Cmp(ttd) != 0 {
		t.Fatalf("Incorrect TerminalTotalDifficulty is not %d", ttd)
	}

	// Shanghai Chain Config Test
	ttd = big.NewInt(0)
	chainConfig, err = BuildChainConfig(ttd, beaconChainGenesisTime, slotsPerEpoch, slotTime, &config.ForkConfig{
		AltairForkEpoch:    big.NewInt(0),
		BellatrixForkEpoch: big.NewInt(200),
		CapellaForkEpoch:   big.NewInt(200),
		DenebForkEpoch:     nil,
	})
	if err != nil {
		t.Fatalf("Error producing shanghaiChainConfig: %v", err)
	}
	if chainConfig.ShanghaiTime == nil {
		t.Fatal("ShanghaiTime is nil")
	}
	if *chainConfig.ShanghaiTime != beaconChainGenesisTime+(slotsPerEpoch*slotTime*200) {
		t.Fatalf("Incorrect ShanghaiTime is not %d", beaconChainGenesisTime+(slotsPerEpoch*slotTime))
	}
	if chainConfig.CancunTime != nil {
		t.Fatal("CancunTime is not nil")
	}
	if chainConfig.TerminalTotalDifficulty.Cmp(ttd) != 0 {
		t.Fatalf("Incorrect TerminalTotalDifficulty is not %d", ttd)
	}

	// Incorrectly configured epoch
	ttd = big.NewInt(0)
	chainConfig, err = BuildChainConfig(ttd, beaconChainGenesisTime, slotsPerEpoch, slotTime, &config.ForkConfig{
		AltairForkEpoch:    big.NewInt(0),
		BellatrixForkEpoch: big.NewInt(200),
		CapellaForkEpoch:   big.NewInt(199),
		DenebForkEpoch:     nil,
	})
	if err == nil {
		t.Fatalf("Expected error producing chainConfig")
	}
}

func TestGetGenesisFromFile(t *testing.T) {
	// Load params.yml
	beaconConfig, err := params.UnmarshalConfig(taparams.ConfigContent, nil)
	if err != nil {
		panic(err)
	}

	var genesis core.Genesis
	// Load genesis.json
	if err = json.Unmarshal(taparams.GenesisContent, &genesis); err != nil {
		panic(err)
	}

	genesisState := &GenesisState{
		ForkName:         "deneb",
		BeaconConfig:     beaconConfig,
		NumValidators:    64,
		GenesisTimeDelay: 3,
		Genesis:          &genesis,
	}

	executionGenesis, _, err := BuildExecutionGenesis(genesisState)
	assert.NoError(t, err)

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
	encoded, err := marshalFn(executionGenesis.GenesisState)
	assert.NoError(t, err)
	t.Log("encoded length: ", len(encoded))
}
