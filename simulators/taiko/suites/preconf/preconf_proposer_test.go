package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"taiko/bindings/ontake"
	"taiko/bindings/pacaya"
	"taiko/common/clients"
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
		L1Endpoint:        "ws://localhost:8545",
		L2Endpoint:        "ws://localhost:6046",
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:  "http://localhost:6051",
		JwtSecret:         "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
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
	err := preconferProposer(rpccli, preconfURL)
	assert.NoError(t, err)
}
