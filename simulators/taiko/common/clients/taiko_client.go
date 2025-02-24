package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/hive/hivesim"
	tkutils "github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/utils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/urfave/cli/v2"
	"sync"
	"sync/atomic"
)

type State struct {
	*rpc.Client

	L1Head atomic.Pointer[types.Header]
	L2Head atomic.Pointer[types.Header]

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewState(rpcCli *rpc.Client) (state *State, err error) {
	subCtx, cancel := context.WithCancel(context.Background())

	state = &State{
		Client: rpcCli,
		L1Head: atomic.Pointer[types.Header]{},
		L2Head: atomic.Pointer[types.Header]{},

		ctx:    subCtx,
		cancel: cancel,
	}

	go state.loop()

	return state, nil
}

func (s *State) Close() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
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
		}
	}
}

func NewTaikoClient[T tkutils.SubcommandApplication](client T, flags []cli.Flag) error {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:  "client",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return client.InitFromCli(context.Background(), c)
			},
		},
	}
	return app.Run([]string{"taiko-client", "client"})
}

func GetClientConfig(envs hivesim.Params) *rpc.ClientConfig {
	return &rpc.ClientConfig{
		L1Endpoint:                  envs["L1_WS"],
		L2Endpoint:                  envs["L2_WS"],
		TaikoL1Address:              common.HexToAddress(envs["TAIKO_INBOX"]),
		TaikoWrapperAddress:         common.HexToAddress(envs["TAIKO_WRAPPER"]),
		ForcedInclusionStoreAddress: common.HexToAddress(envs["FORCED_INCLUSION_STORE"]),
		ProverSetAddress:            common.HexToAddress(envs["PROVER_SET"]),
		TaikoL2Address:              common.HexToAddress(envs["TAIKO_ANCHOR"]),
		TaikoTokenAddress:           common.HexToAddress(envs["TAIKO_TOKEN"]),
		L2EngineEndpoint:            envs["L2_AUTH"],
		JwtSecret:                   envs["JWT_SECRET"],
	}
}
