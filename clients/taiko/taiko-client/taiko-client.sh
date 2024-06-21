#!/bin/bash

# - HIVE_TAIKO_CLIENT_ROLE=driver

# Set default values
HIVE_TAIKO_CLIENT_ROLE={$HIVE_TAIKO_CLIENT_ROLE=-driver}

set -e

if [ "$HIVE_TAIKO_CLIENT_ROLE" == "driver" ];then
  echo "starting driver..."
  sh /driver.sh
elif [ "$HIVE_TAIKO_CLIENT_ROLE" == "proposer" ]; then
  echo "starting proposer..."
  sh /proposer.sh
elif [ "$HIVE_TAIKO_CLIENT_ROLE" == "prover" ]; then
  echo "starting prover..."
  sh /prover.sh
else
  echo "unexpect taiko-client role: $HIVE_TAIKO_CLIENT_ROLE"
  exit 1
fi