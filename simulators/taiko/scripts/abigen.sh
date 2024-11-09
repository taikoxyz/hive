#!/bin/sh

if [ ! -d "$TAIKO_MONO_DIR/packages/protocol/out" ]; then
    echo "ABI not generated in protocol package yet. Please run npm install && npx hardhat compile in ../protocol"
    exit 1
fi

paths=("TaikoToken.sol" "TaikoL1.sol" "TaikoL2.sol" "GuardianProver.sol" "ProverSet.sol" "LibProving.sol" "DevnetTierRouter.sol")

names=("TaikoToken" "TaikoL1" "TaikoL2" "GuardianProver" "ProverSet" "LibProving" "DevnetTierRouter")

for (( i = 0; i < ${#paths[@]}; ++i ));
do
    lower=$(echo "${names[i]}" | tr '[:upper:]' '[:lower:]')
    mkdir -p bindings/$lower
    jq .abi "$TAIKO_MONO_DIR"/packages/protocol/out/${paths[i]}/${names[i]}.json > bindings/$lower/${names[i]}.json
    abigen --abi bindings/$lower/${names[i]}.json \
    --pkg $lower \
    --type ${names[i]} \
    --out bindings/$lower/${names[i]}.go
    rm bindings/$lower/${names[i]}.json
done

exit 0
