package clients

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

// L1Origin represents a L1Origin of a L2 block.
type L1Origin struct {
	BlockID       *big.Int    `json:"blockID" gencodec:"required"`
	L2BlockHash   common.Hash `json:"l2BlockHash"`
	L1BlockHeight *big.Int    `json:"l1BlockHeight" gencodec:"required"`
	L1BlockHash   common.Hash `json:"l1BlockHash" gencodec:"required"`
}

// MarshalJSON marshals as JSON.
func (l L1Origin) MarshalJSON() ([]byte, error) {
	type L1Origin struct {
		BlockID       *math.HexOrDecimal256 `json:"blockID" gencodec:"required"`
		L2BlockHash   common.Hash           `json:"l2BlockHash"`
		L1BlockHeight *math.HexOrDecimal256 `json:"l1BlockHeight" gencodec:"required"`
		L1BlockHash   common.Hash           `json:"l1BlockHash" gencodec:"required"`
	}
	var enc L1Origin
	enc.BlockID = (*math.HexOrDecimal256)(l.BlockID)
	enc.L2BlockHash = l.L2BlockHash
	enc.L1BlockHeight = (*math.HexOrDecimal256)(l.L1BlockHeight)
	enc.L1BlockHash = l.L1BlockHash
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (l *L1Origin) UnmarshalJSON(input []byte) error {
	type L1Origin struct {
		BlockID       *math.HexOrDecimal256 `json:"blockID" gencodec:"required"`
		L2BlockHash   *common.Hash          `json:"l2BlockHash"`
		L1BlockHeight *math.HexOrDecimal256 `json:"l1BlockHeight" gencodec:"required"`
		L1BlockHash   *common.Hash          `json:"l1BlockHash" gencodec:"required"`
	}
	var dec L1Origin
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.BlockID == nil {
		return errors.New("missing required field 'blockID' for L1Origin")
	}
	l.BlockID = (*big.Int)(dec.BlockID)
	if dec.L2BlockHash != nil {
		l.L2BlockHash = *dec.L2BlockHash
	}
	if dec.L1BlockHeight == nil {
		return errors.New("missing required field 'l1BlockHeight' for L1Origin")
	}
	l.L1BlockHeight = (*big.Int)(dec.L1BlockHeight)
	if dec.L1BlockHash == nil {
		return errors.New("missing required field 'l1BlockHash' for L1Origin")
	}
	l.L1BlockHash = *dec.L1BlockHash
	return nil
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
