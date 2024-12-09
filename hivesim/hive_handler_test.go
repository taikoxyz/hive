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
		runHive(t, "taiko-genesis/l2-full-sync", clientGroups)
	})

	t.Run("taiko-reorg/taiko-reorg", func(t *testing.T) {
		runHive(t, "taiko-reorg/taiko-reorg", [][]string{clientGroups[0]})
	})

	t.Run("taiko-blob/blob-l1-beacon", func(t *testing.T) {
		runHive(t, "taiko-blob/blob-l1-beacon", [][]string{
			{
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
			}})
	})

	t.Run("taiko-blob/blob-server", func(t *testing.T) {
		runHive(t, "taiko-blob/blob-server", [][]string{{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-reth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		}})
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
	// ./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "taiko-genesis/l2-snap-sync"
	t.Run("taiko-genesis/l2-snap-sync/clusters(3)", func(t *testing.T) {
		runHive(t, "taiko-genesis/l2-snap-sync", clientGroups)
	})
	// ./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "taiko-genesis/l2-full-sync"
	t.Run("taiko-genesis/l2-full-sync/clusters(3)", func(t *testing.T) {
		runHive(t, "taiko-genesis/l2-full-sync", clientGroups)
	})
	// ./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "taiko-reorg/taiko-reorg"
	t.Run("taiko-reorg/taiko-reorg", func(t *testing.T) {
		runHive(t, "taiko-reorg/taiko-reorg", [][]string{clientGroups[0]})
	})
	// ./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "taiko-blob/blob-l1-beacon"
	t.Run("taiko-blob/blob-l1-beacon", func(t *testing.T) {
		runHive(t, "taiko-blob/blob-l1-beacon", [][]string{{
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
		}})
	})
	// ./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,storage/redis,storage/postgres,blobscan/blobscan-api,blobscan/blobscan-indexer --sim taiko --sim.limit "taiko-blob/blob-server"
	t.Run("taiko-blob/blob-server", func(t *testing.T) {
		runHive(t, "taiko-blob/blob-server", [][]string{{
			"geth",
			"prysm/prysm-bn",
			"prysm/prysm-vc",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		}})
	})
	// ./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "taiko/preconf"
	t.Run("taiko/preconf", func(t *testing.T) {
		runHive(t, "taiko/preconf", [][]string{{
			"anvil",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		}})
	})
}

func runHive(t *testing.T, pattern string, clientGroups [][]string) {
	handler, err := NewHiveFramework(&HiveConfig{
		DockerPull:      true,
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
