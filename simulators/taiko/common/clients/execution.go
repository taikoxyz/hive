package clients

import (
	"context"
	"encoding/json"
	"fmt"
	api "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/marioevz/eth-clients/clients/execution"
	"math/big"
	"net"
	"strings"
	"sync"
	"taiko/common/utils"
	"time"
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
	Config ExecutionClientConfig

	proxy     *execution.Proxy
	latestfcu *api.ForkchoiceStateV1

	engineClient *rpc.Client
	httpClient   *ethclient.Client

	startupComplete bool
}

func (ec *ExecutionClient) EngineURL() string {
	return fmt.Sprintf("http://%v:%d", ec.NetworkIP(), EthEngineRPC)
}

func (ec *ExecutionClient) GetLatestForkchoiceUpdated(
	ctx context.Context,
) (*api.ForkchoiceStateV1, error) {
	if ec.latestfcu != nil {
		return ec.latestfcu, nil
	}
	// Try to reconstruct by querying it from the client
	forkchoiceState := &api.ForkchoiceStateV1{}
	errs := make(chan error, 3)
	var wg sync.WaitGroup

	type labelBlockHashTask struct {
		label string
		dest  *common.Hash
	}

	for _, t := range []*labelBlockHashTask{
		{
			label: "latest",
			dest:  &forkchoiceState.HeadBlockHash,
		},
		{
			label: "safe",
			dest:  &forkchoiceState.SafeBlockHash,
		},
		{
			label: "finalized",
			dest:  &forkchoiceState.FinalizedBlockHash,
		},
	} {
		wg.Add(1)
		t := t
		go func(t *labelBlockHashTask) {
			defer wg.Done()
			if res, err := ec.HeaderByLabel(
				ctx,
				t.label,
			); err != nil {
				ec.Logf(
					"Error trying to fetch label %s from client: %v",
					t.label,
					err,
				)
			} else if err == nil && res != nil && res.Number != nil {
				*t.dest = res.Hash()
			}
		}(t)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return nil, err
	default:
	}

	return forkchoiceState, nil
}

func (ec *ExecutionClient) EngineForkchoiceUpdated(
	parentCtx context.Context,
	fcState *api.ForkchoiceStateV1,
	pAttributes *api.PayloadAttributes,
	version int,
) (*api.ForkChoiceResponse, error) {
	var result api.ForkChoiceResponse
	request := fmt.Sprintf("engine_forkchoiceUpdatedV%d", version)
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	err := ec.engineClient.CallContext(
		ctx,
		&result,
		request,
		fcState,
		pAttributes,
	)
	return &result, err
}

func (ec *ExecutionClient) EngineGetPayload(
	parentCtx context.Context,
	payloadID *api.PayloadID,
	version int,
) (*api.ExecutableData, *big.Int, *api.BlobsBundleV1, *bool, error) {
	var (
		rpcString = fmt.Sprintf("engine_getPayloadV%d", version)
	)
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	if version >= 2 {
		var response api.ExecutionPayloadEnvelope
		err := ec.engineClient.CallContext(
			ctx,
			&response,
			rpcString,
			payloadID,
		)
		return response.ExecutionPayload, response.BlockValue, response.BlobsBundle, &response.Override, err
	} else {
		var executableData api.ExecutableData
		err := ec.engineClient.CallContext(ctx, &executableData, rpcString, payloadID)
		return &executableData, common.Big0, nil, nil, err
	}
}

func (ec *ExecutionClient) EngineNewPayload(
	parentCtx context.Context,
	payload *api.ExecutableData,
	version int,
) (*api.PayloadStatusV1, error) {
	var result api.PayloadStatusV1
	request := fmt.Sprintf("engine_newPayloadV%d", version)
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	err := ec.engineClient.CallContext(ctx, &result, request, payload)
	return &result, err
}

// Eth RPC
// Helper structs to fetch the TotalDifficulty
type TD struct {
	TotalDifficulty *hexutil.Big `json:"totalDifficulty"`
}
type TotalDifficultyHeader struct {
	types.Header
	TD
}

func (tdh *TotalDifficultyHeader) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &tdh.Header); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &tdh.TD); err != nil {
		return err
	}
	return nil
}

func (ec *ExecutionClient) HeaderByHash(
	parentCtx context.Context,
	h common.Hash,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.httpClient.HeaderByHash(ctx, h)
}

func (ec *ExecutionClient) HeaderByNumber(
	parentCtx context.Context,
	n *big.Int,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.httpClient.HeaderByNumber(ctx, n)
}

func (ec *ExecutionClient) HeaderByLabel(
	parentCtx context.Context,
	l string,
) (*types.Header, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	h := new(types.Header)
	client := ec.httpClient.Client()
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
	return ec.httpClient.BlockByHash(ctx, h)
}

func (ec *ExecutionClient) BlockByNumber(
	parentCtx context.Context,
	n *big.Int,
) (*types.Block, error) {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return ec.httpClient.BlockByNumber(ctx, n)
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

	client := ec.httpClient.Client()
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

// Return subset of clients that are part of an specific subnet
func (all ExecutionClients) Subnet(subnet string) ExecutionClients {
	if subnet == "" {
		return all
	}
	res := make(ExecutionClients, 0)
	for _, ec := range all {
		if ec.Config.Subnet == subnet {
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

// Interface used to provide a proxy for a given execution client
type ProxyProvider interface {
	Proxy() *execution.Proxy
}
type Proxies []ProxyProvider

func (all Proxies) Running() []*execution.Proxy {
	res := make([]*execution.Proxy, 0)
	for _, pp := range all {
		if p := pp.Proxy(); p != nil {
			res = append(res, p)
		}
	}
	return res
}
