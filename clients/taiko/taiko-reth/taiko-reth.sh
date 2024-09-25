#!/bin/bash

# Set default values
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

FLAGS="-d --datadir=/reth/data "
FLAGS="$FLAGS --chain=$HIVE_NETWORK_ID "
FLAGS="$FLAGS --http --http.addr=0.0.0.0 --http.corsdomain=* --http.api=admin,debug,eth,net,web3,txpool,taiko "
FLAGS="$FLAGS --ws --ws.addr=0.0.0.0 --ws.origins=* --ws.api=admin,debug,eth,net,web3,txpool,taiko "

FLAGS="$FLAGS --authrpc.addr=0.0.0.0 "
if [ "$HIVE_TAIKO_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO_JWT_SECRET" >/reth/jwt.hex
  FLAGS="$FLAGS --authrpc.jwtsecret=/reth/jwt.hex "
fi

if [ "$HIVE_BOOTNODE" != "" ]; then
  FLAGS="$FLAGS --bootnodes=$HIVE_BOOTNODE"
fi

echo FLAGS: "$FLAGS"

echo "Starting taiko reth..."
taiko-reth node $FLAGS
