package pacaya

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func NewPacayaL1Clients(l1Cli *ethclient.Client) (*PacayaL1Clients, error) {
	taikoInbox, err := taikoinbox.NewTaikoInbox(common.HexToAddress(params.ParamByKey("TAIKO_INBOX")), l1Cli)
	if err != nil {
		return nil, err
	}

	taikoToken, err := taikotoken.NewTaikoToken(common.HexToAddress(params.ParamByKey("TAIKO_TOKEN")), l1Cli)
	if err != nil {
		return nil, err
	}

	proverSet, err := proverset.NewProverSet(common.HexToAddress(params.ParamByKey("PROVER_SET")), l1Cli)
	if err != nil {
		return nil, err
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
}

func NewPacayaL2Clients(l2cli *ethclient.Client) (*PacayaL2Clients, error) {
	taikoAnchor, err := taikoanchor.NewTaikoAnchor(common.HexToAddress(params.ParamByKey("TAIKO_ANCHOR")), l2cli)
	if err != nil {
		return nil, err
	}

	return &PacayaL2Clients{
		TaikoAnchor: taikoAnchor,
	}, nil
}

// PacayaClients contains all smart contract clients for Pacaya fork.
type PacayaClients struct {
	*PacayaL1Clients
	*PacayaL2Clients
	ForkHeight uint64
}

func NewPacayaClients(l1cli, l2cli *ethclient.Client) (*PacayaClients, error) {
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

func (p *PacayaClients) SetForkHeight(height uint64) {
	p.ForkHeight = height
}
