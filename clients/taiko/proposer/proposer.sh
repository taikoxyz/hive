#!/bin/sh

. /taiko-client/common.sh

# default env
export JWT_SECRET="/taiko-client/jwt.hex"

check_env "L1_WS"
check_env "L2_AUTH"
check_env "L2_HTTP"
check_env "TAIKO_L1"
check_env "TAIKO_L2"
check_env "TAIKO_TOKEN"
check_env "L1_PROPOSER_PRIV_KEY"
check_env "L2_SUGGESTED_FEE_RECIPIENT"

# Run proposer
echo "Starting proposer..."
exec taiko-client proposer
