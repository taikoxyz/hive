package clients

import (
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"taiko/bindings/devnettierprovider"
	"taiko/bindings/taikol1"
)

var (
	errEmptyTiersList = errors.New("empty proof tiers list in protocol")
)

var (
	AnvilPort int64 = 8545
)

type L1BlockInfo struct {
	Header   *types.Header
	Snapshot string
}

type BlockProposed struct {
	BlockId uint64
	MinTier uint16
}

type AnvilClient struct {
	SecondsPerSlot uint64
	*EthNode

	reorgCh    chan struct{}
	reorgCache map[uint64]*L1BlockInfo

	taikoL1 *taikol1.TaikoL1

	tiers map[uint16]*devnettierprovider.ITierProviderTier

	proposedEvents []*BlockProposed
}
