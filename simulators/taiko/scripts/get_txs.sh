#!/bin/bash

# get docker env
. scripts/docker_env.sh

END_BLOCK=$(cast block-number --rpc-url "$L1_PROBE_URL")

for ((block = 0; block <= "$END_BLOCK"; block++)); do
  echo "Checking block $block"
  # Get all the transactions.
  transactions=$(cast block $block --json --rpc-url "$L1_PROBE_URL" | jq -r '.transactions[]')

  rm -rf params/l1contract_txs.txt

  # Store txs.
  for tx in $transactions; do
    echo "tx hash: $tx"
    if [ "$tx" != "" ]; then
      cast tx "$tx" --json --rpc-url "$L1_PROBE_URL" >>params/l1contract_txs.txt
    fi
  done
done
