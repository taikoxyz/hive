#!/bin/bash

source /taiko-client/common.sh

# default env
export JWT_SECRET="/taiko-client/nodes/jwt.hex"

# taiko-client's necessary common envs.
check_env "L1_WS"
check_env "L1_BEACON"
check_env "L2_WS"
check_env "L2_AUTH"
check_env "TAIKO_L1"
check_env "TAIKO_L2"

if [ "$P2P_SYNC" = "true" ]; then
  check_env "P2P_CHECK_POINT_SYNC_URL"
fi

# Run driver
echo "Starting driver..."
exec taiko-client driver
