#!/bin/bash

# load tool commands.
source /prysm/common.sh

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/prysm

check_env "HIVE_TAIKO2_FEE_RECEIPT"
check_env "HIVE_TAIKO2_BN_API_IP"

mkdir -p $EXECUTION_DIR/data/vc
mkdir -p $EXECUTION_DIR/data/validators

#echo walletpassword >$EXECUTION_DIR/wallet.pass
#for keystore_path in /hive/input/keystores/*; do
#  pubkey=$(basename "$keystore_path")
#
#  validator accounts import \
#    --accept-terms-of-use=true \
#    --wallet-dir="$EXECUTION_DIR/data/validators" \
#    --keys-dir="/hive/input/keystores/$pubkey" \
#    --account-password-file="/hive/input/secrets/$pubkey" \
#    --wallet-password-file="$EXECUTION_DIR/wallet.pass"
#done
#ls $EXECUTION_DIR/data/validators

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

validator \
  --verbosity="$LOG" \
  --chain-config-file="/hive/input/config.yaml" \
  --beacon-rpc-provider="$HIVE_TAIKO2_BN_API_IP:4000" \
  --accept-terms-of-use=true \
  --interop-start-index=0 \
  --interop-num-validators="${HIVE_TAIKO2_NUM_VALIDATORS:-64}" \
  --datadir="$EXECUTION_DIR/data/vc" \
  --force-clear-db \
  $builder_option
