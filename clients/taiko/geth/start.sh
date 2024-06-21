#!/bin/bash

DATA_DIR=/l1-geth-data

if [ ! -f "$DATA_DIR/GENESIS_LOCK" ];then
  # Add lock file.
  mkdir -p $DATA_DIR
  touch "$DATA_DIR/GENESIS_LOCK"

  # Get current time.
  GENESIS_TIME=$(echo "$(date +%s) / 100 * 100" | bc)
  echo "CURRENT_TIME=$(date +%s)"
  echo "GENESIS_TIME=$GENESIS_TIME"

  # Reset genesis.json
  prysmctl \
    testnet \
    generate-genesis \
    --fork=deneb \
    --num-validators=64 \
    --genesis-time="$GENESIS_TIME" \
    --genesis-time-delay=100 \
    --output-ssz=genesis.ssz \
    --chain-config-file=config.yml \
    --geth-genesis-json-in=genesis.json \
    --geth-genesis-json-out=genesis.json

  cat genesis.json

  # Init geth.
  geth init --datadir=$DATA_DIR genesis.json

  # Move keystore file into data/keystore.
  mv keyfile.json $DATA_DIR/keystore
fi


# Run geth service.
export DEVNET=true && \
  nohup geth \
  --http \
  --http.api=debug,eth,net,web3,txpool,miner \
  --http.addr=0.0.0.0  \
  --http.vhosts=* \
  --http.corsdomain=* \
  --ws \
  --ws.api=debug,eth,net,web3,txpool,miner \
  --ws.addr=0.0.0.0 \
  --ws.origins=* \
  --authrpc.vhosts=* \
  --authrpc.addr=0.0.0.0 \
  --authrpc.jwtsecret=jwtsecret \
  --datadir=$DATA_DIR \
  --allow-insecure-unlock \
  --unlock=0x123463a4b065722e99115d6c222f267d9cabb524 \
  --password=geth_password.txt \
  --nodiscover \
  --gcmode=archive \
  --syncmode=full 2>&1 &

echo "Starting deploy l1 contract..."
sh /execution/deploy_l1_contract.sh

# wait until geth finished.
wait "$(ps aux | grep geth | grep -v grep | awk '{print $2}')"