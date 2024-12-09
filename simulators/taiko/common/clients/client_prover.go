package clients

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"golang.org/x/net/context"
	"math/big"
	"sync"
	"sync/atomic"
)

type ProverClient struct {
	*HiveManagedClient
}

type State struct {
	*rpc.Client

	L1Head atomic.Pointer[types.Header]
	L2Head atomic.Pointer[types.Header]

	ProposedBlockID chan *big.Int

	LatestL1Origin    atomic.Pointer[rawdb.L1Origin]
	CanonicalL1Origin atomic.Pointer[rawdb.L1Origin]

	err atomic.Value

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewState(rpcCli *rpc.Client) (state *State, err error) {
	subCtx, cancel := context.WithCancel(context.Background())

	state = &State{
		Client:            rpcCli,
		L1Head:            atomic.Pointer[types.Header]{},
		L2Head:            atomic.Pointer[types.Header]{},
		ProposedBlockID:   make(chan *big.Int, 1),
		LatestL1Origin:    atomic.Pointer[rawdb.L1Origin]{},
		CanonicalL1Origin: atomic.Pointer[rawdb.L1Origin]{},

		ctx:    subCtx,
		cancel: cancel,
	}

	go state.loop()

	return state, nil
}

func (s *State) Close() {
	s.cancel()
	s.wg.Wait()
}

func (s *State) StateError() error {
	if val := s.err.Load(); val == nil {
		return nil
	} else {
		return val.(error)
	}
}

func (s *State) loop() {
	s.wg.Add(1)
	defer s.wg.Done()

	l1HeadCh := make(chan *types.Header, 3)
	l2HeadCh := make(chan *types.Header, 3)
	blockProposedCh := make(chan *bindings.TaikoL1ClientBlockProposed, 10)
	blockProposedV2Ch := make(chan *bindings.TaikoL1ClientBlockProposedV2, 10)

	l1Sub := rpc.SubscribeChainHead(s.L1, l1HeadCh)
	l2Sub := rpc.SubscribeChainHead(s.L2, l2HeadCh)
	l2BlockProposedSub := rpc.SubscribeBlockProposed(s.TaikoL1, blockProposedCh)
	l2BlockProposedV2Sub := rpc.SubscribeBlockProposedV2(s.TaikoL1, blockProposedV2Ch)

	defer func() {
		l1Sub.Unsubscribe()
		l2Sub.Unsubscribe()
		l2BlockProposedSub.Unsubscribe()
		l2BlockProposedV2Sub.Unsubscribe()
	}()

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case e := <-blockProposedCh:
			select {
			case <-s.ProposedBlockID:
				s.ProposedBlockID <- e.BlockId
			default:
				s.ProposedBlockID <- e.BlockId
			}
		case e := <-blockProposedV2Ch:
			select {
			case <-s.ProposedBlockID:
				s.ProposedBlockID <- e.BlockId
			default:
				s.ProposedBlockID <- e.BlockId
			}
		case head := <-l1HeadCh:
			s.L1Head.Store(head)
		case head := <-l2HeadCh:
			s.L2Head.Store(head)
			l1Origin, err := s.L2.L1OriginByID(ctx, head.Number)
			if err != nil {
				s.err.Store(err)
				continue
			}
			s.LatestL1Origin.Store(l1Origin)

			l1Origin, err = s.L2.HeadL1Origin(ctx)
			if err != nil {
				s.err.Store(err)
				continue
			}
			s.CanonicalL1Origin.Store(l1Origin)
		}
	}
}
