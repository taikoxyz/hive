#!/bin/bash

# load tool commands.
source /prysm/common.sh

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/prysm

check_env "HIVE_TAIKO_FEE_RECEIPT"
check_env "HIVE_TAIKO_BN_API_IP"

mkdir -p $EXECUTION_DIR/data/vc

echo config.yaml:
cat /hive/input/config.yaml

LOG=info
case "$HIVE_LOGLEVEL" in
0) LOG=fatal ;;
1) LOG=error ;;
2) LOG=warn ;;
3) LOG=info ;;
4) LOG=debug ;;
5) LOG=trace ;;
esac

builder_option=$([[ "$HIVE_TAIKO_BUILDER_ENDPOINT" == "" ]] && echo "" || echo "--enable-builder")
echo BUILDER=$builder_option

echo Starting Prysm Validator Client

validator \
  --verbosity="$LOG" \
  --chain-config-file="/hive/input/config.yaml" \
  --beacon-rpc-provider="$HIVE_TAIKO_BN_API_IP:4000" \
  --accept-terms-of-use=true \
  --interop-start-index=0 \
  --interop-num-validators="${HIVE_TAIKO_NUM_VALIDATORS:-64}" \
  --datadir="$EXECUTION_DIR/data/vc" \
  --force-clear-db \
  $builder_option
