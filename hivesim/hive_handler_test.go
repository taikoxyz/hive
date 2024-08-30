package hivesim

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHiveHandler(t *testing.T) {
	clientGroups := [][]string{
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
		{
			"taiko/taiko-geth",
			"taiko/driver",
		},
	}

	// Multi clusters test.
	t.Run("taiko-genesis/l2-snap-sync/clusters(3)", func(t *testing.T) {
		testDenebGenesis(t, "taiko-genesis/l2-snap-sync", clientGroups)
	})
	t.Run("taiko-genesis/l2-full-sync/clusters(3)", func(t *testing.T) {
		testDenebGenesis(t, "taiko-genesis/l2-full-sync", clientGroups)
	})

	t.Run("taiko-reorg/taiko-reorg/clusters(3)", func(t *testing.T) {
		testDenebReorg(t, "taiko-reorg/taiko-reorg", clientGroups)
	})
}

func testDenebGenesis(t *testing.T, pattern string, clientGroups [][]string) {
	handler, err := NewHiveFramework(&HiveConfig{
		BuildOutput:     false,
		ContainerOutput: true,
		BaseDir:         "/Users/huan/projects/taiko/hive",
		SimPattern:      "taiko",
		SimTestPattern:  pattern,
		ClientGroups:    clientGroups,
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}

func testDenebReorg(t *testing.T, pattern string, clientGroups [][]string) {
	handler, err := NewHiveFramework(&HiveConfig{
		BuildOutput:     true,
		ContainerOutput: true,
		BaseDir:         "/Users/huan/projects/taiko/hive",
		SimPattern:      "taiko",
		SimTestPattern:  pattern,
		ClientGroups:    clientGroups,
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}
