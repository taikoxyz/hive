#!/bin/bash

# load tool commands.
. /geth/common.sh

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/geth

check_env "HIVE_TAIKO_CLIQUE_PRIVATEKEY"
check_env "HIVE_TAIKO_CLIQUE_ADDRESS"

echo genesis.json:
cat /hive/input/genesis.json

# Initialize the local testchain with the genesis state
echo "Initializing database with genesis state..."
geth init /hive/input/genesis.json

if [ "$HIVE_TAIKO_CLIQUE_PRIVATEKEY" != "" ]; then
  echo "Importing clique key..."
  echo "secret" >$EXECUTION_DIR/geth_password.txt
  geth account import --password $EXECUTION_DIR/geth_password.txt <(echo "$HIVE_TAIKO_CLIQUE_PRIVATEKEY")
else
  echo "clique private key is not set, exiting..."
  exit 1
fi

echo "Starting geth..."
geth \
  --bootnodes="$HIVE_BOOTNODE" \
  --http \
  --http.api=admin,debug,eth,net,web3,txpool,miner \
  --http.addr=0.0.0.0 \
  --http.vhosts=* \
  --http.corsdomain=* \
  --ws \
  --ws.api=admin,debug,eth,net,web3,txpool,miner \
  --ws.addr=0.0.0.0 \
  --ws.origins=* \
  --authrpc.vhosts=* \
  --authrpc.addr=0.0.0.0 \
  --authrpc.jwtsecret=/tmp/jwt.hex \
  --allow-insecure-unlock \
  --unlock="$HIVE_TAIKO_CLIQUE_ADDRESS" \
  --password=$EXECUTION_DIR/geth_password.txt \
  --nodiscover \
  --gcmode=archive \
  --state.scheme=path \
  --gcmode=full
