#!/bin/bash

# Set default values
HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

FLAGS="--taiko --taiko --datadir /geth/data --state.scheme path"
# Configure http.
FLAGS="$FLAGS --http --http.addr 0.0.0.0 --http.api admin,debug,eth,miner,net,personal,txpool,web3"
# Configure ws.
FLAGS="$FLAGS --ws --ws.addr 0.0.0.0 --ws.api admin,debug,eth,miner,net,personal,txpool,web3"
# Configure auth rpc
FLAGS="$FLAGS --authrpc.addr 0.0.0.0 --authrpc.port 8551 --allow-insecure-unlock --authrpc.jwtsecret /geth/jwt.hex"
# Configure log level
FLAGS="$FLAGS --verbosity $HIVE_LOGLEVEL"
# If a specific network ID is requested, use that
FLAGS="$FLAGS --networkid $HIVE_NETWORK_ID"

# It doesn't make sense to dial out, use only a pre-set bootnode.
if [ "$HIVE_BOOTNODE" != "" ]; then
  FLAGS="$FLAGS --bootnodes $HIVE_BOOTNODE"
else
  FLAGS="$FLAGS --nodiscover"
fi

# Handle any client mode or operation requests
if [ "$HIVE_NODETYPE" == "archive" ]; then
  FLAGS="$FLAGS --syncmode full --gcmode archive"
fi
if [ "$HIVE_NODETYPE" == "full" ]; then
  FLAGS="$FLAGS --syncmode full"
fi
if [ "$HIVE_NODETYPE" == "snap" ]; then
  FLAGS="$FLAGS --syncmode snap"
fi

# Don't immediately abort, some imports are meant to fail
set +e

echo "Running taiko-geth with flags $FLAGS"
geth $FLAGS
