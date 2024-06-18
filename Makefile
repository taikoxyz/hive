.PHONY: hive hivechain hiveview

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

docker:
	@echo "Building docker images..."
	@docker build -t hive/clients/taiko-geth:latest ./clients/taiko-geth
	@docker build -t hive/clients/go-ethereum:latest ./clients/go-ethereum
	@docker build -t hive/clients/prysm-bn:latest ./clients/prysm-bn
	@docker build -t hive/clients/prysm-vc:latest ./clients/prysm-vc

clean:
	@echo "Cleaning binaries..."
	@rm -rf build/bin/*
	@echo "Cleaning chain data..."
	@rm -rf chain
	@echo "Cleaning workspace..."
	@rm -rf workspace
	@echo "Cleaning clients docker images..."
	@docker images | grep "hive/clients/*" | awk '{print $3}' | xargs docker rmi
