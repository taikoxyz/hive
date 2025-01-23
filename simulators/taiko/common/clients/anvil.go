package clients

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

	reorgCh    chan struct{}
	reorgCache map[uint64]*L1BlockInfo
}

func (a *AnvilClient) Start() error {
	if err := a.EthNode.Start(); err != nil {
		return err
	}
	return a.EthNode.InitL1Clients()
}
