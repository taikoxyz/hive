#!/bin/bash

# Immediately abort the script on any error encountered
set -e

EXECUTION_DIR=/prysm

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

echo "bootnodes: ${HIVE_TAIKO2_BOOTNODE_ENRS}"

echo config.yaml:
cat /hive/input/config.yaml

CONTAINER_IP=$(hostname -i | awk '{print $1;}')
metrics_option=$([[ "$HIVE_TAIKO2_METRICS_PORT" == "" ]] && echo "--disable-monitoring=true" || echo "--disable-monitoring=false --monitoring-host=0.0.0.0 --monitoring-port=$HIVE_TAIKO2_METRICS_PORT")

if [ -z "$HIVE_TAIKO2_JWT_SECRET" ]; then
  echo "$HIVE_TAIKO2_JWT_SECRET" >$EXECUTION_DIR/jwtsecret
fi

if [[ "$HIVE_TAIKO2_BOOTNODE_ENRS" == "" ]]; then
  bootnode_option=""
else
  bootnode_option=""
  for bn in ${HIVE_TAIKO2_BOOTNODE_ENRS//,/ }; do
    trimmed_bn=${bn//=/}
    bootnode_option="$bootnode_option --bootstrap-node=$trimmed_bn"
  done
fi
builder_option=$([[ "$HIVE_TAIKO2_BUILDER_ENDPOINT" == "" ]] && echo "" || echo "--http-mev-relay=$HIVE_TAIKO2_BUILDER_ENDPOINT")
echo BUILDER=$builder_option

echo Starting Prysm Beacon Node

/beacon-chain \
  --verbosity="$LOG" \
  --accept-terms-of-use=true \
  --datadir=/data/beacon \
  --chain-config-file=/hive/input/config.yaml \
  --genesis-state=/hive/input/genesis.ssz \
  $bootnode_option \
  --p2p-tcp-port="${HIVE_TAIKO2_P2P_TCP_PORT:-13000}" \
  --p2p-udp-port="${HIVE_TAIKO2_P2P_UDP_PORT:-12000}" \
  --p2p-host-ip="${CONTAINER_IP}" \
  --p2p-local-ip="${CONTAINER_IP}" \
  --execution-endpoint="$HIVE_TAIKO2_ETH1_RPC_ADDRS" \
  --jwt-secret=$EXECUTION_DIR/jwtsecret \
  --min-sync-peers=1 \
  --subscribe-all-subnets=true \
  $metrics_option \
  $builder_option \
  --deposit-contract="${HIVE_TAIKO2_DEPOSIT_CONTRACT_ADDRESS:-0x1111111111111111111111111111111111111111}" \
  --contract-deployment-block="${HIVE_TAIKO2_DEPOSIT_DEPLOY_BLOCK_NUMBER:-0}" \
  --grpc-gateway-host=0.0.0.0 \
  --rpc-host=0.0.0.0 \
  --rpc-port="${HIVE_TAIKO2_BN_GRPC_PORT:-4000}" \
  --grpc-gateway-host=0.0.0.0 --grpc-gateway-port="${HIVE_TAIKO2_BN_API_PORT:-3500}" --grpc-gateway-corsdomain="*" \
  --suggested-fee-recipient="${HIVE_TAIKO2_FEE_RECEIPT}" \
  --force-clear-db
# NOTE: gRPC/RPC ports are inverted to allow the simulator to access the REST API
