package clients

import "github.com/ethereum/go-ethereum/core/rawdb"

var (
	AnvilPort int64 = 8545
)

type L1BlockInfo struct {
	L2Number uint64
	Snapshot string
}

type BlockProposed struct {
	BlockId uint64
	MinTier uint16
}

type AnvilClient struct {
	SecondsPerSlot uint64
	*EthNode

	reorgCh     chan struct{}
	reorgPoints map[uint64]string
	l1Origins   map[uint64]*rawdb.L1Origin
}

func (a *AnvilClient) Start() error {
	return a.EthNode.Start()
}

func (a *AnvilClient) Shutdown() error {
	a.StopRecordReorgPoints()
	return a.EthNode.Shutdown()
}
