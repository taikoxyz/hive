.PHONY: hive hivechain hiveview hive_alpine_dependency hive_go_dependency

all: hive hivechain hiveview hive_alpine_dependency hive_go_dependency

test_base: test_full_sync

test_reorg: test_reorg_reorg test_reorg_shorter test_reorg_longer

test_preconf: test_preconf_preconf test_preconf_forced_inclusion test_preconf_reorg

test_blob: test_blob_beacon test_blob_server

hive:
	@echo "Building hive..."
	@go build -o build/bin/hive ./

hivechain:
	@echo "Building hivechain..."
	@go build -o build/bin/hivechain ./cmd/hivechain

hiveview:
	@echo "Building hiveview..."
	@go build -o build/bin/hiveview ./cmd/hiveview

hive_alpine_dependency:
	@echo "Building hive_alpine_dependency..."
	@docker build --no-cache -t hive_alpine_dependency:latest --target hive_alpine_dependency .

hive_go_dependency:
	@echo "Building hive_go_dependency..."
	@docker build --no-cache -t hive_go_dependency:latest --target hive_go_dependency .

hive_prysmctl:
	@echo "Building prysmctl"
	@docker build --no-cache -t hive_prysmctl:latest clients/prysm/prysmctl

test_reorg_reorg: hive_prysmctl
	@echo "Running reorg/reorg_reorg test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "reorg/reorg_reorg"

test_reorg_longer: hive_prysmctl
	@echo "Running reorg/reorg_longer test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "reorg/reorg_longer"

test_reorg_shorter: hive_prysmctl
	@echo "Running reorg/reorg_shorter test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "reorg/reorg_shorter"

test_preconf_preconf: hive_prysmctl
	@echo "Running preconf/preconf test..."
	./build/bin/hive --docker.output --sim taiko --sim.limit "preconf/preconf" --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/taiko-geth,taiko/driver

test_preconf_reorg: hive_prysmctl
	@echo "Running preconf/reorg test..."
	./build/bin/hive --docker.output --sim taiko --sim.limit "preconf/reorg" --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer

test_preconf_forced_inclusion: hive_prysmctl
	@echo "Running preconf/forced-inclusion test..."
	./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer --sim taiko --sim.limit "preconf/forced-inclusion"

test_full_sync: hive_prysmctl
	@echo "Running base/fullsync test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "base/fullsync"

test_blob_beacon: hive_prysmctl
	@echo "Running blob/blob-l1-beacon test..."
	./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "blob/blob-l1-beacon"

test_blob_server: hive_prysmctl
	@echo "Running blob/blob-server test..."
	./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,storage/postgres,blobscan/blobscan-api,blobscan/blobscan-indexer --sim taiko --sim.limit "blob/blob-server"