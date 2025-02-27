package preconf

import (
	"context"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	eparams "github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"taiko/bindings/pacaya/forcedinclusionstore"
	"taiko/common/clients"
	"taiko/common/testnet"
	tn "taiko/common/testnet"
	"taiko/common/utils"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests, &ForcedInclusionTestSpec{
		PreconfTestSpec: PreconfTestSpec{
			suite_base.BaseTestSpec{
				Name:  "forced-inclusion",
				Debug: false,
			},
		},
	})
}

type ForcedInclusionTestSpec struct {
	PreconfTestSpec
}

func (f *ForcedInclusionTestSpec) GetTestnetConfig() *testnet.Config {
	cfg := f.PreconfTestSpec.GetTestnetConfig()
	cfg.Network = "network_forced_inclusion"

	return cfg
}

func (f *ForcedInclusionTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node     = testnet.Nodes[0]
		l2Geth   = node.L2EthClient
		driver   = node.DriverClient
		proposer = node.ProposerClient
	)

	if err := node.Start(); err != nil {
		t.Fatalf("BaseTestSpec failed to start 0 node: %v", err)
	}

	// Wait until over the pacaya hardfork number.
	t.FailIfNotNil(l2Geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight), fmt.Sprintf("cannot wait for pacaya hardfork number: %d", proposer.PacayaClients.ForkHeight))

	// For Debug
	if f.Debug {
		proposer.Shutdown()
		time.Sleep(time.Minute * 120)
	}

	err := storeForcedInclusion(driver.Envs, driver.Client)
	t.FailIfNotNil(err, "cannot store forced inclusion")

	for range time.Tick(time.Second) {
		select {
		case <-time.After(30 * time.Second):
			t.Fatalf("cannot store forced inclusion")
		default:
		}
		head1, err := driver.PacayaClients.ForcedInclusionStore.Head(nil)
		t.FailIfNotNil(err, "cannot get forced inclusion store head")
		if head1 > 0 {
			t.Logf("forced inclusion stored, new head: %d", head1)
			return
		}
	}
}

func storeForcedInclusion(envs hivesim.Params, rpccli *rpc.Client) error {
	envs["L1_PROPOSER_PRIV_KEY"] = params.ChainAuths[len(params.ChainAuths)-1].Key

	mockClient := &clients.MockClient{Envs: envs}
	if err := clients.NewTaikoClient(mockClient, flags.ProposerFlags); err != nil {
		return err
	}

	signedTxs, err := utils.CreateL2Txs(context.Background(), rpccli.L2, true)
	if err != nil {
		return err
	}

	txBytes, err := utils.EncodeAndCompressTxList(signedTxs)
	if err != nil {
		return err
	}
	blobs, err := clients.MakeBlobs(txBytes)
	if err != nil {
		return err
	}

	fee, err := rpccli.PacayaClients.ForcedInclusionStore.FeeInGwei(nil)
	if err != nil {
		return err
	}
	feeInGwei := big.NewInt(0).SetUint64(fee)
	feeInGwei = feeInGwei.Mul(feeInGwei, big.NewInt(eparams.GWei))

	abi, err := forcedinclusionstore.ForcedInclusionStoreMetaData.GetAbi()
	if err != nil {
		return err
	}

	for index, blob := range blobs {
		input, err := abi.Pack("storeForcedInclusion", uint8(index), uint32(0), uint32(eth.BlobSize))
		if err != nil {
			return err
		}

		to := common.HexToAddress(params.ParamByKey("FORCED_INCLUSION_STORE"))
		if _, err = mockClient.Send(context.Background(), txmgr.TxCandidate{
			TxData:   input,
			Blobs:    []*eth.Blob{blob},
			To:       &to,
			Value:    feeInGwei,
			GasLimit: mockClient.ProposeBlockTxGasLimit,
		}); err != nil {
			return err
		}
	}

	return nil
}
