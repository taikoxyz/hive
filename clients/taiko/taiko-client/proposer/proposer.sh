#!/bin/bash

source /proposer/common.sh

# Set default values
BLOCK_PROPOSAL_FEE=${BLOCK_PROPOSAL_FEE:-1}

check_env "L1_NODE_WS_ENDPOINT"
check_env "L1_PROPOSER_PRIVATE_KEY"
check_env "L2_NODE_HTTP_ENDPOINT"
check_env "L2_EXECUTION_ENGINE_AUTH_ENDPOINT"
check_env "TAIKO_L1_ADDRESS"
check_env "TAIKO_TOKEN_L1_ADDRESS"
check_env "ASSIGNMENT_HOOK_L1_ADDRESS"
check_env "PROVER_ENDPOINTS"
check_env "TAIKO_L2_ADDRESS"
check_env "L2_SUGGESTED_FEE_RECIPIENT"

ARGS="--l1.ws ${L1_NODE_WS_ENDPOINT}
  --l2.http ${L2_NODE_HTTP_ENDPOINT}
  --l2.auth ${L2_EXECUTION_ENGINE_AUTH_ENDPOINT}
  --taikoL1 ${TAIKO_L1_ADDRESS}
  --taikoL2 ${TAIKO_L2_ADDRESS}
  --taikoToken ${TAIKO_TOKEN_L1_ADDRESS}
  --assignmentHookAddress ${ASSIGNMENT_HOOK_L1_ADDRESS}
  --jwtSecret /proposer/jwt.hex
  --l1.proposerPrivKey ${L1_PROPOSER_PRIVATE_KEY}
  --l2.suggestedFeeRecipient ${L2_SUGGESTED_FEE_RECIPIENT}
  --proverEndpoints ${PROVER_ENDPOINTS}
  --tierFee.optimistic ${BLOCK_PROPOSAL_FEE}
  --tierFee.sgx ${BLOCK_PROPOSAL_FEE}"

if [ -n "$TXPOOL_LOCALS" ]; then
  ARGS="${ARGS} --txPool.localsOnly"
  ARGS="${ARGS} --txPool.locals ${TXPOOL_LOCALS}"
fi

if [ -n "$MAX_TIER_FEE_BUMPS" ]; then
  ARGS="${ARGS} --tierFee.maxPriceBumps ${MAX_TIER_FEE_BUMPS}"
fi

if [ -n "$BLOCK_BUILDER_TIP" ]; then
  ARGS="${ARGS} --l1.blockBuilderTip ${BLOCK_BUILDER_TIP}"
fi

if [ "$BLOB_ALLOWED" == "true" ]; then
  ARGS="${ARGS} --l1.blobAllowed"
fi

# TXMGR Settings
if [ -n "$TX_FEE_LIMIT_MULTIPLIER" ]; then
  ARGS="${ARGS} --tx.feeLimitMultiplier ${TX_FEE_LIMIT_MULTIPLIER}"
fi

if [ -n "$TX_FEE_LIMIT_THRESHOLD" ]; then
  ARGS="${ARGS} --tx.feeLimitThreshold ${TX_FEE_LIMIT_THRESHOLD}"
fi

if [ -n "$TX_GAS_LIMIT" ]; then
  ARGS="${ARGS} --tx.gasLimit ${TX_GAS_LIMIT}"
fi

if [ -n "$TX_MIN_BASEFEE" ]; then
  ARGS="${ARGS} --tx.minBaseFee ${TX_MIN_BASEFEE}"
fi

if [ -n "$TX_MIN_TIP_CAP" ]; then
  ARGS="${ARGS} --tx.minTipCap ${TX_MIN_TIP_CAP}"
fi

if [ -n "$TX_NOT_IN_MEMPOOL" ]; then
  ARGS="${ARGS} --tx.notInMempoolTimeout ${TX_NOT_IN_MEMPOOL}"
fi

if [ -n "$TX_NUM_CONFIRMATIONS" ]; then
  ARGS="${ARGS} --tx.numConfirmations ${TX_NUM_CONFIRMATIONS}"
fi

if [ -n "$TX_RECEIPT_QUERY" ]; then
  ARGS="${ARGS} --tx.receiptQueryInterval ${TX_RECEIPT_QUERY}"
fi

if [ -n "$TX_RESUBMISSION" ]; then
  ARGS="${ARGS} --tx.resubmissionTimeout ${TX_RESUBMISSION}"
fi

if [ -n "$TX_SAFE_ABORT_NONCE_TOO_LOW" ]; then
  ARGS="${ARGS} --tx.safeAbortNonceTooLowCount ${TX_SAFE_ABORT_NONCE_TOO_LOW}"
fi

if [ -n "$TX_SEND_TIMEOUT" ]; then
  ARGS="${ARGS} --tx.sendTimeout ${TX_SEND_TIMEOUT}"
fi

# Run proposer
echo "Starting proposer..."
exec taiko-client proposer "${ARGS}"
