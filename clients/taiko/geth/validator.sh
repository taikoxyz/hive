#!/bin/bash

CONSENSUS_DIR=/devnet/consensus

until ps aux | grep "beacon-chain " | grep -v grep 2>/dev/null; do
  sleep 1
done

ehco "Starting validator..."
validator \
  --beacon-rpc-provider=localhost:4000 \
  --datadir=$CONSENSUS_DIR/validatordata \
  --accept-terms-of-use \
  --interop-num-validators=3 \
  --interop-start-index=0 \
  --chain-config-file=$CONSENSUS_DIR/config.yml
