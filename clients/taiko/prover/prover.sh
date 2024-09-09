#!/bin/sh

. /taiko-client/common.sh

check_env "L1_HTTP"
check_env "L1_WS"
check_env "L2_HTTP"
check_env "L2_WS"
check_env "TAIKO_L1_ADDRESS"
check_env "TAIKO_L2_ADDRESS"
check_env "TAIKO_TOKEN_ADDRESS"
check_env "L1_PROVER_PRIV_KEY"
#check_env "GUARDIAN_PROVER_MINORITY"
#check_env "GUARDIAN_PROVER_MAJORITY"
#check_env "GUARDIAN_PROVER_CONTRACT"
check_env "HIVE_TAIKO_JWT_SECRET"

if [ "$HIVE_TAIKO_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO_JWT_SECRET" >/taiko-client/jwt.hex
  export JWT_SECRET="/taiko-client/jwt.hex"
fi

# Run prover
echo "Starting prover..."
exec taiko-client prover
