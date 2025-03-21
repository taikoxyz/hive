#!/bin/bash

# load tool commands.
source common.sh

# Immediately abort the script on any error encountered
set -e

check_env "HIVE_TAIKO_ETH1_RPC_ADDRS"

LOG=info
case "$HIVE_LOGLEVEL" in
0) LOG=fatal ;;
1) LOG=error ;;
2) LOG=warn ;;
3) LOG=info ;;
4) LOG=debug ;;
5) LOG=trace ;;
esac

beacon-chain \
  --chain-id=32382 \
  --verbosity=$LOG \
  --datadir=beacon \
  --genesis-state=genesis.ssz \
  --chain-config-file=config.yml \
  --min-sync-peers=0 \
  --interop-eth1data-votes \
  --contract-deployment-block=0 \
  --rpc-host=0.0.0.0 \
  --grpc-gateway-host=0.0.0.0 \
  --execution-endpoint=$HIVE_TAIKO_ETH1_RPC_ADDRS \
  --accept-terms-of-use \
  --jwt-secret=/tmp/jwt.hex \
  --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
  --minimum-peers-per-subnet=0 \
  --enable-debug-rpc-endpoints \
  --blob-retention-epochs=409600
