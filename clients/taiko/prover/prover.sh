#!/bin/sh

HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}

. /taiko-client/common.sh

check_env "L1_HTTP"
check_env "L1_WS"
check_env "L2_HTTP"
check_env "L2_WS"
check_env "TAIKO_INBOX"
check_env "TAIKO_ANCHOR"
check_env "TAIKO_TOKEN"
check_env "L1_PROVER_PRIV_KEY"
#check_env "GUARDIAN_PROVER_MINORITY"
#check_env "GUARDIAN_PROVER_MAJORITY"
#check_env "GUARDIAN_PROVER_CONTRACT"

# Run prover
echo "Starting prover..."
nohup taiko-client prover --verbosity $HIVE_LOGLEVEL 2>&1 &

touch /taiko-client/nohup.out
tail -f /taiko-client/nohup.out
