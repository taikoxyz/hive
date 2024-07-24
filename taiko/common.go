package taiko

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/hive/hivesim"
	"math/big"
)

var networkCreated = make(map[hivesim.SuiteID]bool)

// CreateOrConnectNetwork ensures there is a separate network to be able to send the client traffic
// from two separate IP addrs.
func CreateOrConnectNetwork(t *hivesim.T, container, network string) {
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true
	}

	if err := t.Sim.ConnectContainer(t.SuiteID, network, container); err != nil {
		t.Fatal("can't connect container to network:", err)
	}
}

// Revert roll back to specified height.
func Revert(ctx context.Context, url, jwt string, number uint64) error {
	cli, err := rpc.DialOptions(ctx, url, rpc.WithHTTPAuth(node.NewJWTAuth(common.HexToHash(jwt))))
	if err != nil {
		return err
	}

	client := ethclient.NewClient(cli)
	header, err := client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(number))
	if err != nil {
		return err
	}

	update := engine.ForkchoiceStateV1{
		HeadBlockHash:      header.Hash(),
		FinalizedBlockHash: header.Hash(),
		SafeBlockHash:      header.Hash(),
	}

	var result engine.ForkChoiceResponse
	if err = cli.CallContext(ctx, &result, "engine_forkchoiceUpdatedV2", update, nil); err != nil {
		return err
	}
	if result.PayloadStatus.Status != engine.VALID {
		return fmt.Errorf("payload status is not valid")
	}

	return nil
}
