package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
)

func (ec *EthNode) BlockNumber(ctx context.Context) *big.Int {
	number, err := ec.EthClient.BlockNumber(ctx)
	ec.FailIfNotNil(err, fmt.Sprintf("%s: failed to get block number, chainID: %d, err: %v", ec.ClientType(), ec.EthClient.ChainID.Uint64(), err))
	return big.NewInt(int64(number))
}

func (ec *EthNode) RevertTaikoGeth(ctx context.Context, number *big.Int) {
	header, err := ec.EthClient.HeaderByNumber(ctx, number)
	ec.FailIfNotNil(err, fmt.Sprintf("%s: failed to get block header", ec.ClientType()))

	authClient, err := rpc.NewJWTEngineClient(ec.EngineURL(), string(common.Hex2Bytes("c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7")))
	ec.FailIfNotNil(err, fmt.Sprintf("%s: failed to create auth client", ec.ClientType()))

	result, err := authClient.ForkchoiceUpdate(
		ctx,
		&engine.ForkchoiceStateV1{
			HeadBlockHash: header.Hash(),
		},
		nil,
	)
	ec.FailIfNotNil(err, fmt.Sprintf("%s: failed to forkchoice updated to latest block, target_number: %d", ec.ClientType(), number.Uint64()))

	if result.PayloadStatus.Status != engine.VALID {
		ec.Fatalf("%s: failed to update forkchoice to latest block, target_number: %d, status: %s", ec.ClientType(), number.Uint64(), result.PayloadStatus.Status)
	}
}
