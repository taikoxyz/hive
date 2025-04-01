package suite_reorg

import (
	"context"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"taiko/params"
	suite_base "taiko/suites/base"
	"time"
)

func init() {
	Tests = append(Tests,
		ReorgFinalizeTestSpec{
			ReorgTestSpec{BaseTestSpec: suite_base.BaseTestSpec{Name: "reorg_finalize"}},
		},
	)
}

type ReorgFinalizeTestSpec struct {
	ReorgTestSpec
}

func (r ReorgFinalizeTestSpec) Verify(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	var (
		node    = testnet.Nodes[0]
		anvil   = node.AnvilClient
		l2eth   = node.L2EthClient
		prover  = node.ProverClient
		propoer = node.ProposerClient
	)

	if err := node.Start(); err != nil {
		t.Fatalf("%s: failed to start the first node, err: %v", r.Name, err)
	}

	// For debug
	if testnet.DevDebug {
		node.DriverClient.Shutdown()
		time.Sleep(time.Hour * 2)
	}

	r.params = &clients.ReorgParams{
		DelayTime:   1,
		DelayNumber: 0,
	}

	// Start recording l1 chain.
	anvil.StartRecordReorgPoints(ctx, l2eth.EthClient)

	// Verify blocks.
	var lastVerifiedBlockID uint64
	for range time.Tick(time.Second) {
		prover.VerifyBlocks(params.L1Auths[0])
		lastVerifiedBlockID = prover.GetLastVerifiedBlockId(ctx)
		if lastVerifiedBlockID >= propoer.PacayaClients.ForkHeight-1 {
			node.ProverClient.Shutdown()
			break
		}
	}

	// Verify l1 origin
	t.FailIfNotNil(verifyL1Origin(ctx, lastVerifiedBlockID, anvil.EthClient, l2eth.EthClient), "check finalize status before reorg")

	// Set reorg point before the fork height.
	r.reorgStart = lastVerifiedBlockID / 2
	r.reorgDeep = 4

	// Reorg l1 chain.
	r.reorg(ctx, t, node)

	// Verify l1 origin
	t.FailIfNotNil(verifyL1Origin(ctx, lastVerifiedBlockID, anvil.EthClient, l2eth.EthClient), "check finalize status after reorg")

	// verify finalize
	prover.WaitLatestVerifiedNumber(ctx, time.Second*5, lastVerifiedBlockID)
}

func verifyL1Origin(ctx context.Context, lastVerifiedBlockID uint64, l1cli, l2cli *rpc.EthClient) error {
	fHeader, err := l2cli.HeaderByNumber(ctx, big.NewInt(int64(lastVerifiedBlockID)))
	if err != nil {
		return err
	}
	if fHeader.Number.Uint64() != lastVerifiedBlockID {
		return fmt.Errorf("the finalize number is not right, expect_number: %d, actual_number: %d", lastVerifiedBlockID, fHeader.Number.Uint64())
	}

	for number := uint64(1); number <= lastVerifiedBlockID; number++ {
		l1Origin, err := l2cli.L1OriginByID(ctx, big.NewInt(int64(number)))
		if err != nil {
			return fmt.Errorf("l1 origin: %w", err)
		}
		if l1Origin.L1BlockHeight == nil {
			return fmt.Errorf("l1 origin block height is nil, l2_number: %d", number)
		}
		header, err := l1cli.HeaderByNumber(ctx, l1Origin.L1BlockHeight)
		if err != nil {
			return err
		}
		if header.Hash() != l1Origin.L1BlockHash {
			return fmt.Errorf("the finalize hash is not right, expect_hash: %s, actual_hash: %s", header.Hash(), l1Origin.L1BlockHash.String())
		}
	}
	return nil
}
