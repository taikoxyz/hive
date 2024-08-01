#!/bin/bash

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/prysm

mkdir -p $EXECUTION_DIR/data/vc
mkdir -p $EXECUTION_DIR/data/validators

ls $EXECUTION_DIR/data/validators

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

builder_option=$([[ "$HIVE_TAIKO2_BUILDER_ENDPOINT" == "" ]] && echo "" || echo "--enable-builder")
echo BUILDER=$builder_option

echo Starting Prysm Validator Client

/validator \
  --verbosity="$LOG" \
  --accept-terms-of-use=true \
  --enable-beacon-rest-api=true \
  --beacon-rest-api-provider="http://$HIVE_TAIKO2_BN_API_IP:${HIVE_TAIKO2_BN_API_PORT:-3500}" \
  --datadir="$EXECUTION_DIR/data/vc" \
  --chain-config-file="/hive/input/config.yaml" \
  --suggested-fee-recipient="${HIVE_TAIKO2_FEE_RECEIPT}" \
  --force-clear-db \
  $builder_option
