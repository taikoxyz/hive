package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// L1Origin represents a L1Origin of a L2 block.
type L1Origin struct {
	BlockID       *big.Int    `json:"blockID" gencodec:"required"`
	L2BlockHash   common.Hash `json:"l2BlockHash"`
	L1BlockHeight *big.Int    `json:"l1BlockHeight" gencodec:"required"`
	L1BlockHash   common.Hash `json:"l1BlockHash" gencodec:"required"`
}

type TaikoGethClient struct {
	*EthNode
}

func (t *TaikoGethClient) L1OriginByID(ctx context.Context, blockID *big.Int) (*L1Origin, error) {
	rpcClient := t.RPClient()
	var res *L1Origin
	if err := rpcClient.CallContext(ctx, &res, "taiko_l1OriginByID", hexutil.EncodeBig(blockID)); err != nil {
		return nil, err
	}

	return res, nil
}
