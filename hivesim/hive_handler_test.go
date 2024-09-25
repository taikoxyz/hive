package hivesim

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaikoRethHandler(t *testing.T) {
	clientGroups := [][]string{
		{
			"anvil",
			"taiko/taiko-reth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		},
		{
			"taiko/taiko-reth",
			"taiko/driver",
		},
		{
			"taiko/taiko-reth",
			"taiko/driver",
		},
	}

	// Multi clusters test.
	t.Run("taiko-genesis/l2-full-sync/clusters(3)", func(t *testing.T) {
		testDenebGenesis(t, "taiko-genesis/l2-full-sync", clientGroups)
	})

	t.Run("taiko-reorg/taiko-reorg", func(t *testing.T) {
		testDenebReorg(t, "taiko-reorg/taiko-reorg", [][]string{clientGroups[0]})
	})

	t.Run("taiko-blob/blob-l1-beacon", func(t *testing.T) {
		testBlobScan(t, "taiko-blob/blob-l1-beacon", []string{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-reth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
			"storage/redis",
			"storage/postgres",
			"blobscan/blobscan-api",
			"blobscan/blobscan-indexer",
		})
	})

	t.Run("taiko-blob/blob-server", func(t *testing.T) {
		testBlobScan(t, "taiko-blob/blob-server", []string{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-reth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		})
	})
}

func TestTaikoGethHandler(t *testing.T) {
	clientGroups := [][]string{
		{
			"anvil",
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

	t.Run("taiko-reorg/taiko-reorg", func(t *testing.T) {
		testDenebReorg(t, "taiko-reorg/taiko-reorg", [][]string{clientGroups[0]})
	})

	t.Run("taiko-blob/blob-l1-beacon", func(t *testing.T) {
		testBlobScan(t, "taiko-blob/blob-l1-beacon", []string{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
			"storage/redis",
			"storage/postgres",
			"blobscan/blobscan-api",
			"blobscan/blobscan-indexer",
		})
	})

	t.Run("taiko-blob/blob-server", func(t *testing.T) {
		testBlobScan(t, "taiko-blob/blob-server", []string{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		})
	})
}

func testBlobScan(t *testing.T, pattern string, clients []string) {
	handler, err := NewHiveFramework(&HiveConfig{
		BuildOutput:     false,
		ContainerOutput: true,
		BaseDir:         "/Users/huan/projects/taiko/hive",
		SimPattern:      "taiko",
		SimTestPattern:  pattern,
		ClientGroups:    [][]string{clients},
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
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
