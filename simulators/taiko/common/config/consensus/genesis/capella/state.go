package capella

import (
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"taiko/common/config/consensus/genesis/interfaces"
	"taiko/common/utils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
)

type GenesisStateView struct {
	*capella.BeaconStateView
	Spec *common.Spec
}

func NewBeaconStateView(spec *common.Spec) interfaces.StateViewGenesis {
	return &GenesisStateView{
		BeaconStateView: capella.NewBeaconStateView(spec),
		Spec:            spec,
	}
}

func (g *GenesisStateView) ToJson() ([]byte, error) {
	raw, err := g.BeaconStateView.Raw(g.Spec)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(raw, "", "  ")
}

func (g *GenesisStateView) ForkVersion() common.Version {
	return g.Spec.CAPELLA_FORK_VERSION
}

func (g *GenesisStateView) PreviousForkVersion() common.Version {
	return g.Spec.BELLATRIX_FORK_VERSION
}

func (g *GenesisStateView) EmptyBodyRoot() common.Root {
	return capella.BeaconBlockBodyType(configs.Mainnet).
		New().
		HashTreeRoot(tree.GetHashFn())
}

func (g *GenesisStateView) SetGenesisExecutionHeader(
	executionGenesis *types.Block,
) error {
	ttd := uint256.Int(g.Spec.TERMINAL_TOTAL_DIFFICULTY)

	if executionGenesis.Difficulty().Cmp(ttd.ToBig()) < 0 {
		return fmt.Errorf(
			"eth1 genesis difficulty %d is less than terminal total difficulty %d",
			executionGenesis.Difficulty(),
			ttd,
		)
	}

	extra := executionGenesis.Extra()
	if len(extra) <= utils.ExtraVanity+utils.ExtraSeal ||
		(len(extra)-utils.ExtraVanity-utils.ExtraSeal)%ethcommon.AddressLength != 0 {
		return fmt.Errorf(
			"extra data(%d) is not right",
			len(extra),
		)
	}
	if len(executionGenesis.Transactions()) != 0 {
		return fmt.Errorf(
			"expected no transactions in genesis execution payload",
		)
	}
	if len(executionGenesis.Withdrawals()) != 0 {
		return fmt.Errorf(
			"expected no withdrawals in genesis execution payload",
		)
	}

	baseFee, overflow := uint256.FromBig(executionGenesis.BaseFee())
	if overflow {
		return fmt.Errorf("basefee larger than 2^256-1")
	}

	execPayloadHeader := &capella.ExecutionPayloadHeader{
		ParentHash:    common.Root(executionGenesis.ParentHash()),
		FeeRecipient:  common.Eth1Address(executionGenesis.Coinbase()),
		StateRoot:     common.Bytes32(executionGenesis.Root()),
		ReceiptsRoot:  common.Bytes32(executionGenesis.ReceiptHash()),
		LogsBloom:     common.LogsBloom(executionGenesis.Bloom()),
		PrevRandao:    common.Bytes32{},
		BlockNumber:   view.Uint64View(executionGenesis.NumberU64()),
		GasLimit:      view.Uint64View(executionGenesis.GasLimit()),
		GasUsed:       view.Uint64View(executionGenesis.GasUsed()),
		Timestamp:     common.Timestamp(executionGenesis.Time()),
		ExtraData:     extra,
		BaseFeePerGas: view.Uint256View(*baseFee),
		BlockHash:     common.Root(executionGenesis.Hash()),
		// empty transactions root
		TransactionsRoot: common.PayloadTransactionsType(g.Spec).
			DefaultNode().
			MerkleRoot(tree.GetHashFn()),
		// empty withdrawals root
		WithdrawalsRoot: common.WithdrawalsType(g.Spec).
			DefaultNode().
			MerkleRoot(tree.GetHashFn()),
	}

	return g.SetLatestExecutionPayloadHeader(execPayloadHeader)
}
