#!/bin/bash

PRIVATE_KEY=${PRIVATE_KEY:-0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}

forge script script/DeployOnL1.s.sol:DeployOnL1 \
  --fork-url "http://localhost:8545" \
  --broadcast \
  --ffi \
  -vvvvv \
  --evm-version cancun \
  --private-key "$PRIVATE_KEY" \
  --block-gas-limit 100000000 \
  --legacy