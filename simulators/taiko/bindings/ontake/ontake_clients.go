package ontake

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"taiko/bindings/ontake/forkrouter"
	"taiko/bindings/ontake/guardianprover"
	"taiko/bindings/ontake/libproposing"
	"taiko/bindings/ontake/proverset"
	"taiko/bindings/ontake/taikol1"
	"taiko/bindings/ontake/taikol2"
	"taiko/bindings/ontake/taikotoken"
	"taiko/params"
)

type OntakeL1Clients struct {
	TaikoL1                *taikol1.TaikoL1
	LibProposing           *libproposing.LibProposing
	TaikoToken             *taikotoken.TaikoToken
	GuardianProverMajority *guardianprover.GuardianProver
	GuardianProverMinority *guardianprover.GuardianProver
	ProverSet              *proverset.ProverSet
	ForkRouter             *forkrouter.ForkRouter
}

func NewOntakeL1Clients(l1Cli *rpc.EthClient) (*OntakeL1Clients, error) {
	taikoL1, err := taikol1.NewTaikoL1(params.ParamToAddress("TAIKO_INBOX"), l1Cli)
	if err != nil {
		return nil, err
	}

	taikoToken, err := taikotoken.NewTaikoToken(params.ParamToAddress("TAIKO_TOKEN"), l1Cli)
	if err != nil {
		return nil, err
	}

	libProposing, err := libproposing.NewLibProposing(params.ParamToAddress("TAIKO_INBOX"), l1Cli)
	if err != nil {
		return nil, err
	}

	guardianProverMajority, err := guardianprover.NewGuardianProver(params.ParamToAddress("GUARDIAN_PROVER_MAJORITY"), l1Cli)
	if err != nil {
		return nil, err
	}

	guardianProverMinority, err := guardianprover.NewGuardianProver(params.ParamToAddress("GUARDIAN_PROVER_MINORITY"), l1Cli)
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

	forkRouter, err := forkrouter.NewForkRouter(params.ParamToAddress("TAIKO_INBOX"), l1Cli)
	if err != nil {
		return nil, err
	}

	return &OntakeL1Clients{
		TaikoL1:                taikoL1,
		TaikoToken:             taikoToken,
		LibProposing:           libProposing,
		GuardianProverMajority: guardianProverMajority,
		GuardianProverMinority: guardianProverMinority,
		ProverSet:              proverSet,
		ForkRouter:             forkRouter,
	}, nil
}

type OntakeL2Clients struct {
	TaikoL2 *taikol2.TaikoL2
}

func NewOntakeL2Clients(l2Cli *rpc.EthClient) (*OntakeL2Clients, error) {
	taikoL2, err := taikol2.NewTaikoL2(params.ParamToAddress("TAIKO_ANCHOR"), l2Cli)
	if err != nil {
		return nil, err
	}

	return &OntakeL2Clients{
		TaikoL2: taikoL2,
	}, nil
}

// OntakeClients contains all smart contract clients for Ontake fork.
type OntakeClients struct {
	*OntakeL1Clients
	*OntakeL2Clients
}

func NewOntakeClients(l1cli, l2cli *rpc.EthClient) (*OntakeClients, error) {
	l1Clients, err := NewOntakeL1Clients(l1cli)
	if err != nil {
		return nil, err
	}

	l2Clients, err := NewOntakeL2Clients(l2cli)
	if err != nil {
		return nil, err
	}

	return &OntakeClients{
		OntakeL1Clients: l1Clients,
		OntakeL2Clients: l2Clients,
	}, nil
}
