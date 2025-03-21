#!/bin/bash

# load tool commands.
source common.sh

# Immediately abort the script on any error encountered
set -e

check_env "HIVE_TAIKO_BN_API_IP"

LOG=info
case "$HIVE_LOGLEVEL" in
0) LOG=fatal ;;
1) LOG=error ;;
2) LOG=warn ;;
3) LOG=info ;;
4) LOG=debug ;;
5) LOG=trace ;;
esac

echo Starting Prysm Validator Client
validator \
  --verbosity=warn \
  --beacon-rpc-provider="$HIVE_TAIKO_BN_API_IP:4000" \
  --datadir=/l1-validator \
  --accept-terms-of-use \
  --interop-num-validators=1 \
  --interop-start-index=0 \
  --chain-config-file=config.yml
