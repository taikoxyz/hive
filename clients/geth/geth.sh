#!/bin/bash

HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}

# load tool commands.
. common.sh

# Immediately abort the script on any error encountered
set -e

# Reset genesis.json
prysmctl \
  testnet \
  generate-genesis \
  --fork=fulu \
  --num-validators=1 \
  --genesis-time="1742366600" \
  --output-ssz=genesis.ssz \
  --chain-config-file=config.yml \
  --geth-genesis-json-in=genesis.json \
  --geth-genesis-json-out=genesis.json

echo genesis.json:
cat genesis.json

# Initialize the local testchain with the genesis state
echo "Initializing database with genesis state..."
geth init --datadir=data --state.scheme=hash genesis.json

# Move keystore file into data/keystore.
mv keyfile.json data/keystore

echo "Starting geth..."
geth \
  --verbosity "$HIVE_LOGLEVEL" \
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
  --unlock=0x123463a4b065722e99115d6c222f267d9cabb524 \
  --password=geth_password.txt \
  --nodiscover \
  --gcmode=archive \
  --gcmode=full
