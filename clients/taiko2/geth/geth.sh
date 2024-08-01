#!/bin/bash

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/geth

# Initialize the local testchain with the genesis state
echo "Initializing database with genesis state..."
geth --datadir $EXECUTION_DIR/data init /hive/input/genesis.json

if [ -z "$HIVE_TAIKO2_JWT_SECRET" ]; then
  echo "$HIVE_TAIKO2_JWT_SECRET" >$EXECUTION_DIR/jwtsecret
fi

if [ -z "$HIVE_TAIKO2_CLIQUE_PRIVATEKEY" ]; then
  echo "Importing clique key..."
  echo "secret" >$EXECUTION_DIR/geth_password.txt
  geth --datadir $EXECUTION_DIR/data account import --password $EXECUTION_DIR/geth_password.txt <(echo "$HIVE_TAIKO2_CLIQUE_PRIVATEKEY")
else
  echo "clique private key is not set, exiting..."
  exit 1
fi

echo "Starting geth..."
export DEVNET=true &&
  geth \
    --http \
    --http.api=debug,eth,net,web3,txpool,miner \
    --http.addr=0.0.0.0 \
    --http.vhosts=* \
    --http.corsdomain=* \
    --ws \
    --ws.api=debug,eth,net,web3,txpool,miner \
    --ws.addr=0.0.0.0 \
    --ws.origins=* \
    --authrpc.vhosts=* \
    --authrpc.addr=0.0.0.0 \
    --authrpc.jwtsecret=$EXECUTION_DIR/jwtsecret \
    --datadir=$EXECUTION_DIR/data \
    --allow-insecure-unlock \
    --unlock="$HIVE_TAIKO2_CLIQUE_ADDRESS" \
    --password=$EXECUTION_DIR/geth_password.txt \
    --nodiscover \
    --gcmode=archive \
    --syncmode=full
