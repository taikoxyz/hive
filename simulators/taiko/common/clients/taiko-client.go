package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TaikoClientConfig struct {
	BeaconSync bool
}

type DriverClient struct {
	*HiveManagedClient
}

type ProposerClient struct {
	*HiveManagedClient
}

type ProverClient struct {
	*HiveManagedClient
}

type TaikoClient struct {
	l1Client *ethclient.Client
}

func (t *TaikoClient) GetTiers(ctx context.Context) {

}
