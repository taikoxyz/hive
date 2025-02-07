#!/bin/bash

TAIKO_MONO_DIR=${TAIKO_MONO_DIR:-$HOME/projects/taiko/taiko-mono}

if [ ! -d "$TAIKO_MONO_DIR/packages/protocol/out" ]; then
    echo "ABI not generated in protocol package yet. Please run npm install && npx hardhat compile in ../protocol"
    exit 1
fi

paths=(
    "TaikoInbox"
    "TaikoAnchor"
    "TaikoToken"
    "ResolverBase"
    "ProverSet"
    "ForkRouter"
    "ComposeVerifier"
)

bindings_path=bindings/pacaya

for (( i = 0; i < ${#paths[@]}; ++i ));
do
    mkdir -p $bindings_path/$lower
    latest_name=$(basename ${paths[i]})
    lower=$(echo "${latest_name}" | tr '[:upper:]' '[:lower:]')
    jq .abi "$TAIKO_MONO_DIR"/packages/protocol/out/${paths[i]}.sol/${latest_name}.json > $bindings_path/$lower/${latest_name}.json
    abigen --abi $bindings_path/$lower/${latest_name}.json \
    --pkg $lower \
    --type ${latest_name} \
    --out $bindings_path/$lower/${latest_name}.go
    rm $bindings_path/$lower/${latest_name}.json
done

exit 0
