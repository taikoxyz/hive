#!/bin/bash

rm -rf params/l1contract_txs.txt

L1_PROBE_URL=${L1_PROBE_URL:-http://localhost:8545}
START_BLOCK=0
END_BLOCK=$(cast block-number --rpc-url "$L1_PROBE_URL")

for ((block = $START_BLOCK; block <= $END_BLOCK; block++)); do
  echo "Checking block $block"
  # Get all the transactions.
  transactions=$(cast block $block --json --rpc-url "$L1_PROBE_URL" | jq -r '.transactions[]')

  # Store txs.
  for tx in $transactions; do
    echo "tx hash: $tx"
    if [ "$tx" != "" ]; then
      cast tx "$tx" --json --rpc-url "$L1_PROBE_URL" >>params/l1contract_txs.txt
    fi
  done
done
