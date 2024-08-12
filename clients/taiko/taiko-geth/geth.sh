#!/bin/bash

# Set default values
HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

if [ "$HIVE_TAIKO_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO_JWT_SECRET" >/geth/jwt.hex
fi

echo "Starting taiko geth..."
geth \
  --taiko \
  --datadir=/geth/data \
  --networkid="$HIVE_NETWORK_ID" \
  --verbosity="$HIVE_LOGLEVEL" \
  --state.scheme=path \
  --http \
  --http.api=debug,eth,net,web3,txpool,miner,taiko \
  --http.addr=0.0.0.0 \
  --http.vhosts=* \
  --http.corsdomain=* \
  --ws \
  --ws.api=debug,eth,net,web3,txpool,miner,taiko \
  --ws.addr=0.0.0.0 \
  --ws.origins=* \
  --authrpc.vhosts=* \
  --authrpc.addr=0.0.0.0 \
  --authrpc.jwtsecret=/geth/jwt.hex \
  --allow-insecure-unlock \
  --nodiscover
