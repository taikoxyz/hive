package execution_config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"math/big"
	"taiko/common/config"
	consensus_config "taiko/common/config/consensus"
	"taiko/common/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/hive/hivesim"
	cl "taiko/common/config/consensus"
)

var (
	CLIQUE_PERIOD_DEFAULT        = uint64(2)
	DEFAULT_CLIQUE_PRIVATE_KEY   = "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622"
	DEFAULT_CLIQUE_MINER_ADDRESS = "123463a4B065722E99115D6c222f267d9cABb524"
	DEFAULT_ETHASH_MINER_ADDRESS = "1212121212121212121212121212121212121212"
)

type ExecutionConsensus interface {
	HiveParams(int) hivesim.Params
	DifficultyPerBlock() *big.Int
	SecondsPerBlock() uint64
}

type ExecutionEthashConsensus struct {
	MinerAddress string
	MiningNodes  int
}

func (c ExecutionEthashConsensus) Configure(*core.Genesis) error {
	// Nothing to do here...
	return nil
}

func (c ExecutionEthashConsensus) HiveParams(node int) hivesim.Params {
	if c.MinerAddress == "" {
		c.MinerAddress = DEFAULT_ETHASH_MINER_ADDRESS
	}
	if c.MiningNodes == 0 {
		// Default is that only one node is a miner
		c.MiningNodes = 1
	}
	if node < c.MiningNodes {
		return hivesim.Params{"HIVE_MINER": c.MinerAddress}
	}
	return hivesim.Params{}
}

func (c ExecutionEthashConsensus) DifficultyPerBlock() *big.Int {
	// Approximately 0x20000
	return big.NewInt(131072)
}

func (c ExecutionEthashConsensus) SecondsPerBlock() uint64 {
	// It is really hard to approxmate this value
	return 10
}

// A pre-existing chain is imported by the client, and it is not
// expected that the client mines or produces any blocks.
type ExecutionPreChain struct{}

func (c ExecutionPreChain) Configure(*core.Genesis) error {
	return nil
}

func (c ExecutionPreChain) HiveParams(node int) hivesim.Params {
	return hivesim.Params{}
}

func (c ExecutionPreChain) DifficultyPerBlock() *big.Int {
	// Approximately 0x20000
	return big.NewInt(131072)
}

func (c ExecutionPreChain) SecondsPerBlock() uint64 {
	return 1
}

// A pre-existing chain is imported by the client, and it is not
// expected that the client mines or produces any blocks.
type ExecutionPostMergeGenesis struct{}

func (c ExecutionPostMergeGenesis) Configure(*core.Genesis) error {
	return nil
}

func (c ExecutionPostMergeGenesis) HiveParams(node int) hivesim.Params {
	return hivesim.Params{}
}

func (c ExecutionPostMergeGenesis) DifficultyPerBlock() *big.Int {
	return big.NewInt(0)
}

func (c ExecutionPostMergeGenesis) SecondsPerBlock() uint64 {
	return 12
}

type ExecutionCliqueConsensus struct {
	CliquePeriod     uint64
	CliquePrivateKey string
	CliqueAddress    string
}

func (c ExecutionCliqueConsensus) Configure(genesis *core.Genesis) error {
	if c.CliquePeriod == 0 {
		c.CliquePeriod = CLIQUE_PERIOD_DEFAULT
	}
	if c.CliqueAddress == "" {
		c.CliqueAddress = DEFAULT_CLIQUE_MINER_ADDRESS
	}
	genesis.Config.Clique = &params.CliqueConfig{
		Period: c.CliquePeriod,
		Epoch:  0,
	}

	genesis.ExtraData = make([]byte, utils.ExtraVanity+utils.ExtraSeal+common.AddressLength)
	minerAddr := common.HexToAddress(c.CliqueAddress)
	copy(genesis.ExtraData[utils.ExtraVanity:], minerAddr[:])
	return nil
}

func (c ExecutionCliqueConsensus) HiveParams(node int) hivesim.Params {
	//if node > 0 {
	//	return hivesim.Params{}
	//}
	if c.CliquePrivateKey == "" || c.CliqueAddress == "" {
		c.CliquePrivateKey = DEFAULT_CLIQUE_PRIVATE_KEY
		c.CliqueAddress = DEFAULT_CLIQUE_MINER_ADDRESS
	}
	return hivesim.Params{
		"HIVE_TAIKO2_CLIQUE_PRIVATEKEY": c.CliquePrivateKey,
		"HIVE_TAIKO2_CLIQUE_ADDRESS":    c.CliqueAddress,
	}
}

func (c ExecutionCliqueConsensus) DifficultyPerBlock() *big.Int {
	return big.NewInt(2)
}

func (c ExecutionCliqueConsensus) SecondsPerBlock() uint64 {
	if c.CliquePeriod == 0 {
		return CLIQUE_PERIOD_DEFAULT
	}
	return c.CliquePeriod
}

func BuildChainConfig(
	ttd *big.Int,
	beaconChainGenesisTime uint64,
	slotsPerEpoch uint64,
	secondsPerSlot uint64,
	config *config.ForkConfig,
) (*params.ChainConfig, error) {
	chainConfig := &params.ChainConfig{
		ChainID:                 big.NewInt(7),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: ttd,
		Clique:                  nil,
	}

	// Configure post-merge forks
	var (
		previousForkTime = config.BellatrixForkEpoch
		previousFork     = "bellatrix"
	)
	for _, forkConfig := range []struct {
		ForkName       string
		BeaconForkTime *big.Int
		ChainConfig    **uint64
	}{
		{
			ForkName:       "capella",
			BeaconForkTime: config.CapellaForkEpoch,
			ChainConfig:    &chainConfig.ShanghaiTime,
		},
		{
			ForkName:       "deneb",
			BeaconForkTime: config.DenebForkEpoch,
			ChainConfig:    &chainConfig.CancunTime,
		},
	} {
		if forkConfig.BeaconForkTime != nil {
			if previousForkTime == nil {
				return nil, fmt.Errorf("fork '%s' has a time but previous fork '%s' does not", forkConfig.ForkName, previousFork)
			}
			if forkConfig.BeaconForkTime.Cmp(previousForkTime) < 0 {
				return nil, fmt.Errorf("fork '%s' has a time before previous fork '%s'", forkConfig.ForkName, previousFork)
			}
			timestamp := beaconChainGenesisTime + (forkConfig.BeaconForkTime.Uint64() * secondsPerSlot * slotsPerEpoch)
			*forkConfig.ChainConfig = &timestamp
		}
		previousForkTime = forkConfig.BeaconForkTime
		previousFork = forkConfig.ForkName
	}
	return chainConfig, nil
}

type ExecutionGenesis struct {
	GenesisState   state.BeaconState
	Genesis        *core.Genesis
	Block          *types.Block
	Hash           common.Hash
	DepositAddress string
}

func BuildExecutionGenesis(generateGenesisStateFlags *GenesisState) (*ExecutionGenesis, *consensus_config.Spec, error) {
	genesisState, genesis, err := generateBeaconState(context.Background(), generateGenesisStateFlags)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("validator length: ", len(genesisState.Validators()))

	genesisBlock := genesis.ToBlock()
	return &ExecutionGenesis{
			GenesisState:   genesisState,
			Genesis:        genesis,
			Block:          genesisBlock,
			Hash:           genesisBlock.Hash(),
			DepositAddress: generateGenesisStateFlags.BeaconConfig.DepositContractAddress,
		}, &cl.Spec{
			BeaconChainConfig: *generateGenesisStateFlags.BeaconConfig,
		}, nil
}

func (genesis *ExecutionGenesis) NetworkID() uint64 {
	return 7
}

func (genesis *ExecutionGenesis) ChainID() uint64 {
	return genesis.Genesis.Config.ChainID.Uint64()
}

func (genesis *ExecutionGenesis) IsPostMerge() bool {
	return genesis.Block.Difficulty().Cmp(genesis.Genesis.Config.TerminalTotalDifficulty) >= 0
}

func (conf *ExecutionGenesis) ToParams() hivesim.Params {
	params := hivesim.Params{
		"HIVE_DEPOSIT_CONTRACT_ADDRESS": conf.DepositAddress,
		"HIVE_NETWORK_ID":               fmt.Sprintf("%d", conf.NetworkID()),
		"HIVE_CHAIN_ID":                 conf.Genesis.Config.ChainID.String(),
		"HIVE_FORK_HOMESTEAD":           conf.Genesis.Config.HomesteadBlock.String(),
		//"HIVE_FORK_DAO_BLOCK":           conf.Genesis.Config.DAOForkBlock.String(),  // nil error, not used anyway
		"HIVE_FORK_TANGERINE":            conf.Genesis.Config.EIP150Block.String(),
		"HIVE_FORK_SPURIOUS":             conf.Genesis.Config.EIP155Block.String(), // also eip558
		"HIVE_FORK_BYZANTIUM":            conf.Genesis.Config.ByzantiumBlock.String(),
		"HIVE_FORK_CONSTANTINOPLE":       conf.Genesis.Config.ConstantinopleBlock.String(),
		"HIVE_FORK_PETERSBURG":           conf.Genesis.Config.PetersburgBlock.String(),
		"HIVE_FORK_ISTANBUL":             conf.Genesis.Config.IstanbulBlock.String(),
		"HIVE_FORK_MUIRGLACIER":          conf.Genesis.Config.MuirGlacierBlock.String(),
		"HIVE_FORK_BERLIN":               conf.Genesis.Config.BerlinBlock.String(),
		"HIVE_FORK_LONDON":               conf.Genesis.Config.LondonBlock.String(),
		"HIVE_FORK_ARROWGLACIER":         conf.Genesis.Config.ArrowGlacierBlock.String(),
		"HIVE_MERGE_BLOCK_ID":            conf.Genesis.Config.MergeNetsplitBlock.String(),
		"HIVE_TERMINAL_TOTAL_DIFFICULTY": conf.Genesis.Config.TerminalTotalDifficulty.String(),
	}
	if conf.Genesis.Config.ShanghaiTime != nil {
		params["HIVE_SHANGHAI_TIMESTAMP"] = fmt.Sprint(*conf.Genesis.Config.ShanghaiTime)
	}
	if conf.Genesis.Config.CancunTime != nil {
		params["HIVE_CANCUN_TIMESTAMP"] = fmt.Sprint(*conf.Genesis.Config.CancunTime)
	}
	if conf.Genesis.Config.Clique != nil {
		params["HIVE_CLIQUE_PERIOD"] = fmt.Sprint(conf.Genesis.Config.Clique.Period)
	}
	return params
}

func ExecutionBundle(genesis *core.Genesis) (hivesim.StartOption, error) {
	out, err := json.Marshal(genesis)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize genesis state: %v", err)
	}
	return hivesim.WithDynamicFile(
		"/hive/input/genesis.json",
		config.BytesSource(out),
	), nil
}

func ChainBundle(chain []*types.Block) (hivesim.StartOption, error) {
	var buf bytes.Buffer
	for _, block := range chain {
		if err := block.EncodeRLP(&buf); err != nil {
			return nil, err
		}
	}
	return hivesim.WithDynamicFile(
		"/chain.rlp",
		config.BytesSource(buf.Bytes()),
	), nil
}
