package suite_base

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"taiko/common/clients"
	tn "taiko/common/testnet"
	"time"
)

var Deneb = "deneb"

func (ts BaseTestSpec) VerifyNodes(ctx context.Context, t *hivesim.T, testnet *tn.Testnet) {
	for _, node := range testnet.Nodes {
		ts.verify(ctx, t, node)
	}
}

func (ts BaseTestSpec) verify(ctx context.Context, t *hivesim.T, node *clients.Node) {
	var (
		l1Eth     = node.L1EthClient
		beacon    = node.BeaconClient
		validator = node.ValidatorClient
		l2Eth     = node.L2EthClient
		driver    = node.DriverClient
		proposer  = node.ProposerClient
		prover    = node.ProverClient
	)

	if l1Eth != nil && !l1Eth.IsRunning() {
		t.Fatalf("l1eth node is not running!")
		// Verify l1eth node run successfully.
		err := l1Eth.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l1geth number, err: %v", err)
		}
	}
	if beacon != nil && !beacon.IsRunning() {
		t.Fatalf("beacon node is not running!")
	}
	if validator != nil && !validator.IsRunning() {
		t.Fatalf("validator node is not running!")
	}

	if l2Eth != nil && !l2Eth.IsRunning() {
		t.Fatalf("l2eth node is not running!")
		// Verify l2eth node run successfully.
		err := l2Eth.VerifyNumber(ctx, time.Second*200, 1)
		if err != nil {
			t.Fatalf("failed to verify l2geth number, err: %v", err)
		}
	}
	if driver != nil && !driver.IsRunning() {
		t.Fatalf("driver node is not running!")
	}
	if proposer != nil && !proposer.IsRunning() {
		t.Fatalf("proposer node is not running!")
	}
	if prover != nil && !prover.IsRunning() {
		t.Fatalf("prover node is not running!")
	}
}
