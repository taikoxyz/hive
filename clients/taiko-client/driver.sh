#!/bin/bash

source /common.sh

# Set default values
DISABLE_P2P_SYNC=${DISABLE_P2P_SYNC:-false}
P2P_SYNC_URL=${P2P_SYNC_URL:-https://rpc.mainnet.taiko.xyz}
TAIKO_L1_ADDRESS=${TAIKO_L1_ADDRESS:-0x06a9Ab27c7e2255df1815E6CC0168d7755Feb19a}
TAIKO_L2_ADDRESS=${TAIKO_L2_ADDRESS:-0x1670000000000000000000000000000000010001}

check_env "L1_ENDPOINT_WS"
check_env "L1_BEACON_HTTP"
check_env "L2_EXECUTION_ENGINE"

# Run driver
echo "Starting driver..."
if [ "$DISABLE_P2P_SYNC" = "false" ]; then
    exec taiko-client driver \
        --l1.ws "${L1_ENDPOINT_WS}" \
        --l2.ws ws://"${L2_EXECUTION_ENGINE}":8546 \
        --l1.beacon "${L1_BEACON_HTTP}" \
        --l2.auth http://"${L2_EXECUTION_ENGINE}":8551 \
        --taikoL1 "${TAIKO_L1_ADDRESS}" \
        --taikoL2 "${TAIKO_L2_ADDRESS}" \
        --jwtSecret /jwt.hex \
        --p2p.sync \
        --p2p.checkPointSyncUrl "${P2P_SYNC_URL}"
else
    exec taiko-client driver \
        --l1.ws "${L1_ENDPOINT_WS}" \
        --l2.ws ws://"${L2_EXECUTION_ENGINE}":8546 \
        --l1.beacon "${L1_BEACON_HTTP}" \
        --l2.auth http://"${L2_EXECUTION_ENGINE}":8551 \
        --taikoL1 "${TAIKO_L1_ADDRESS}" \
        --taikoL2 "${TAIKO_L2_ADDRESS}" \
        --jwtSecret /jwt.hex
fi