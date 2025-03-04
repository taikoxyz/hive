package preconf

import (
	"context"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	params.SetEnvParams("TEST_L1_BEACON", "true")
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
	l2Geth.WaitLatestNumber(ctx, time.Second*30, proposer.PacayaClients.ForkHeight)

	// stop the proposer.
	proposer.PauseClient()

	events := proposer.GetProposeEvents(ctx, 0)
	var latestBlockId = events[len(events)-1].Info.LastBlockId
	if latestBlockId > proposer.PacayaClients.ForkHeight {
		l2Geth.WaitLatestNumber(ctx, time.Second*8, latestBlockId)
	}

	signedTxs, err := storeForcedInclusion(driver.Envs, driver.Client)
	t.FailIfNotNil(err, "cannot store forced inclusion")

	// restart the proposer.
	proposer.UnpauseClient()

	l1Number := events[len(events)-1].Raw.BlockNumber + 1
	for range time.Tick(time.Second) {
		select {
		case <-time.After(time.Second * 10):
			t.Fatalf("cannot get latest block")
		default:
		}
		events = proposer.GetProposeEvents(ctx, l1Number)
		if len(events) >= 1 {
			break
		}
	}

	t.Log("event index 0, blocks length: ", len(events[0].Info.Blocks), events[0].Info.LastBlockId)
	t.True(len(events[0].TxList) == 0, "forced inclusion txs should be empty")
	t.True(len(events[0].Info.BlobHashes) == 1, "forced inclusion blob hashes should be 1")
	t.True(len(events[0].Info.Blocks) == 1, "forced inclusion blocks should be 1")

	l2Geth.WaitLatestNumber(ctx, time.Second*6, events[0].Info.LastBlockId)

	l2Block, err := l2Geth.EthClient.BlockByNumber(ctx, big.NewInt(int64(events[0].Info.LastBlockId)))
	t.FailIfNotNil(err, "cannot get latest block")

	t.True(l2Block.Transactions().Len()-1 <= 512 && l2Block.Transactions().Len() == len(signedTxs)+1)

	txsMap := make(map[common.Hash]bool)
	for _, tx := range signedTxs {
		txsMap[tx.Hash()] = true
	}

	for _, tx := range l2Block.Transactions()[1:] {
		if _, ok := txsMap[tx.Hash()]; !ok {
			t.Fatalf("tx %s not found in block", tx.Hash())
		}
	}
}

func storeForcedInclusion(envs hivesim.Params, rpccli *rpc.Client) (types.Transactions, error) {
	envs["L1_PROPOSER_PRIV_KEY"] = params.ChainAuths[len(params.ChainAuths)-1].Key

	mockClient := &clients.MockClient{Envs: envs}
	if err := clients.NewTaikoClient(mockClient, flags.ProposerFlags); err != nil {
		return nil, err
	}

	signedTxs, err := utils.CreateL2Txs(context.Background(), rpccli.L2, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := utils.EncodeAndCompressTxList(signedTxs)
	if err != nil {
		return nil, err
	}
	blobs, err := clients.MakeBlobs(txBytes)
	if err != nil {
		return nil, err
	}

	fee, err := rpccli.PacayaClients.ForcedInclusionStore.FeeInGwei(nil)
	if err != nil {
		return nil, err
	}
	feeInGwei := big.NewInt(0).SetUint64(fee)
	feeInGwei = feeInGwei.Mul(feeInGwei, big.NewInt(eparams.GWei))

	abi, err := forcedinclusionstore.ForcedInclusionStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	blob := blobs[0]
	input, err := abi.Pack("storeForcedInclusion", uint8(0), uint32(0), uint32(len(txBytes)))
	if err != nil {
		return nil, err
	}

	to := common.HexToAddress(params.ParamByKey("FORCED_INCLUSION_STORE"))
	if _, err = mockClient.Send(context.Background(), txmgr.TxCandidate{
		TxData:   input,
		Blobs:    []*eth.Blob{blob},
		To:       &to,
		Value:    feeInGwei,
		GasLimit: mockClient.ProposeBlockTxGasLimit,
	}); err != nil {
		return nil, err
	}

	return signedTxs, nil
}
