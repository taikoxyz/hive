package suite_base

import (
	"github.com/ethereum/go-ethereum/common"
	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"taiko2/common/utils"
)

var Deneb = "deneb"

var (
	blobTxAccounts   = utils.TestAccounts[:len(utils.TestAccounts)/2]
	normalTxAccounts = utils.TestAccounts[len(utils.TestAccounts)/2:]
)

func WithdrawalAddress(vI beacon.ValidatorIndex) common.Address {
	return common.Address{byte(vI + 0x100)}
}
