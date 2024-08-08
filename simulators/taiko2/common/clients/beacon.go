package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/client/beaconapi"
	"github.com/protolambda/eth2api/client/nodeapi"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/ztyp/tree"
	"net/http"
	"strings"
	"sync"
	"taiko2/common/config/consensus"
	"taiko2/common/utils"
	"time"
)

const PortBeaconTCP = 9000

type BeaconClientConfig struct {
	ClientIndex             int
	TerminalTotalDifficulty int64
	Spec                    *consensus_config.Spec
	GenesisValidatorsRoot   *tree.Root
	GenesisTime             *common.Timestamp
	Subnet                  string
	Network                 string
}

type BeaconClient struct {
	*HiveManagedClient
	Logger    utils.Logging
	Config    *BeaconClientConfig
	networkIP string

	api *eth2api.Eth2HttpClient
}

func (bn *BeaconClient) Logf(format string, values ...interface{}) {
	if l := bn.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (bn *BeaconClient) BeaconURL() string {
	return fmt.Sprintf("http://%s:%d", bn.NetworkIP(), BeaconPort)
}

func (bn *BeaconClient) NetworkIP() string {
	if bn.networkIP != "" {
		return bn.networkIP
	}

	t := bn.T
	var err error
	bn.networkIP, err = t.Sim.ContainerNetworkIP(t.SuiteID, bn.Config.Network, bn.Client.Container)
	if err != nil {
		bn.Logf("Failed to get network IP: %v", err)
		return bn.HiveManagedClient.GetHost()
	}
	return bn.networkIP
}

func (bn *BeaconClient) Start() error {
	bn.Logf("Starting beacon client %d", bn.Config.ClientIndex)
	if !bn.IsRunning() {
		if err := bn.HiveManagedClient.Start(); err != nil {
			return err
		}
	}

	return bn.Init(context.Background())
}

func (bn *BeaconClient) Init(ctx context.Context) error {
	if bn.api == nil {
		bn.api = &eth2api.Eth2HttpClient{
			Addr:  bn.GetAddress(),
			Cli:   &http.Client{},
			Codec: eth2api.JSONCodec{},
		}
	}

	var wg sync.WaitGroup
	var errs = make(chan error, 2)

	if bn.Config.GenesisTime == nil || bn.Config.GenesisValidatorsRoot == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if gen, err := bn.GenesisConfig(ctx); err == nil &&
					gen != nil {
					bn.Config.GenesisTime = &gen.GenesisTime
					bn.Config.GenesisValidatorsRoot = &gen.GenesisValidatorsRoot
					return
				}
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-time.After(time.Second):
				}
			}
		}()
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func (bn *BeaconClient) Shutdown() error {
	return bn.HiveManagedClient.Shutdown()
}

func (bn *BeaconClient) ENR(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	var out eth2api.NetworkIdentity
	if err := nodeapi.Identity(ctx, bn.api, &out); err != nil {
		return "", err
	}
	return out.ENR, nil
}

func (bn *BeaconClient) P2PAddr(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	var out eth2api.NetworkIdentity
	if err := nodeapi.Identity(ctx, bn.api, &out); err != nil {
		return "", err
	}
	ip := bn.GetIP()
	if ip != nil {
		return fmt.Sprintf(
			"/ip4/%s/tcp/%d/p2p/%s",
			ip.String(),
			PortBeaconTCP,
			out.PeerID,
		), nil
	} else {
		return fmt.Sprintf(
			"/dns/%s/tcp/%d/p2p/%s",
			bn.GetHost(),
			PortBeaconTCP,
			out.PeerID,
		), nil
	}
}

func (bn *BeaconClient) BeaconAPIURL() (string, error) {
	if bn.api == nil {
		return "", fmt.Errorf("api not initialized")
	}
	return bn.api.Addr, nil
}

func (bn *BeaconClient) EnodeURL() (string, error) {
	return "", errors.New(
		"beacon node does not have an discv4 Enode URL, use ENR or multi-address instead",
	)
}

func (bn *BeaconClient) ClientName() string {
	name := bn.ClientType()
	if len(name) > 3 && name[len(name)-3:] == "-bn" {
		name = name[:len(name)-3]
	}
	return name
}

func (bn *BeaconClient) API() *eth2api.Eth2HttpClient {
	return bn.api
}

func (bn *BeaconClient) GenesisConfig(
	parentCtx context.Context,
) (*eth2api.GenesisResponse, error) {
	var (
		dest   = new(eth2api.GenesisResponse)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()

	exists, err = beaconapi.Genesis(ctx, bn.api, dest)

	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return dest, err
}

func (bn *BeaconClient) BlockV2Root(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (tree.Root, error) {
	var (
		root   tree.Root
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	root, exists, err = beaconapi.BlockRoot(ctx, bn.api, blockId)
	if !exists {
		return root, fmt.Errorf(
			"endpoint not found on beacon client",
		)
	}
	return root, err
}

type BlockV2OptimisticResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}

func (bn *BeaconClient) BlockIsOptimistic(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (bool, error) {
	var (
		blockOptResp = new(BlockV2OptimisticResponse)
		exists       bool
		err          error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = eth2api.SimpleRequest(
		ctx,
		bn.api,
		eth2api.FmtGET("/httpClient/v2/beacon/blocks/%s", blockId.BlockId()),
		blockOptResp,
	)
	if !exists {
		return false, fmt.Errorf("endpoint not found on beacon client")
	}
	return blockOptResp.ExecutionOptimistic, err
}

func (bn *BeaconClient) BlockHeader(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (*eth2api.BeaconBlockHeaderAndInfo, error) {
	var (
		headInfo = new(eth2api.BeaconBlockHeaderAndInfo)
		exists   bool
		err      error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.BlockHeader(ctx, bn.api, blockId, headInfo)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return headInfo, err
}

func (bn *BeaconClient) BlobSidecars(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) ([]deneb.BlobSidecar, error) {
	var (
		blobSidecars = new([]deneb.BlobSidecar)
		exists       bool
		err          error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.BlobSidecars(ctx, bn.api, blockId, blobSidecars)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return *blobSidecars, err
}

func (bn *BeaconClient) StateValidator(
	parentCtx context.Context,
	stateId eth2api.StateId,
	validatorId eth2api.ValidatorId,
) (*eth2api.ValidatorResponse, error) {
	var (
		stateValidatorResponse = new(eth2api.ValidatorResponse)
		exists                 bool
		err                    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.StateValidator(
		ctx,
		bn.api,
		stateId,
		validatorId,
		stateValidatorResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return stateValidatorResponse, err
}

type BeaconClients []*BeaconClient

// Return subset of clients that are currently running
func (all BeaconClients) Running() BeaconClients {
	res := make(BeaconClients, 0)
	for _, bc := range all {
		if bc.IsRunning() {
			res = append(res, bc)
		}
	}
	return res
}

// Return subset of clients that are part of an specific subnet
func (all BeaconClients) Subnet(subnet string) BeaconClients {
	if subnet == "" {
		return all
	}
	res := make(BeaconClients, 0)
	for _, bn := range all {
		if bn.Config.Subnet == subnet {
			res = append(res, bn)
		}
	}
	return res
}

// Returns comma-separated ENRs of all running beacon nodes
func (beacons BeaconClients) ENRs(parentCtx context.Context) (string, error) {
	if len(beacons) == 0 {
		return "", nil
	}
	enrs := make([]string, 0)
	for _, bn := range beacons {
		if bn.IsRunning() {
			enr, err := bn.ENR(parentCtx)
			if err != nil {
				return "", err
			}
			enrs = append(enrs, enr)
		}
	}
	return strings.Join(enrs, ","), nil
}

// Returns comma-separated P2PAddr of all running beacon nodes
func (beacons BeaconClients) P2PAddrs(
	parentCtx context.Context,
) (string, error) {
	if len(beacons) == 0 {
		return "", nil
	}
	staticPeers := make([]string, 0)
	for _, bn := range beacons {
		if bn.IsRunning() {
			p2p, err := bn.P2PAddr(parentCtx)
			if err != nil {
				return "", err
			}
			staticPeers = append(staticPeers, p2p)
		}
	}
	return strings.Join(staticPeers, ","), nil
}
