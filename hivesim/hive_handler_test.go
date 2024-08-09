package hivesim

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCC(t *testing.T) {
	handler, err := NewHiveFramework(&HiveConfig{
		BaseDir:        "/Users/huan/projects/taiko/hive",
		SimPattern:     "taiko2",
		SimTestPattern: "eth2-deneb-testnet/test-deneb-genesis",
		Clients: []string{
			"taiko2/geth",
			"taiko2/prysm-bn",
			"taiko2/prysm-vc",
		},
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}
