#!/bin/sh

set -e

DATA_DIR=/l1-beacon-chain

if [ ! -f "$DATA_DIR/GENESIS_LOCK" ];then
  # Add lock file.
  mkdir -p $DATA_DIR
  touch "$DATA_DIR/GENESIS_LOCK"

  # Get current time.
  GENESIS_TIME=$(echo "$(date +%s) / 100 * 100" | bc)
  echo "CURRENT_TIME=$(date +%s)"
  echo "GENESIS_TIME=$GENESIS_TIME"

  # Reset genesis.ssz
  prysmctl \
    testnet \
    generate-genesis \
    --fork=deneb \
    --num-validators=64 \
    --genesis-time="$GENESIS_TIME" \
    --genesis-time-delay=100 \
    --output-ssz=$DATA_DIR/genesis.ssz \
    --chain-config-file=/config.yml \
    --geth-genesis-json-in=/genesis.json \
    --geth-genesis-json-out=/genesis.json

  cat /genesis.json
fi

beacon-chain \
  --verbosity=debug \
  --datadir=$DATA_DIR \
  --min-sync-peers=0 \
  --genesis-state=$DATA_DIR/genesis.ssz \
  --interop-eth1data-votes \
  --chain-config-file=/config.yml \
  --contract-deployment-block=0 \
  --chain-id=32382 \
  --rpc-host=0.0.0.0 \
  --grpc-gateway-host=0.0.0.0 \
  --execution-endpoint=http://localhost:8551 \
  --accept-terms-of-use \
  --jwt-secret=/jwtsecret \
  --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
  --minimum-peers-per-subnet=0 \
  --enable-debug-rpc-endpoints \
  --blob-retention-epochs=409600
