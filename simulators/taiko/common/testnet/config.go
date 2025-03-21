package testnet

import (
	"github.com/ethereum/hive/hivesim"
	"math/big"
	"taiko/common/clients"
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)

	MAINNET_SLOT_TIME int64 = 12
	MINIMAL_SLOT_TIME int64 = 6

	// Clients that support the minimal slot time in hive
	MINIMAL_SLOT_TIME_CLIENTS = []string{
		"lighthouse",
		"teku",
		"prysm",
		"lodestar",
		"grandine",
	}
)

type CreateConfig func(index int, nodes clients.Nodes) (hivesim.Params, error)

type Config struct {
	// Network
	Network string

	// 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail
	LogLevel string

	CreateConfig CreateConfig `json:"-"`

	// open debug flag
	DevDebug  bool
	DevExpose bool
}
