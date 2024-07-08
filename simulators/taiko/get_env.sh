#!/bin/bash

# get deployed contract address.
DEPLOYMENT_JSON=$(cat /taiko/taiko-mono/packages/protocol/deployments/deploy_l1.json)
export TAIKO_L1=$(echo "$DEPLOYMENT_JSON" | jq '.taiko' | sed 's/\"//g')
export TAIKO_TOKEN=$(echo "$DEPLOYMENT_JSON" | jq '.taiko_token' | sed 's/\"//g')
export PROVER_SET=$(echo "$DEPLOYMENT_JSON" | jq '.prover_set' | sed 's/\"//g')
export GUARDIAN_PROVER_MINORITY=$(echo "$DEPLOYMENT_JSON" | jq '.guardian_prover_minority' | sed 's/\"//g')

# show the integration test environment variables.
# L1_BEACON_HTTP_ENDPOINT=$L1_BEACON_HTTP_ENDPOINT
echo "L2_SUGGESTED_FEE_RECIPIENT=$L2_SUGGESTED_FEE_RECIPIENT
TAIKO_L1=$TAIKO_L1
TAIKO_L2=$TAIKO_L2
TAIKO_TOKEN=$TAIKO_TOKEN
PROVER_SET=$PROVER_SET
GUARDIAN_PROVER_MINORITY=$GUARDIAN_PROVER_MINORITY" >/taiko/.env

# show contract addresses.
cat /taiko/.env
