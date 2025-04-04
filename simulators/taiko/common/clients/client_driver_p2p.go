package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/hivesim"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"log/slog"
	"os"
	"taiko/params"
	"time"
)

func (d *DriverClient) GetPreconfP2PNode() string {
	sk, err := ParsePriv(params.ChainAuths[d.Index].Key)
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to parse private key %d", d.ClientType(), d.Index))

	idB, err := peer.IDFromPublicKey(sk.GetPublic())
	d.FailIfNotNil(err, fmt.Sprintf("%s: failed to parse public key %d", d.ClientType(), d.Index))

	return fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", d.NetworkIP(), PreconfP2pPort, idB)
}

type P2PNode struct {
	*p2p.NodeP2P
	p2pSigner p2p.Signer
}

func NewP2PNode(ctx context.Context, client *rpc.Client, index int, staticPeer string) (*P2PNode, error) {

	sk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKey := common.Bytes2Hex(crypto.FromECDSA(sk))

	envs := hivesim.Params{}
	envs["PRECONFIRMATION_P2P_PRIV_RAW"] = privateKey
	envs["PRECONFIRMATION_P2P_STATIC"] = staticPeer
	envs["PRECONFIRMATION_P2P_LISTEN_TCP_PORT"] = fmt.Sprintf("%d", PreconfP2pPort)
	envs["PRECONFIRMATION_P2P_SEQUENCER_KEY"] = params.ParamByKey("L1_PROPOSER_PRIV_KEY")
	for k, v := range envs {
		_ = os.Setenv(k, v)
	}

	driverClient := &MockDriver{}
	if err := NewTaikoClient(
		driverClient, flags.DriverFlags,
		"--p2p.listen.tcp", fmt.Sprintf("%d", PreconfP2pPort),
		"--p2p.discovery.path", "memory",
		"--p2p.peerstore.path", "memory",
		"--p2p.sequencer.key", params.ParamByKey("L1_PROPOSER_PRIV_KEY"),
		"--p2p.priv.raw", privateKey,
		"--p2p.bootnodes", "enode://869d07b5932f17e8490990f75a3f94195e9504ddb6b85f7189e5a9c0a8fff8b00aecf6f3ac450ecba6cdabdb5858788a94bde2b613e0f2d82e9b395355f76d1a@34.65.67.101:30305",
		"--p2p.static", staticPeer,
	); err != nil {
		return nil, fmt.Errorf("failed to create taiko client: %w", err)
	}
	cfg := driverClient.Config

	preconfBlockServer := &MockPreconfBlockAPIServer{Address: common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")}

	glogger := log.NewGlogHandler(log.NewTerminalHandler(os.Stdout, false))
	glogger.Verbosity(slog.LevelInfo)
	log.SetDefault(log.NewLogger(glogger))

	p2pNode, err := p2p.NewNodeP2P(
		ctx,
		&rollup.Config{L1ChainID: client.L1.ChainID, L2ChainID: client.L2.ChainID, Taiko: true},
		log.Root(),
		cfg.P2PConfigs,
		preconfBlockServer,
		nil,
		preconfBlockServer,
		metrics.NewMetrics("client"),
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create p2p node: %w", err)
	}

	p2pSigner, err := cfg.P2PSignerConfigs.SetupSigner(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create p2p signer: %w", err)
	}

	fmt.Printf("P2PNode config, host: %s, static_peer: %s\n", p2pNode.Host().ID().String(), cfg.P2PConfigs.StaticPeers)

	node := &P2PNode{
		NodeP2P:   p2pNode,
		p2pSigner: p2pSigner,
	}

	return node, nil
}

func (p *P2PNode) Close() {
	p.NodeP2P.Close()
	p.p2pSigner.Close()
}

func (p *P2PNode) PublishL2Payload(ctx context.Context, sendBody *eth.ExecutionPayloadEnvelope) error {
	body, _ := json.Marshal(sendBody)
	fmt.Printf("successfully published preconf block, peer_count: %d, content: %s\n", len(p.Peers()), string(body))

	err := p.GossipOut().PublishL2Payload(ctx,
		sendBody,
		p.p2pSigner,
	)
	if err != nil {
		return fmt.Errorf("failed to publish payload: %w", err)
	}

	return nil
}

func (p *P2PNode) WaitConnected(ctx context.Context, timeout time.Duration) error {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for p2p node connection")
		case <-tick.C:
			if len(p.Peers()) > 0 {
				return nil
			}
		}
	}
}

func GetP2PEnvs(index int) hivesim.Params {
	var envs = hivesim.Params{}
	envs["PRECONFIRMATION_SERVER_PORT"] = fmt.Sprintf("%d", PreconfServerPort)
	envs["PRECONFIRMATION_SERVER_SIGNATURE_CHECK"] = "true"
	envs["PRECONFIRMATION_P2P_DISCOVERY_PATH"] = "memory"
	envs["PRECONFIRMATION_P2P_PEERSTORE_PATH"] = "memory"
	envs["PRECONFIRMATION_P2P_SEQUENCER_KEY"] = envs["L1_PROPOSER_PRIV_KEY"]
	envs["PRECONFIRMATION_P2P_PRIV_RAW"] = params.ChainAuths[index].Key
	envs["PRECONFIRMATION_P2P_BOOTNODES"] = "enode://869d07b5932f17e8490990f75a3f94195e9504ddb6b85f7189e5a9c0a8fff8b00aecf6f3ac450ecba6cdabdb5858788a94bde2b613e0f2d82e9b395355f76d1a@34.65.67.101:30305"

	return envs
}
