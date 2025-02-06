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
	_, err := rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:        "ws://localhost:8545",
		L2Endpoint:        "ws://localhost:6046",
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L2EngineEndpoint:  "http://localhost:6051",
		JwtSecret:         os.Getenv("JWT_SECRET"),
	})
	assert.NoError(t, err)
}
