package chaingenerators

import (
	"github.com/ethereum/go-ethereum/core/types"
	el "taiko2/common/config/execution"
)

type ChainGenerator interface {
	Generate(*el.ExecutionGenesis) ([]*types.Block, error)
}
