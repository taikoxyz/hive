#!/bin/bash

rm params/l1contract_txs.txt /dev/null 2>&1

L1_NODE_HTTP_ENDPOINT=${L1_NODE_HTTP_ENDPOINT:-http://localhost:8545}
START_BLOCK=0
END_BLOCK=$(cast block-number --rpc-url "$L1_NODE_HTTP_ENDPOINT")

for ((block = $START_BLOCK; block <= $END_BLOCK; block++)); do
  echo "Checking block $block"
  # Get all the transactions.
  transactions=$(cast block $block --json --rpc-url "$L1_NODE_HTTP_ENDPOINT" | jq -r '.transactions[]')

  # Store txs.
  for tx in $transactions; do
    echo "tx hash: $tx"
    if [ "$tx" != "" ]; then
      cast tx "$tx" --json --rpc-url "$L1_NODE_HTTP_ENDPOINT" >>params/l1contract_txs.txt
    fi
  done
done
