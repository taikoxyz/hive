#!/bin/sh

HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}

. /anvil/common.sh

CONTAINER_IP=$(hostname -i | awk '{print $1;}')
echo "anvil CONTAINER_IP: $CONTAINER_IP"

echo genesis.json:
cat genesis.json

anvil --init genesis.json --host 0.0.0.0 --accounts 20
