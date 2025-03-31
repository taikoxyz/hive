package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
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
