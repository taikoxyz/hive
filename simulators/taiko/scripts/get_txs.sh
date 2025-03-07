#!/bin/bash

params_path=params
rm -f $params_path/l1contract_txs.txt

# get docker env
. scripts/docker_env.sh

END_BLOCK=$(cast block-number --rpc-url "$L1_PROBE_URL")

for block in $(seq 0 "$END_BLOCK"); do
  echo "Checking block $block"
  # Get all the transactions.
  transactions=$(cast block "$block" --json --rpc-url "$L1_PROBE_URL" | jq -r '.transactions[]')

  # Store txs.
  for tx in $transactions; do
    echo "tx hash: $tx"
    if [ "$tx" != "" ]; then
      cast tx "$tx" --json --rpc-url "$L1_PROBE_URL" >>$params_path/l1contract_txs.txt
    fi
  done
done
