package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"net"
	"strings"
	"taiko/common/utils"
)

type ExecutionProxyConfig struct {
	Host net.IP
	Port int
}

type ExecutionClientConfig struct {
	ClientIndex int
	ProxyConfig *ExecutionProxyConfig
	Subnet      string
	JWTSecret   [32]byte
	Network     string
}

type ExecutionClient struct {
	*EthNode
}

func (ec *ExecutionClient) EngineURL() string {
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), EthEngineRPC)
}

func (ec *ExecutionClient) HTTPClient() *rpc.EthClient {
	return ec.EthClient
}

func (ec *ExecutionClient) HeaderByHash(
	parentCtx context.Context,
	h common.Hash,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.EthClient.HeaderByHash(ctx, h)
}

func (ec *ExecutionClient) HeaderByNumber(
	parentCtx context.Context,
	n *big.Int,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.EthClient.HeaderByNumber(ctx, n)
}

func (ec *ExecutionClient) HeaderByLabel(
	parentCtx context.Context,
	l string,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	h := new(types.Header)
	client := ec.EthClient
	err := client.CallContext(
		ctx,
		h,
		"eth_getBlockByNumber",
		l,
		false,
	)
	return h, err
}

func (ec *ExecutionClient) BlockByHash(
	parentCtx context.Context,
	h common.Hash,
) (*types.Block, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.EthClient.BlockByHash(ctx, h)
}

func (ec *ExecutionClient) BlockByNumber(
	parentCtx context.Context,
	n *big.Int,
) (*types.Block, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.EthClient.BlockByNumber(ctx, n)
}

type BinaryMarshable interface {
	MarshalBinary() ([]byte, error)
}

func (ec *ExecutionClient) SendTransaction(
	parentCtx context.Context,
	tx BinaryMarshable,
) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()

	client := ec.EthClient
	return client.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(data))
}

type ExecutionClients []*ExecutionClient

// Return subset of clients that are currently running
func (all ExecutionClients) Running() ExecutionClients {
	res := make(ExecutionClients, 0)
	for _, ec := range all {
		if ec.IsRunning() {
			res = append(res, ec)
		}
	}
	return res
}

// Returns comma-separated Bootnodes of all running execution nodes
func (all ExecutionClients) Enodes() (string, error) {
	if len(all) == 0 {
		return "", nil
	}
	enodes := make([]string, 0)
	for _, en := range all {
		if en.IsRunning() {
			enode, err := en.GetEnodeURL()
			if err != nil {
				return "", err
			}
			enodes = append(enodes, enode)
		}
	}
	return strings.Join(enodes, ","), nil
}

// Returns true if all head hashes match
func (all ExecutionClients) CheckHeads(
	l utils.Logging,
	parentCtx context.Context,
) (bool, error) {
	if len(all) <= 1 {
		return false, fmt.Errorf(
			"attempted to check the heads of a single or zero clients matched",
		)
	}

	header, err := all[0].HeaderByNumber(parentCtx, nil)
	if err != nil || header == nil {
		return false, err
	}
	baseHash := header.Hash()

	for _, en := range all[1:] {
		header, err = en.HeaderByNumber(parentCtx, nil)
		if err != nil || header == nil {
			return false, err
		}
		h := header.Hash()
		if h != baseHash {
			if l != nil {
				l.Logf("Hash mismatch between heads: %s != %s\n", h, baseHash)
			}
			return false, nil
		} else if l != nil {
			l.Logf("Hash match between heads: %s == %s\n", h, baseHash)
		}
	}
	return true, nil
}
