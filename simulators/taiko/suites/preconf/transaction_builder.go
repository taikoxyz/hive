package preconf

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/utils"
	"taiko/params"
)

// CalldataTransactionBuilder is responsible for building a TaikoL1.proposeBlock transaction with txList
// bytes saved in calldata.
type CalldataTransactionBuilder struct {
	rpc                     *rpc.Client
	proposerPrivateKey      *ecdsa.PrivateKey
	l2SuggestedFeeRecipient common.Address
	taikoL1Address          common.Address
	proverSetAddress        common.Address
	gasLimit                uint64
	chainConfig             *config.ChainConfig
	revertProtectionEnabled bool
}

// NewCalldataTransactionBuilder creates a new CalldataTransactionBuilder instance based on giving configurations.
func NewCalldataTransactionBuilder(
	rpc *rpc.Client,
	gasLimit uint64,
	chainConfig *config.ChainConfig,
	revertProtectionEnabled bool,
) *CalldataTransactionBuilder {
	return &CalldataTransactionBuilder{
		rpc,
		params.ParamToPriv("L1_PROPOSER_PRIV_KEY"),
		params.ParamToAddress("L2_SUGGESTED_FEE_RECIPIENT"),
		params.ParamToAddress("TAIKO_INBOX"),
		params.ParamToAddress("PROVER_SET"),
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
	// ABI encode the TaikoInbox.proposeBatch / ProverSet.proposeBatch parameters.
	var (
		to            = &b.taikoL1Address
		data          []byte
		encodedParams []byte
		blockParams   []pacayaBindings.ITaikoInboxBlockParams
		allTxs        types.Transactions
	)

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

	if encodedParams, err = encoding.EncodeBatchParams(&encoding.BatchParams{
		AnchorBlockId:            anchorHeader.Number.Uint64(),
		LastBlockTimestamp:       anchorHeader.Time,
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
		to = &b.proverSetAddress

		if data, err = encoding.ProverSetPacayaABI.Pack("proposeBatch", encodedParams, txListsBytes); err != nil {
			return nil, err
		}
	} else {
		if data, err = encoding.TaikoInboxABI.Pack("proposeBatch", encodedParams, txListsBytes); err != nil {
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
