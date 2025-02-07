package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"os"
	"testing"
)

func TestCC(t *testing.T) {
	rpcCli, err := rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:        "ws://localhost:8545",
		L2Endpoint:        "ws://localhost:6046",
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:  "http://localhost:6051",
		JwtSecret:         "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7",
	})
	assert.NoError(t, err)

	code, err := rpcCli.L2.CodeAt(nil, common.HexToAddress(os.Getenv("TAIKO_INBOX")), nil)
	assert.NoError(t, err)
	t.Log(len(code))

	res, err := rpcCli.PacayaClients.TaikoInbox.GetStats2(nil)
	assert.NoError(t, err)
	t.Log(res.Paused)
}
