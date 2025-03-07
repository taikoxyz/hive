package preconf

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/utils"
	"math/big"
)

// CalldataTransactionBuilder is responsible for building a TaikoL1.proposeBlock transaction with txList
// bytes saved in calldata.
type CalldataTransactionBuilder struct {
	rpc                     *rpc.Client
	proposerPrivateKey      *ecdsa.PrivateKey
	l2SuggestedFeeRecipient common.Address
	taikoL1Address          common.Address
	taikoWrapperAddress     common.Address
	proverSetAddress        common.Address
	gasLimit                uint64
	chainConfig             *config.ChainConfig
	revertProtectionEnabled bool
}

// NewCalldataTransactionBuilder creates a new CalldataTransactionBuilder instance based on giving configurations.
func NewCalldataTransactionBuilder(
	envs hivesim.Params,
	rpc *rpc.Client,
	gasLimit uint64,
	chainConfig *config.ChainConfig,
	revertProtectionEnabled bool,
) *CalldataTransactionBuilder {
	proposerPrivateKey, err := crypto.ToECDSA(common.FromHex(envs["L1_PROPOSER_PRIV_KEY"]))
	if err != nil {
		return nil
	}
	return &CalldataTransactionBuilder{
		rpc,
		proposerPrivateKey,
		common.HexToAddress(envs["L2_SUGGESTED_FEE_RECIPIENT"]),
		common.HexToAddress(envs["TAIKO_INBOX"]),
		common.HexToAddress(envs["TAIKO_WRAPPER"]),
		common.HexToAddress(envs["PROVER_SET"]),
		gasLimit,
		chainConfig,
		revertProtectionEnabled,
	}
}

// BuildPacaya implements the ProposeBlocksTransactionBuilder interface.
func (b *CalldataTransactionBuilder) BuildPacaya(
	ctx context.Context,
	txBatch []types.Transactions,
	anchorHeader *types.Header,
) (*txmgr.TxCandidate, error) {
	// ABI encode the TaikoWrapper.proposeBatch / ProverSet.proposeBatch parameters.
	var (
		to                    = &b.taikoWrapperAddress
		proposer              = crypto.PubkeyToAddress(b.proposerPrivateKey.PublicKey)
		data                  []byte
		encodedParams         []byte
		blockParams           []pacayaBindings.ITaikoInboxBlockParams
		forcedInclusionParams *encoding.BatchParams
		allTxs                types.Transactions
	)

	forcedInclusion, minTxsPerForcedInclusion, err := b.rpc.GetForcedInclusionPacaya(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forced inclusion: %w", err)
	}

	if b.proverSetAddress != rpc.ZeroAddress {
		to = &b.proverSetAddress
		proposer = b.proverSetAddress
	}

	if forcedInclusion != nil {
		blobParams, blockParams := buildParamsForForcedInclusion(forcedInclusion, minTxsPerForcedInclusion)
		forcedInclusionParams = &encoding.BatchParams{
			Proposer:                 proposer,
			Coinbase:                 b.l2SuggestedFeeRecipient,
			RevertIfNotFirstProposal: b.revertProtectionEnabled,
			BlobParams:               *blobParams,
			Blocks:                   blockParams,
		}
	}

	for _, txs := range txBatch {
		allTxs = append(allTxs, txs...)
		blockParams = append(blockParams, pacayaBindings.ITaikoInboxBlockParams{
			NumTransactions: uint16(len(txs)),
			TimeShift:       0,
			SignalSlots:     make([][32]byte, 0),
		})
	}

	rlpEncoded, err := rlp.EncodeToBytes(allTxs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode transactions: %w", err)
	}
	txListsBytes, err := utils.Compress(rlpEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to compress transactions: %w", err)
	}

	if encodedParams, err = encoding.EncodeBatchParamsWithForcedInclusion(
		forcedInclusionParams,
		&encoding.BatchParams{
			AnchorBlockId:            anchorHeader.Number.Uint64(),
			LastBlockTimestamp:       anchorHeader.Time,
			Proposer:                 proposer,
			Coinbase:                 b.l2SuggestedFeeRecipient,
			RevertIfNotFirstProposal: b.revertProtectionEnabled,
			BlobParams: encoding.BlobParams{
				ByteOffset: 0,
				ByteSize:   uint32(len(txListsBytes)),
			},
			Blocks: blockParams,
		}); err != nil {
		return nil, err
	}

	if b.proverSetAddress != rpc.ZeroAddress {
		if data, err = encoding.ProverSetPacayaABI.Pack("proposeBatch", encodedParams, txListsBytes); err != nil {
			return nil, err
		}
	} else {
		if data, err = encoding.TaikoWrapperABI.Pack("proposeBatch", encodedParams, txListsBytes); err != nil {
			return nil, err
		}
	}

	return &txmgr.TxCandidate{
		TxData:   data,
		Blobs:    nil,
		To:       to,
		GasLimit: b.gasLimit,
	}, nil
}

// buildParamsForForcedInclusion builds the blob params and the block params
// for the given forced inclusion.
func buildParamsForForcedInclusion(
	forcedInclusion *pacayaBindings.IForcedInclusionStoreForcedInclusion,
	minTxsPerForcedInclusion *big.Int,
) (*encoding.BlobParams, []pacayaBindings.ITaikoInboxBlockParams) {
	if forcedInclusion == nil {
		return nil, nil
	}
	return &encoding.BlobParams{
			BlobHashes: [][32]byte{forcedInclusion.BlobHash},
			NumBlobs:   0,
			ByteOffset: forcedInclusion.BlobByteOffset,
			ByteSize:   forcedInclusion.BlobByteSize,
		}, []pacayaBindings.ITaikoInboxBlockParams{
			{
				NumTransactions: uint16(minTxsPerForcedInclusion.Uint64()),
				TimeShift:       0,
				SignalSlots:     make([][32]byte, 0),
			},
		}
}
