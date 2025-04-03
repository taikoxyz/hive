package clients

import (
	"context"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	preconfBlocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	preconfblocks "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/preconf_blocks"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
	"os"
	"taiko/params"
)

func (d *DriverClient) GetPreconfP2PNode() string {
	sk, err := ParsePriv(params.ChainAuths[d.Index].Key)
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to parse private key %d", d.ClientType(), d.Index))

	idB, err := peer.IDFromPublicKey(sk.GetPublic())
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to parse public key %d", d.ClientType(), d.Index))

	return fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", d.NetworkIP(), PreconfP2pPort, idB)
}

func (d *DriverClient) PublishL2Payload(ctx context.Context, msg *preconfblocks.BuildPreconfBlockRequestBody) {
	difficulty, err := encoding.CalculatePacayaDifficulty(new(big.Int).SetUint64(msg.ExecutableData.Number))
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to calculate difficulty", d.ClientType()))

	baseFee := uint256.NewInt(msg.ExecutableData.BaseFeePerGas)
	err = d.p2pNode.GossipOut().PublishL2Payload(ctx,
		&eth.ExecutionPayloadEnvelope{
			ExecutionPayload: &eth.ExecutionPayload{
				ParentHash:    msg.ExecutableData.ParentHash,
				FeeRecipient:  msg.ExecutableData.FeeRecipient,
				PrevRandao:    eth.Bytes32(difficulty[:]),
				BlockNumber:   eth.Uint64Quantity(msg.ExecutableData.Number),
				GasLimit:      eth.Uint64Quantity(msg.ExecutableData.GasLimit),
				Timestamp:     eth.Uint64Quantity(msg.ExecutableData.Timestamp),
				ExtraData:     eth.BytesMax32(msg.ExecutableData.ExtraData),
				BaseFeePerGas: eth.Uint256Quantity(*baseFee),
				Transactions:  []eth.Data{msg.ExecutableData.Transactions},
			},
		},
		d.p2pSigner,
	)
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to publish payload", d.ClientType()))
}

func GetP2PNode(ctx context.Context, staticPeer string) (*p2p.NodeP2P, p2p.Signer, error) {
	// Set static peer.
	_ = os.Setenv("PRECONFIRMATION_P2P_STATIC", staticPeer)

	driverClient := &MockDriver{}
	if err := NewTaikoClient(driverClient, flags.DriverFlags); err != nil {
		return nil, nil, fmt.Errorf("failed to create taiko client: %w", err)
	}
	cfg := driverClient.Config

	client, err := rpc.NewClient(ctx, cfg.ClientConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init taiko client: %w", err)
	}

	preconfBlockServer, err := preconfBlocks.New(
		cfg.PreconfBlockServerCORSOrigins,
		cfg.PreconfBlockServerJWTSecret,
		&MockPreconfBlockChainSyncer{},
		client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create preconf block server: %w", err)
	}

	p2pNode, err := p2p.NewNodeP2P(
		ctx,
		&rollup.Config{L1ChainID: client.L1.ChainID, L2ChainID: client.L2.ChainID, Taiko: true},
		log.Root(),
		driverClient.P2PConfigs,
		preconfBlockServer,
		nil,
		preconfBlockServer,
		metrics.NewMetrics("client"),
		false,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create p2p node: %w", err)
	}

	p2pSigner, err := cfg.P2PSignerConfigs.SetupSigner(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create p2p signer: %w", err)
	}

	return p2pNode, p2pSigner, nil
}
