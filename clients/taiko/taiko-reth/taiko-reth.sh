#!/bin/bash

# Set default values
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

FLAGS="node -d --datadir=/reth/data --engine.legacy "

# Add network id
case "$HIVE_NETWORK_ID" in
  167000)
    FLAGS="$FLAGS --chain=main "
    ;;
  167001)
    FLAGS="$FLAGS --chain=internal-l2-a "
    ;;
  167002)
    FLAGS="$FLAGS --chain=internal-l2-b "
    ;;
esac

FLAGS="$FLAGS --http --http.addr=0.0.0.0 --http.corsdomain=* --http.api=admin,debug,eth,net,web3,txpool,taiko "
FLAGS="$FLAGS --ws --ws.addr=0.0.0.0 --ws.origins=* --ws.api=admin,debug,eth,net,web3,txpool,taiko "

FLAGS="$FLAGS --authrpc.addr=0.0.0.0 --authrpc.jwtsecret=/tmp/jwt.hex"

if [ "$HIVE_BOOTNODE" != "" ]; then
  FLAGS="$FLAGS --port 30303 --addr 0.0.0.0 --trusted-peers=$HIVE_BOOTNODE " # --enable-discv5-discovery --trusted-peers
fi

echo FLAGS: "$FLAGS"

echo "Starting taiko reth..."
nohup taiko-reth $FLAGS 2>&1 &

touch /reth/nohup.out
tail -f /reth/nohup.out
