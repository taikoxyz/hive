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
	a.reorgCache = make(map[uint64]*L1BlockInfo)
	a.reorgCh = make(chan struct{})
	return a.EthNode.Start()
}

func (a *AnvilClient) Shutdown() error {
	if a.reorgCh != nil {
		close(a.reorgCh)
	}
	return a.EthNode.Shutdown()
}
