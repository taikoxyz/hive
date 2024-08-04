#!/bin/bash

EXECUTION_DIR=/devnet/execution

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
    --datadir=$EXECUTION_DIR \
    --allow-insecure-unlock \
    --unlock=0x123463a4b065722e99115d6c222f267d9cabb524 \
    --password=$EXECUTION_DIR/geth_password.txt \
    --nodiscover \
    --gcmode=archive \
    --syncmode=full
