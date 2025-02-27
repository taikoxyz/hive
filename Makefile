.PHONY: hive hivechain hiveview hive_alpine_dependency hive_go_dependency

all: hive hivechain hiveview hive_alpine_dependency hive_go_dependency

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
	docker build --no-cache -t hive_alpine_dependency:latest --target hive_alpine_dependency .

hive_go_dependency:
	docker build --no-cache -t hive_go_dependency:latest --target hive_go_dependency .

test_reorg_reorg:
	@echo "Running reorg/reorg test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "reorg/reorg"

test_preconf_preconf:
	@echo "Running preconf/preconf test..."
	./build/bin/hive --docker.output --sim taiko --sim.limit "preconf/preconf" --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/taiko-geth,taiko/driver

test_preconf_reorg:
	@echo "Running preconf/reorg test..."
	./build/bin/hive --docker.output --sim taiko --sim.limit "preconf/reorg" --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer

test_preconf_forced_inclusion:
	@echo "Running preconf/forced-inclusion test..."
	#./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer --sim taiko --sim.limit "preconf/forced-inclusion"
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer --sim taiko --sim.limit "preconf/forced-inclusion"

test_full_sync:
	@echo "Running base/fullsync test..."
	./build/bin/hive --docker.output --client anvil,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,taiko/taiko-geth,taiko/driver,taiko/taiko-geth,taiko/driver --sim taiko --sim.limit "base/fullsync"

test_blob_l1_beacon:
	@echo "Running blob/blob-l1-beacon test..."
	./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover --sim taiko --sim.limit "blob/blob-l1-beacon"

test_blob_server:
	@echo "Running blob/blob-server test..."
	./build/bin/hive --docker.output --client geth,prysm/prysm-bn,prysm/prysm-vc,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover,storage/redis,storage/postgres,blobscan/blobscan-api,blobscan/blobscan-indexer --sim taiko --sim.limit "blob/blob-server"