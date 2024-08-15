#!/bin/sh

. /anvil/common.sh

CONTAINER_IP=$(hostname -i | awk '{print $1;}')
echo "anvil CONTAINER_IP: $CONTAINER_IP"

echo genesis.json:
cat /hive/input/genesis.json

# Generate the version.txt file.
RUN anvil --version >/version.txt

anvil \
  --init /hive/input/genesis.json \
  --host 0.0.0.0 \
  --accounts 20
#--block-time "$HIVE_TAIKO2_BLOCK_TIME"
