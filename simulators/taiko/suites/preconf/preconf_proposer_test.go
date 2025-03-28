package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"taiko/common/clients"
	"taiko/params"
	"testing"
)

var (
	rpccli       *rpc.Client
	l1Cli, l2Cli *rpc.EthClient
	preconfURL   = fmt.Sprintf("http://localhost:%d", clients.PreconfServerPort)
	l1Ontake     *ontake.OntakeL1Clients
	l1Pacaya     *pacaya.PacayaL1Clients
)

func init() {
	var err error
	rpccli, err = rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:                  "ws://localhost:8545",
		L2Endpoint:                  "ws://localhost:6046",
		TaikoL1Address:              common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoWrapperAddress:         common.HexToAddress(os.Getenv("TAIKO_WRAPPER")),
		ForcedInclusionStoreAddress: common.HexToAddress(os.Getenv("FORCED_INCLUSION_STORE")),
		ProverSetAddress:            common.HexToAddress(os.Getenv("PROVER_SET")),
		TaikoL2Address:              common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress:           common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:            "http://localhost:6051",
		JwtSecret:                   "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
	})
	if err != nil {
		panic(err)
	}

	l1Cli, l2Cli = rpccli.L1, rpccli.L2
	l1Ontake, err = ontake.NewOntakeL1Clients(l1Cli)
	if err != nil {
		panic(err)
	}

	l1Pacaya, err = pacaya.NewPacayaL1Clients(l2Cli)
	if err != nil {
		panic(err)
	}
}

func TestPreconferProposer(t *testing.T) {
	head, err := rpccli.PacayaClients.ForcedInclusionStore.Head(nil)
	assert.NoError(t, err)
	t.Log(head)

	envs := params.EnvParams()
	envs["L1_WS"] = "ws://localhost:8545"

	_, err = storeForcedInclusion(envs, rpccli)
	assert.NoError(t, err)

	head, err = rpccli.PacayaClients.ForcedInclusionStore.Head(nil)
	assert.NoError(t, err)
	t.Log(head)
}

func TestCC(t *testing.T) {
	head, err := rpccli.PacayaClients.ForcedInclusionStore.Head(nil)
	assert.NoError(t, err)
	t.Log(head)
}

func TestDD(t *testing.T) {
	enr := "enr:-Ja4QCZM8pihrIWZOr9A647xApl6fKJkJvTlMZwca9pJE10IK0BIKTvNc8kch2GKLc1nv0uAPDigCrImFe6yqgRC7oyGAZU67I-KgmlkgnY0h29wc3RhY2uE2ZgKAIlzZWNwMjU2azGhArpXNNj3CRcZRx5_fta53xcNxwzGYcoF5ohgGtmE8Giwg3RjcIIkBoN1ZHCCdmE"
	node, err := enode.Parse(enode.ValidSchemes, enr)
	assert.NoError(t, err)
	t.Log(node.URLv4())
}

func TestPreconferBlock(t *testing.T) {
	_, _, _, err := preconferBlock(0, rpccli, preconfURL, 5)
	assert.NoError(t, err)
}
