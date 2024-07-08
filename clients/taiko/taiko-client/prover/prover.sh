#!/bin/bash

source /taiko-client/common.sh

# default env
export JWT_SECRET="/taiko-client/nodes/jwt.hex"

check_env "L1_HTTP"
check_env "L1_WS"
check_env "L2_HTTP"
check_env "L2_WS"
check_env "PROVER_SET"
check_env "TAIKO_L1"
check_env "TAIKO_L2"
check_env "TAIKO_TOKEN"
check_env "L1_PROVER_PRIV_KEY"

# Run prover
echo "Starting prover..."
exec taiko-client prover
