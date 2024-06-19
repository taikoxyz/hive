#!/bin/sh

set -e

validator \
  --verbosity=debug \
  --beacon-rpc-provider=localhost:4000 \
  --datadir=/l1-validator \
  --accept-terms-of-use \
  --interop-num-validators=64 \
  --interop-start-index=0 \
  --chain-config-file=/config.yml
