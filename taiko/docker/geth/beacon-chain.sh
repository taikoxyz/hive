#!/bin/bash

CONSENSUS_DIR=/devnet/consensus

until ps aux | grep "geth " | grep -v grep 2>/dev/null; do
  sleep 1
done

echo "Starting beacon chain..."
beacon-chain \
  --datadir=$CONSENSUS_DIR/beacondata \
  --min-sync-peers=0 \
  --genesis-state=$CONSENSUS_DIR/genesis.ssz \
  --bootstrap-node= \
  --interop-eth1data-votes \
  --chain-config-file=$CONSENSUS_DIR/config.yml \
  --contract-deployment-block=0 \
  --chain-id=32382 \
  --rpc-host=0.0.0.0 \
  --grpc-gateway-host=0.0.0.0 \
  --execution-endpoint=http://localhost:8551 \
  --accept-terms-of-use \
  --jwt-secret=$CONSENSUS_DIR/jwtsecret \
  --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
  --minimum-peers-per-subnet=0 \
  --enable-debug-rpc-endpoints
