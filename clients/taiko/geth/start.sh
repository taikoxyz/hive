#!/bin/bash

DATA_DIR=/devnet

if [ ! -f "$DATA_DIR/GENESIS_LOCK" ]; then

  # Remove old data if exists.
  rm -rf DATA_DIR=/devnet/execution/geth
  rm -rf DATA_DIR=/devnet/consensus/beacondata
  rm -rf DATA_DIR=/devnet/consensus/vaidatordata

  # Add lock file.
  touch "$DATA_DIR/GENESIS_LOCK"

  # Reset genesis.json
  prysmctl \
    testnet \
    generate-genesis \
    --fork=deneb \
    --num-validators=1 \
    --genesis-time-delay=15 \
    --output-ssz=$DATA_DIR/consensus/genesis.ssz \
    --chain-config-file=$DATA_DIR/consensus/config.yml \
    --geth-genesis-json-in=$DATA_DIR/execution/genesis.json \
    --geth-genesis-json-out=$DATA_DIR/execution/genesis.json

  cat $DATA_DIR/execution/genesis.json

  # Init geth.
  rm -rf $DATA_DIR/execution/geth
  geth init --datadir=$DATA_DIR/execution $DATA_DIR/execution/genesis.json
fi

# Run beacon-chain service.
nohup sh /devnet/beacon-chain.sh >/devnet/beacon-chain.log 2>&1 &

# Run validator service.
nohup sh /devnet/validator.sh >/devnet/validator.log 2>&1 &

sh /devnet/geth.sh
