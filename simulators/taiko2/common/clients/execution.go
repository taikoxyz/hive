package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	api "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/hive/taiko"
	"github.com/ethereum/hive/taiko/params"
	"github.com/marioevz/eth-clients/clients"
	"github.com/marioevz/eth-clients/clients/execution"
	spoof "github.com/rauljordan/engine-proxy/proxy"
	"math/big"
	"net"
	"strings"
	"sync"
	"taiko2/common/utils"
	"time"
)

const (
	PortHttpRPC   = 8545
	PortWSRPC     = 8546
	PortEngineRPC = 8551
)

var AllForkchoiceUpdatedCalls = []string{
	"engine_forkchoiceUpdatedV1",
	"engine_forkchoiceUpdatedV2",
	"engine_forkchoiceUpdatedV3",
}

var AllGetPayloadCalls = []string{
	"engine_getPayloadV1",
	"engine_getPayloadV2",
	"engine_getPayloadV3",
}

var AllNewPayloadCalls = []string{
	"engine_newPayloadV1",
	"engine_newPayloadV2",
	"engine_newPayloadV3",
}

var AllEngineCalls = []string{
	"engine_forkchoiceUpdatedV1",
	"engine_forkchoiceUpdatedV2",
	"engine_forkchoiceUpdatedV3",
	"engine_getPayloadV1",
	"engine_getPayloadV2",
	"engine_getPayloadV3",
	"engine_newPayloadV1",
	"engine_newPayloadV2",
	"engine_newPayloadV3",
}

type EnodeClient interface {
	clients.Client
	GetEnodeURL() (string, error)
}

type ExecutionProxyConfig struct {
	Host                   net.IP
	Port                   int
	LogEngineCalls         bool
	TrackForkchoiceUpdated bool
}

type ExecutionClientConfig struct {
	ClientIndex             int
	ProxyConfig             *ExecutionProxyConfig
	TerminalTotalDifficulty int64
	Subnet                  string
	JWTSecret               [32]byte
}

type ExecutionClient struct {
	Client
	Logger utils.Logging
	Config ExecutionClientConfig

	proxy     *execution.Proxy
	latestfcu *api.ForkchoiceStateV1

	engineClient *rpc.Client
	httpClient   *ethclient.Client

	startupComplete bool

	Deploy bool
}

func (ec *ExecutionClient) Logf(format string, values ...interface{}) {
	if l := ec.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (ec *ExecutionClient) HttpURL() string {
	return fmt.Sprintf("http://%v:%d", ec.GetHost(), PortHttpRPC)
}

func (ec *ExecutionClient) WSURL() string {
	return fmt.Sprintf("ws://%v:%d", ec.GetHost(), PortWSRPC)
}

func (ec *ExecutionClient) EngineURL() string {
	return fmt.Sprintf("http://%v:%d", ec.GetHost(), PortEngineRPC)
}

func (ec *ExecutionClient) MustGetEnode() string {
	if enodeClient, ok := ec.Client.(EnodeClient); ok {
		addr, err := enodeClient.GetEnodeURL()
		if err == nil {
			return addr
		}
		panic(err)
	}
	panic(fmt.Errorf("invalid client type"))
}

func (ec *ExecutionClient) ConfiguredTTD() *big.Int {
	return big.NewInt(ec.Config.TerminalTotalDifficulty)
}

func (ec *ExecutionClient) Start() error {
	if !ec.Client.IsRunning() {
		if managedClient, ok := ec.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			if err := managedClient.Start(); err != nil {
				return err
			}
		}
	}

	return ec.Init(context.Background())
}

func (ec *ExecutionClient) Init(ctx context.Context) (err error) {
	if !ec.Client.IsRunning() {
		return fmt.Errorf("execution client not yet launched")
	}
	if !ec.startupComplete {
		defer func() {
			ec.startupComplete = true
		}()

		// Prepare HTTP Client
		ec.engineClient, err = rpc.DialOptions(context.Background(), ec.EngineURL(), rpc.WithHTTPAuth(node.NewJWTAuth(ec.Config.JWTSecret)))
		if err != nil {
			return err
		}

		ec.httpClient, err = ethclient.DialContext(ctx, ec.HttpURL())
		if err != nil {
			return err
		}

		// Prepare proxy
		dest := ec.EngineURL()

		if ec.Config.ProxyConfig != nil {
			p := execution.NewProxy(
				ec.Config.ProxyConfig.Host,
				ec.Config.ProxyConfig.Port,
				dest,
				ec.Config.JWTSecret[:],
			)

			if ec.Config.ProxyConfig.TrackForkchoiceUpdated {
				logCallback := func(req []byte) *spoof.Spoof {
					var (
						fcState api.ForkchoiceStateV1
						pAttr   api.PayloadAttributes
						err     error
					)
					err = execution.UnmarshalFromJsonRPCRequest(
						req,
						&fcState,
						&pAttr,
					)
					if err == nil {
						ec.latestfcu = &fcState
					} else {
						ec.Logf(
							"Error trying to unmarshal forkchoice state: %v. Latest FCU will be nil",
							err,
						)
						ec.latestfcu = nil
					}
					return nil
				}
				p.AddRequestCallbacks(logCallback, AllForkchoiceUpdatedCalls...)
			}

			if ec.Config.ProxyConfig.LogEngineCalls {
				logCallback := func(res []byte, req []byte) *spoof.Spoof {
					ec.Logf(
						"DEBUG: execution client %d, request: %s",
						ec.Config.ClientIndex,
						req,
					)
					ec.Logf(
						"DEBUG: execution client %d, response: %s",
						ec.Config.ClientIndex,
						res,
					)
					return nil
				}
				p.AddResponseCallbacks(logCallback, AllEngineCalls...)
			}

			ec.proxy = p
		}
	}

	var (
		client  *ethclient.Client
		chainID *big.Int
	)
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for chainID == nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 10):
			chainID = big.NewInt(0)
			break
		case <-tick.C:
			ec.Logger.Logf("Waiting for chainID")
			client, err = ethclient.DialContext(ctx, ec.HttpURL())
			if err != nil {
				continue
			}
			chainID, _ = client.ChainID(ctx)
		}
	}

	if ec.Deploy {
		return ec.DeployContracts(ctx, chainID, client)
	}

	return nil
}

func (ec *ExecutionClient) DeployContracts(ctx context.Context, chainID *big.Int, client *ethclient.Client) error {
	sk, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	if err != nil {
		return err
	}

	for _, tx := range params.ContractTxs {
		data, _ := json.Marshal(tx)
		ec.Logf("Deploying contract: %v", string(data))
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			return err
		}

		if err := client.SendTransaction(ctx, signedTx); err != nil {
			return err
		}
	}

	// check results.
	for _, tx := range params.ContractTxs {
		if tx.To() == nil {
			_, err := bind.WaitDeployed(context.Background(), client, tx)
			if err != nil {
				return fmt.Errorf("failed to wait deployed: %v", err)
			}
		} else {
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				return fmt.Errorf("failed to wait mined, hash: %s, err: %v", tx.Hash().String(), err)
			}
			if receipt.Status != types.ReceiptStatusSuccessful {
				return fmt.Errorf("failed to call contract, hash: %s", tx.Hash().String())
			}
		}
	}

	envs := params.EnvParams.Copy()
	envs["L1_HTTP"] = ec.HttpURL()

	// init contracts.
	return taiko.InitTaikoContract(envs)
}

func (ec *ExecutionClient) Shutdown() error {
	if managedClient, ok := ec.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}

func (ec *ExecutionClient) Proxy() *execution.Proxy {
	return ec.proxy
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
			if enodeClient, ok := en.Client.(EnodeClient); ok {
				enode, err := enodeClient.GetEnodeURL()
				if err != nil {
					return "", err
				}
				enodes = append(enodes, enode)
			} else {
				return "", fmt.Errorf("invalid client type")
			}
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
