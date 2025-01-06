#!/bin/sh

. /taiko-client/common.sh

# taiko-client's necessary common envs.
check_env "L1_WS"
#check_env "L1_BEACON"
check_env "L2_WS"
check_env "L2_AUTH"
check_env "TAIKO_L1_ADDRESS"
check_env "TAIKO_L2_ADDRESS"
# check_env "SOFT_BLOCK_SERVER_PORT"

if [ "$P2P_SYNC" != "" ]; then
  check_env "P2P_CHECK_POINT_SYNC_URL"
fi

# Run driver
echo "Starting driver..."
nohup taiko-client driver \
  --verbosity 4 2>&1 &

touch /taiko-client/nohup.out
tail -f /taiko-client/nohup.out
