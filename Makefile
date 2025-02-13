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
