.PHONY: hive hivechain hiveview docker_build docker_clean

all: hive hivechain hiveview

hive:
	@echo "Building hive..."
	@go build -o build/bin/hive ./

hivechain:
	@echo "Building hivechain..."
	@go build -o build/bin/hivechain ./cmd/hivechain

hiveview:
	@echo "Building hiveview..."
	@go build -o build/bin/hiveview ./cmd/hiveview

docker_build:
	@echo "Building docker images..."
	@docker build -t hive/clients/go-ethereum:latest ./clients/go-ethereum
	@docker build -t hive/clients/taiko-geth:latest ./clients/taiko-geth
	@docker build -t hive/clients/taiko-client:latest ./clients/taiko-client
	@docker build -t hive/clients/ethdevnet/validator:latest ./clients/ethdevnet/validator
	@docker build -t hive/clients/ethdevnet/geth:latest ./clients/ethdevnet/geth
	@docker build -t hive/clients/ethdevnet/beacon-chain:latest ./clients/ethdevnet/beacon-chain
	@docker build -t hive/simulators/devp2p:latest ./simulators/devp2p
	@docker build -t hive/simulators/ethereum/rpc:latest ./simulators/ethereum/rpc

docker_clean:
	@echo "Building docker images..."
	@docker images | grep "hive/clients/*" | awk '{print $3}' | xargs docker rmi
	@docker images | grep "hive/simulators/*" | awk '{print $3}' | xargs docker rmi

clean:
	@echo "Cleaning binaries..."
	@rm -rf build/bin/*
	@echo "Cleaning chain data..."
	@rm -rf chain
	@echo "Cleaning workspace..."
	@rm -rf workspace
	@echo "Cleaning clients docker images..."
	@docker images | grep "hive/clients/*" | awk '{print $3}' | xargs docker rmi
