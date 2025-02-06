#!/bin/bash

cd $TAIKO_MONO_DIR && git checkout const_contracts_ontake && cd -

. scripts/docker_env.sh
. scripts/l1_env.sh

echo "Start deploying taiko contracts on l1 chain..."
cd "$TAIKO_MONO_DIR"/packages/protocol && forge script script/layer1/DeployProtocolOnL1.s.sol:DeployProtocolOnL1 \
  --fork-url "$L1_PROBE_URL" \
  --broadcast \
  --ffi \
  -vvvv \
  --evm-version cancun \
  --private-key $PRIVATE_KEY \
  --block-gas-limit 200000000

cd - || exit
