package clients

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"golang.org/x/net/context"
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

	l1Sub := rpc.SubscribeChainHead(s.L1, l1HeadCh)
	l2Sub := rpc.SubscribeChainHead(s.L2, l2HeadCh)
	defer func() {
		l1Sub.Unsubscribe()
		l2Sub.Unsubscribe()
	}()

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
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
