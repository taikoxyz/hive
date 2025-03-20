#!/bin/sh

HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}

. /anvil/common.sh

CONTAINER_IP=$(hostname -i | awk '{print $1;}')
echo "anvil CONTAINER_IP: $CONTAINER_IP"

if [ -f /hive/input/genesis.json ]; then
  echo genesis.json:
  cat /hive/input/genesis.json
  anvil --init /hive/input/genesis.json --host 0.0.0.0 --accounts 20
else
  anvil --host 0.0.0.0 --accounts 20
fi
