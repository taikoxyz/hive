#!/bin/bash

# load tool commands.
source /prysm/common.sh

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/prysm

check_env "HIVE_TAIKO_JWT_SECRET"
check_env "HIVE_TAIKO_FEE_RECEIPT"
check_env "HIVE_TAIKO_DEPOSIT_CONTRACT_ADDRESS"
check_env "HIVE_TAIKO_CHAIN_ID"
check_env "HIVE_TAIKO_ETH1_RPC_ADDRS"

if [ ! -f "/hive/input/genesis.ssz" ]; then
  echo "genesis.ssz file is missing, exiting..."
  exit 1
fi

mkdir -p $EXECUTION_DIR/data/beacon
mkdir -p $EXECUTION_DIR/data/network

LOG=info
case "$HIVE_LOGLEVEL" in
0) LOG=fatal ;;
1) LOG=error ;;
2) LOG=warn ;;
3) LOG=info ;;
4) LOG=debug ;;
5) LOG=trace ;;
esac

echo "bootnodes: ${HIVE_TAIKO_BOOTNODE_ENRS}"

echo config.yaml:
cat /hive/input/config.yaml

CONTAINER_IP=$(hostname -i | awk '{print $1;}')
metrics_option=$([[ "$HIVE_TAIKO_METRICS_PORT" == "" ]] && echo "--disable-monitoring=true" || echo "--disable-monitoring=false --monitoring-host=0.0.0.0 --monitoring-port=$HIVE_TAIKO_METRICS_PORT")

if [ "$HIVE_TAIKO_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO_JWT_SECRET" >$EXECUTION_DIR/jwtsecret
fi

if [[ "$HIVE_TAIKO_BOOTNODE_ENRS" == "" ]]; then
  bootnode_option=""
else
  bootnode_option=""
  for bn in ${HIVE_TAIKO_BOOTNODE_ENRS//,/ }; do
    trimmed_bn=${bn//=/}
    bootnode_option="$bootnode_option --bootstrap-node=$trimmed_bn"
  done
fi
builder_option=$([[ "$HIVE_TAIKO_BUILDER_ENDPOINT" == "" ]] && echo "" || echo "--http-mev-relay=$HIVE_TAIKO_BUILDER_ENDPOINT")
echo BUILDER=$builder_option

echo Starting Prysm Beacon Node

beacon-chain \
  --verbosity="$LOG" \
  --chain-id="${HIVE_TAIKO_CHAIN_ID:-7}" \
  --datadir=/data/beacon \
  --chain-config-file=/hive/input/config.yaml \
  --genesis-state=/hive/input/genesis.ssz \
  --interop-eth1data-votes=true \
  --accept-terms-of-use=true \
  $bootnode_option \
  --execution-endpoint="$HIVE_TAIKO_ETH1_RPC_ADDRS" \
  --jwt-secret=$EXECUTION_DIR/jwtsecret \
  --min-sync-peers=0 \
  --subscribe-all-subnets=true \
  $metrics_option \
  $builder_option \
  --deposit-contract="${HIVE_TAIKO_DEPOSIT_CONTRACT_ADDRESS:-0x1111111111111111111111111111111111111111}" \
  --contract-deployment-block="${HIVE_TAIKO_DEPOSIT_DEPLOY_BLOCK_NUMBER:-0}" \
  --rpc-host="${CONTAINER_IP}" \
  --grpc-gateway-host=0.0.0.0 \
  --suggested-fee-recipient="${HIVE_TAIKO_FEE_RECEIPT}" \
  --force-clear-db
# NOTE: gRPC/RPC ports are inverted to allow the simulator to access the REST API
