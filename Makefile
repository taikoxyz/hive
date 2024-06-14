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

clean:
	@echo "Cleaning..."
	@rm -rf build/bin/*
	@rm -rf chain
	@rm -rf workspace
