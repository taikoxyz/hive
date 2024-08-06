#!/bin/sh

. /taiko-client/common.sh

check_env "L1_WS"
check_env "L2_AUTH"
check_env "L2_HTTP"
check_env "TAIKO_L1"
check_env "TAIKO_L2"
check_env "TAIKO_TOKEN"
check_env "L1_PROPOSER_PRIV_KEY"
check_env "HIVE_TAIKO2_JWT_SECRET"
check_env "HIVE_TAIKO2_FEE_RECEIPT"

if [ "$HIVE_TAIKO2_JWT_SECRET" != "" ]; then
  echo "$HIVE_TAIKO2_JWT_SECRET" >/taiko-client/jwt.hex
  export JWT_SECRET="/taiko-client/jwt.hex"
fi

if [ "$HIVE_TAIKO2_FEE_RECEIPT" != "" ]; then
  export L2_SUGGESTED_FEE_RECIPIENT="$HIVE_TAIKO2_FEE_RECEIPT"
fi

# Run proposer
echo "Starting proposer..."
exec taiko-client proposer
