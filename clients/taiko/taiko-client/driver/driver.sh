#!/bin/bash

source /driver/common.sh

# Set default values
DISABLE_P2P_SYNC=${DISABLE_P2P_SYNC:-true}

check_env "L1_NODE_WS_ENDPOINT"
check_env "L1_BEACON_HTTP"
check_env "L2_NODE_WS_ENDPOINT"
check_env "L2_EXECUTION_ENGINE_AUTH_ENDPOINT"
check_env "TAIKO_L1_ADDRESS"
check_env "TAIKO_L2_ADDRESS"

if [ "$DISABLE_P2P_SYNC" = "false" ]; then
  check_env "P2P_SYNC_URL"
fi

# Run driver
echo "Starting driver..."
if [ "$DISABLE_P2P_SYNC" = "false" ]; then
  exec taiko-client driver \
    --l1.ws "${L1_NODE_WS_ENDPOINT}" \
    --l2.ws "${L2_NODE_WS_ENDPOINT}" \
    --l1.beacon "${L1_BEACON_HTTP}" \
    --l2.auth "${L2_EXECUTION_ENGINE_AUTH_ENDPOINT}" \
    --taikoL1 "${TAIKO_L1_ADDRESS}" \
    --taikoL2 "${TAIKO_L2_ADDRESS}" \
    --jwtSecret /jwt.hex \
    --p2p.sync \
    --p2p.checkPointSyncUrl "${P2P_SYNC_URL}"
else
  exec taiko-client driver \
    --l1.ws "${L1_ENDPOINT_WS}" \
    --l2.ws "${L2_EXECUTION_ENGINE_WS_ENDPOINT}" \
    --l1.beacon "${L1_BEACON_HTTP}" \
    --l2.auth "${L2_EXECUTION_ENGINE_AUTH_ENDPOINT}" \
    --taikoL1 "${TAIKO_L1_ADDRESS}" \
    --taikoL2 "${TAIKO_L2_ADDRESS}" \
    --jwtSecret /jwt.hex
fi
