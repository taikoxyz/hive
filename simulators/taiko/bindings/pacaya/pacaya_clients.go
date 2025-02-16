package pacaya

import (
	"github.com/ethereum/go-ethereum/common"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"taiko/bindings/pacaya/forkrouter"
	"taiko/bindings/pacaya/proverset"
	"taiko/bindings/pacaya/taikoanchor"
	"taiko/bindings/pacaya/taikoinbox"
	"taiko/bindings/pacaya/taikotoken"
	"taiko/params"
)

type PacayaL1Clients struct {
	TaikoInbox *taikoinbox.TaikoInbox
	TaikoToken *taikotoken.TaikoToken
	ProverSet  *proverset.ProverSet
	ForkRouter *forkrouter.ForkRouter
}

func NewPacayaL1Clients(l1Cli *rpc.EthClient) (*PacayaL1Clients, error) {
	taikoInbox, err := taikoinbox.NewTaikoInbox(common.HexToAddress(params.ParamByKey("TAIKO_INBOX")), l1Cli)
	if err != nil {
		return nil, err
	}

	taikoToken, err := taikotoken.NewTaikoToken(common.HexToAddress(params.ParamByKey("TAIKO_TOKEN")), l1Cli)
	if err != nil {
		return nil, err
	}

	var proverSet *proverset.ProverSet
	if params.ParamToAddress("PROVER_SET") != (common.Address{}) {
		proverSet, err = proverset.NewProverSet(params.ParamToAddress("PROVER_SET"), l1Cli)
		if err != nil {
			return nil, err
		}
	}

	forkRouter, err := forkrouter.NewForkRouter(common.HexToAddress(params.ParamByKey("TAIKO_INBOX")), l1Cli)
	if err != nil {
		return nil, err
	}

	return &PacayaL1Clients{
		TaikoInbox: taikoInbox,
		TaikoToken: taikoToken,
		ProverSet:  proverSet,
		ForkRouter: forkRouter,
	}, nil
}

type PacayaL2Clients struct {
	TaikoAnchor *taikoanchor.TaikoAnchor
	ForkHeight  uint64
}

func NewPacayaL2Clients(l2cli *rpc.EthClient) (*PacayaL2Clients, error) {
	taikoAnchor, err := taikoanchor.NewTaikoAnchor(common.HexToAddress(params.ParamByKey("TAIKO_ANCHOR")), l2cli)
	if err != nil {
		return nil, err
	}

	forkHeight := uint64(0)
	switch l2cli.ChainID.Uint64() {
	case ethparams.HeklaNetworkID.Uint64(),
		ethparams.TaikoMainnetNetworkID.Uint64(),
		ethparams.PreconfDevnetNetworkID.Uint64():
	default:
		forkHeight = 10
	}

	return &PacayaL2Clients{
		TaikoAnchor: taikoAnchor,
		ForkHeight:  forkHeight,
	}, nil
}

// PacayaClients contains all smart contract clients for Pacaya fork.
type PacayaClients struct {
	*PacayaL1Clients
	*PacayaL2Clients
}

func NewPacayaClients(l1cli, l2cli *rpc.EthClient) (*PacayaClients, error) {
	l1Clients, err := NewPacayaL1Clients(l1cli)
	if err != nil {
		return nil, err
	}

	l2Clients, err := NewPacayaL2Clients(l2cli)
	if err != nil {
		return nil, err
	}

	return &PacayaClients{
		PacayaL1Clients: l1Clients,
		PacayaL2Clients: l2Clients,
	}, nil
}
