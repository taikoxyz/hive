package clients

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"math/big"
)

func VerifyL1Origin(ctx context.Context, l2ReorgStartNumber uint64, l1client, l2client *rpc.EthClient) error {
	var l1Origins []*rawdb.L1Origin

	l2Number, err := l2client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	for number := l2ReorgStartNumber; number <= l2Number; number++ {
		l1Origin, err := l2client.L1OriginByID(ctx, big.NewInt(int64(number)))
		if err != nil {
			return err
		}
		l1Origins = append(l1Origins, l1Origin)
	}

	for _, l1Origin := range l1Origins {
		header, err := l1client.HeaderByNumber(ctx, l1Origin.L1BlockHeight)
		if err != nil {
			return err
		}
		if header.Hash() != l1Origin.L1BlockHash {
			return fmt.Errorf("l1Origin content is not right, blockID: %d, l1BlockHeight: %d, l1BlockHash: %s", l1Origin.BlockID.Uint64(), l1Origin.L1BlockHeight.Uint64(), header.Hash().String())
		}
	}

	return nil
}

func ParsePriv(data string) (*crypto.Secp256k1PrivateKey, error) {
	if len(data) > 2 && data[:2] == "0x" {
		data = data[2:]
	}
	b, err := hex.DecodeString(data)
	if err != nil {
		return nil, errors.New("p2p priv key is not formatted in hex chars")
	}
	p, err := crypto.UnmarshalSecp256k1PrivateKey(b)
	if err != nil {
		// avoid logging the priv key in the error, but hint at likely input length problem
		return nil, fmt.Errorf("failed to parse priv key from %d bytes", len(b))
	}
	return (p).(*crypto.Secp256k1PrivateKey), nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return ethrpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}
