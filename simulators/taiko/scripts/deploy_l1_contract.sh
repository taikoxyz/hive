#!/bin/bash

HIVE_TAIKO_PATH=$PWD
OLD_FORK_TAIKO_MONO=${OLD_FORK_TAIKO_MONO:-$HOME/projects/taiko/tmp/taiko-mono}
TAIKO_MONO_DIR=${TAIKO_MONO_DIR:-$HOME/projects/taiko/taiko-mono}

echo "HIVE_TAIKO_PATH: $HIVE_TAIKO_PATH"

# stop docker compose
docker compose -f $HIVE_TAIKO_PATH/docker/docker-compose.yml up l1_node l2_pacaya -d --wait
trap "docker compose -f $HIVE_TAIKO_PATH/docker/docker-compose.yml down" EXIT SIGINT SIGTERM ERR

# load l1 chain deploy contracts environment variables
source scripts/deploy_env.sh

# load docker environment variables
source scripts/docker_env.sh

echo "Start deploy l1 ontake contracts ..."
# Deploy v1.9.1 protocol at first
cd ${OLD_FORK_TAIKO_MONO}/packages/protocol &&
  forge script script/layer1/DeployProtocolOnL1.s.sol:DeployProtocolOnL1 \
    --fork-url "$L1_PROBE_URL" \
    --broadcast \
    --ffi \
    -vvvvv \
    --evm-version cancun \
    --private-key "$PRIVATE_KEY" \
    --block-gas-limit 200000000 \
    --legacy

cd - || exit

# Get deployed contract address.
DEPLOYMENT_JSON=$(cat ${OLD_FORK_TAIKO_MONO}/packages/protocol/deployments/deploy_l1.json)
export OLD_FORK=0x1291Be112d480055DaFd8a610b7d1e203891C274
export TAIKO_INBOX=$(echo "$DEPLOYMENT_JSON" | jq '.taiko' | sed 's/\"//g')
export ROLLUP_RESOLVER=$(echo "$DEPLOYMENT_JSON" | jq '.rollup_address_manager' | sed 's/\"//g')
export PROVER_SET=$(echo "$DEPLOYMENT_JSON" | jq '.prover_set' | sed 's/\"//g')
export TAIKO_TOKEN=$(echo "$DEPLOYMENT_JSON" | jq '.taiko_token' | sed 's/\"//g')
export SGX_VERIFIER=$(echo "$DEPLOYMENT_JSON" | jq '.tier_sgx' | sed 's/\"//g')
export RISC0_VERIFIER=$(echo "$DEPLOYMENT_JSON" | jq '.tier_zkvm_risc0' | sed 's/\"//g')
export SP1_VERIFIER=$(echo "$DEPLOYMENT_JSON" | jq '.tier_zkvm_sp1' | sed 's/\"//g')
export SHARED_RESOLVER=$(echo "$DEPLOYMENT_JSON" | jq '.shared_address_manager' | sed 's/\"//g')
export BRIDGE_L1=$(echo "$DEPLOYMENT_JSON" | jq '.bridge' | sed 's/\"//g')
export SIGNAL_SERVICE=$(echo "$DEPLOYMENT_JSON" | jq '.signal_service' | sed 's/\"//g')
export ERC20_VAULT=$(echo "$DEPLOYMENT_JSON" | jq '.erc20_vault' | sed 's/\"//g')
export ERC721_VAULT=$(echo "$DEPLOYMENT_JSON" | jq '.erc721_vault' | sed 's/\"//g')
export ERC1155_VAULT=$(echo "$DEPLOYMENT_JSON" | jq '.erc1155_vault' | sed 's/\"//g')
export QUOTA_MANAGER=0x0000000000000000000000000000000000000000

echo "Start upgrade l1 pacaya contracts ..."
cd ${TAIKO_MONO_DIR}/packages/protocol &&
  PRIVATE_KEY=$PRIVATE_KEY forge script script/layer1/devnet/UpgradeDevnetPacayaL1.s.sol:UpgradeDevnetPacayaL1 \
    --fork-url "$L1_PROBE_URL" \
    --broadcast \
    --ffi \
    -vvvvv \
    --evm-version cancun \
    --private-key "$PRIVATE_KEY" \
    --block-gas-limit 200000000 \
    --legacy

cd - || exit

#PACAYA_DEPLOYMENT_JSON=$(cat ${TAIKO_MONO_DIR}/packages/protocol/deployments/deploy_l1.json)
#export CONTRACT_OWNER=0x0000000000000000000000000000000000000000
#export TAIKO_WRAPPER=$(echo "$PACAYA_DEPLOYMENT_JSON" | jq '.taiko_wrapper' | sed 's/\"//g')
#
#cd ${TAIKO_MONO_DIR}/packages/protocol &&
#  PRIVATE_KEY=$PRIVATE_KEY forge script script/layer1/preconf/DeployPreconfContracts.s.sol:DeployPreconfContracts \
#    --fork-url "$L1_PROBE_URL" \
#    --broadcast \
#    --ffi \
#    -vvvvv \
#    --evm-version cancun \
#    --private-key "$PRIVATE_KEY" \
#    --block-gas-limit 200000000 \
#    --legacy
#
#cd - || exit

# Get envs
sh scripts/get_env.sh

# Get txs
sh scripts/get_txs.sh

# build ontake abigen
sh scripts/ontake_abigen.sh

# build pacaya abigen
sh scripts/pacaya_abigen.sh