#!/bin/sh

. /taiko-client/common.sh

export PROVER_DUMMY=true
export PROVER_CAPACITY=1024
export PROVER_L1_NODE_VERSION=1.0.0
export PROVER_L2_NODE_VERSION=0.1.0
export PROVER_ALLOWANCE=10.0

check_env "L1_HTTP"
check_env "L1_WS"
check_env "L2_HTTP"
check_env "L2_WS"
check_env "TAIKO_L1"
check_env "TAIKO_L2"
check_env "TAIKO_TOKEN"
check_env "L1_PROVER_PRIV_KEY"
check_env "GUARDIAN_PROVER_MINORITY"
check_env "GUARDIAN_PROVER_MAJORITY"
check_env "GUARDIAN_PROVER_CONTRACT"
check_env "HIVE_TAIKO2_JWT_SECRET"

if [ "$HIVE_TAIKO2_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO2_JWT_SECRET" >/taiko-client/jwt.hex
  export JWT_SECRET="/taiko-client/jwt.hex"
fi

# Run prover
echo "Starting prover..."
exec taiko-client prover
