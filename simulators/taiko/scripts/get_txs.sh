#!/bin/bash

rm l1contract_txs.txt

L1_NODE_HTTP_ENDPOINT=${L1_NODE_HTTP_ENDPOINT:-http://localhost:8545}
START_BLOCK=0
END_BLOCK=`cast block-number --rpc-url "$L1_NODE_HTTP_ENDPOINT"`

for ((block=$START_BLOCK; block<=$END_BLOCK; block++))
do
  echo "Checking block $block"
  # 获取区块详细信息并提取交易列表
  transactions=$(cast block $block --json --rpc-url $L1_NODE_HTTP_ENDPOINT | jq -r '.transactions[]')

  # 遍历每个交易
  for tx in $transactions
  do
    echo "tx hash: $tx"
    if [ "$tx" != "" ]; then
      cast tx "$tx" --json --rpc-url $L1_NODE_HTTP_ENDPOINT >> l1contract_txs.txt
    fi
  done
done
