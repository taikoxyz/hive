package hivesim

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHiveFramework(t *testing.T) {
	handler, err := NewHiveFramework(&HiveConfig{
		ContainerOutput: true,
		DockerPull:      false,
		BaseDir:         "/Users/huan/projects/taiko/hive",
		SimPattern:      "taiko",
		SimTestPattern:  "taiko-deneb-testnet/test-deneb-genesis",
		Clients: []string{
			"taiko/geth",
			"taiko/prysm-bn",
			"taiko/prysm-vc",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		},
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}
