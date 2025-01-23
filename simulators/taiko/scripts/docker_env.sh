#!/bin/bash

# check until L1 chain is ready
L1_PROBE_URL=http://localhost:$(docker port l1_node | grep '0.0.0.0' | awk -F ':' '{print $2}')
until cast chain-id --rpc-url "$L1_PROBE_URL" 2>/dev/null; do
  sleep 1
done

# check until L2 chain is ready
L2_PROBE_URL=http://localhost:$(docker port l2_$TAIKO_VERSION | grep "0.0.0.0" | awk -F ':' 'NR==1 {print $2}')
until cast chain-id --rpc-url "$L2_PROBE_URL" 2>/dev/null; do
  sleep 1
done

export L1_PROBE_URL
export L2_PROBE_URL

echo "L1_PROBE_URL: $L1_PROBE_URL"
echo "L2_PROBE_URL: $L2_PROBE_URL"
