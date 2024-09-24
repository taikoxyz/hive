#!/bin/bash

# Immediately abort the script on any error encountered
set -e

FLAGS="-d --datadir=/reth/data "
FLAGS="$FLAGS --chain=internal_devnet_a "
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
