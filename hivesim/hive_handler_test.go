package hivesim

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHiveFramework(t *testing.T) {
	handler, err := NewHiveFramework(&HiveConfig{
		BuildOutput:     true,
		ContainerOutput: true,
		DockerPull:      false,
		BaseDir:         "/Users/huan/projects/taiko/hive",
		SimPattern:      "taiko",
		SimTestPattern:  "taiko-deneb-testnet/test-deneb-genesis",
		ClientGroups: [][]string{
			{
				"taiko/anvil",
				"taiko/taiko-geth",
				"taiko/driver",
				"taiko/proposer",
				"taiko/prover",
			},
			{
				"taiko/taiko-geth",
				"taiko/driver",
			},
		},
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}
