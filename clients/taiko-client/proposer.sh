#!/bin/bash

source /common.sh

# Set default values
TAIKO_L1_ADDRESS=${TAIKO_L1_ADDRESS:-0x06a9Ab27c7e2255df1815E6CC0168d7755Feb19a}
TAIKO_L2_ADDRESS=${TAIKO_L2_ADDRESS:-0x1670000000000000000000000000000000010001}
TAIKO_TOKEN_L1_ADDRESS=${TAIKO_TOKEN_L1_ADDRESS:-0x10dea67478c5F8C5E2D90e5E9B26dBe60c54d800}
ASSIGNMENT_HOOK_L1_ADDRESS=${ASSIGNMENT_HOOK_L1_ADDRESS:-0x537a2f0D3a5879b41BCb5A2afE2EA5c4961796F6}
L2_SUGGESTED_FEE_RECIPIENT=${L2_SUGGESTED_FEE_RECIPIENT:-}
BLOCK_PROPOSAL_FEE=${BLOCK_PROPOSAL_FEE:-1}

check_env "L1_ENDPOINT_WS"
check_env "L1_BEACON_HTTP"
check_env "L2_EXECUTION_ENGINE"
check_env "L1_PROPOSER_PRIVATE_KEY"
check_env "PROVER_ENDPOINTS"

ARGS="--l1.ws ${L1_ENDPOINT_WS}
  --l2.http http://${L2_EXECUTION_ENGINE}:8545
  --l2.auth http://${L2_EXECUTION_ENGINE}:8551
  --taikoL1 ${TAIKO_L1_ADDRESS}
  --taikoL2 ${TAIKO_L2_ADDRESS}
  --taikoToken ${TAIKO_TOKEN_L1_ADDRESS}
  --assignmentHookAddress ${ASSIGNMENT_HOOK_L1_ADDRESS}
  --jwtSecret /jwt.hex
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

