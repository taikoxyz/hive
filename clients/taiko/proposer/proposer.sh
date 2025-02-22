#!/bin/sh

. /taiko-client/common.sh

check_env "L1_WS"
check_env "L2_AUTH"
check_env "L2_HTTP"
check_env "TAIKO_INBOX"
check_env "TAIKO_ANCHOR"
check_env "TAIKO_TOKEN"
check_env "L1_PROPOSER_PRIV_KEY"
check_env "HIVE_TAIKO_FEE_RECEIPT"

if [ "$HIVE_TAIKO_FEE_RECEIPT" != "" ]; then
  export L2_SUGGESTED_FEE_RECIPIENT="$HIVE_TAIKO_FEE_RECEIPT"
fi

# Run proposer
echo "Starting proposer..."
nohup taiko-client proposer --verbosity 4 2>&1 &

touch /taiko-client/nohup.out
tail -f /taiko-client/nohup.out
