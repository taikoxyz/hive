#!/bin/bash

OLD_FORK_TAIKO_MONO=${OLD_FORK_TAIKO_MONO:-$HOME/projects/taiko/tmp/taiko-mono}

if [ ! -d "$OLD_FORK_TAIKO_MONO/packages/protocol/out" ]; then
    echo "ABI not generated in protocol package yet. Please run npm install && npx hardhat compile in ../protocol"
    exit 1
fi

paths=(
    "TaikoL1"
    "LibProving"
    "LibProposing"
    "LibUtils"
    "LibVerifying"
    "TaikoL2"
    "TaikoToken"
    "AddressManager"
    "GuardianProver"
    "ProverSet"
    "MainnetTierRouter"
    "SgxVerifier"
)

bindings_path=bindings/ontake

for (( i = 0; i < ${#paths[@]}; ++i ));
do
    mkdir -p $bindings_path/$lower
    latest_name=$(basename ${paths[i]})
    lower=$(echo "${latest_name}" | tr '[:upper:]' '[:lower:]')
    jq .abi $OLD_FORK_TAIKO_MONO/packages/protocol/out/${paths[i]}.sol/${latest_name}.json > $bindings_path/$lower/${latest_name}.json
    abigen --abi $bindings_path/$lower/${latest_name}.json \
    --pkg $lower \
    --type ${latest_name} \
    --out $bindings_path/$lower/${latest_name}.go
    rm $bindings_path/$lower/${latest_name}.json
done

exit 0
