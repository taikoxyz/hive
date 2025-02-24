#!/bin/bash

# Set default values
HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-4}
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

FLAGS="--syncmode=snap"
# Handle any client mode or operation requests
if [ "$HIVE_NODETYPE" == "archive" ]; then
  FLAGS="--syncmode=full --gcmode=archive"
fi
if [ "$HIVE_NODETYPE" == "full" ]; then
  FLAGS="--syncmode=full"
fi
if [ "$HIVE_NODETYPE" == "light" ]; then
  FLAGS="--syncmode=light"
fi
if [ "$HIVE_NODETYPE" == "snap" ]; then
  FLAGS="--syncmode=snap"
fi
echo FLAGS: "$FLAGS"

echo "Starting taiko geth..."
nohup geth \
  --bootnodes="$HIVE_BOOTNODE" \
  --taiko \
  --verbosity 3 \
  --datadir=/geth/data \
  --networkid="$HIVE_NETWORK_ID" \
  --http \
  --http.api=admin,debug,eth,net,web3,txpool,miner,taiko \
  --http.addr=0.0.0.0 \
  --http.vhosts=* \
  --http.corsdomain=* \
  --ws \
  --ws.api=admin,debug,eth,net,web3,txpool,miner,taiko \
  --ws.addr=0.0.0.0 \
  --ws.origins=* \
  --authrpc.vhosts=* \
  --authrpc.addr=0.0.0.0 \
  --authrpc.jwtsecret=/tmp/jwt.hex \
  --allow-insecure-unlock \
  --state.scheme=path \
  "$FLAGS" 2>&1 &

touch /geth/nohup.out
tail -f /geth/nohup.out
