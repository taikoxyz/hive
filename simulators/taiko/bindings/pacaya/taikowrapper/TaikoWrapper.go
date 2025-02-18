// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package taikowrapper

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IForcedInclusionStoreForcedInclusion is an auto generated low-level Go binding around an user-defined struct.
type IForcedInclusionStoreForcedInclusion struct {
	BlobHash         [32]byte
	FeeInGwei        uint64
	CreatedAtBatchId uint64
	BlobByteOffset   uint32
	BlobByteSize     uint32
}

// ITaikoInboxBatchInfo is an auto generated low-level Go binding around an user-defined struct.
type ITaikoInboxBatchInfo struct {
	TxsHash            [32]byte
	Blocks             []ITaikoInboxBlockParams
	BlobHashes         [][32]byte
	ExtraData          [32]byte
	Coinbase           common.Address
	ProposedIn         uint64
	BlobByteOffset     uint32
	BlobByteSize       uint32
	GasLimit           uint32
	LastBlockId        uint64
	LastBlockTimestamp uint64
	AnchorBlockId      uint64
	AnchorBlockHash    [32]byte
	BaseFeeConfig      LibSharedDataBaseFeeConfig
}

// ITaikoInboxBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type ITaikoInboxBatchMetadata struct {
	InfoHash   [32]byte
	Proposer   common.Address
	BatchId    uint64
	ProposedAt uint64
}

// ITaikoInboxBlockParams is an auto generated low-level Go binding around an user-defined struct.
type ITaikoInboxBlockParams struct {
	NumTransactions uint16
	TimeShift       uint8
	SignalSlots     [][32]byte
}

// LibSharedDataBaseFeeConfig is an auto generated low-level Go binding around an user-defined struct.
type LibSharedDataBaseFeeConfig struct {
	AdjustmentQuotient     uint8
	SharingPctg            uint8
	GasIssuancePerSecond   uint32
	MinGasExcess           uint64
	MaxGasIssuancePerBlock uint32
}

// TaikoWrapperMetaData contains all meta data concerning the TaikoWrapper contract.
var TaikoWrapperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_inbox\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_forcedInclusionStore\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_preconfRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MIN_TXS_PER_FORCED_INCLUSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forcedInclusionStore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIForcedInclusionStore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProposeBatch\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"preconfRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeBatch\",\"inputs\":[{\"name\":\"_params\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_txList\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITaikoInbox.BatchInfo\",\"components\":[{\"name\":\"txsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blocks\",\"type\":\"tuple[]\",\"internalType\":\"structITaikoInbox.BlockParams[]\",\"components\":[{\"name\":\"numTransactions\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"timeShift\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]},{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"extraData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"coinbase\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposedIn\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobByteOffset\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobByteSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"lastBlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastBlockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchorBlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchorBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"baseFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structLibSharedData.BaseFeeConfig\",\"components\":[{\"name\":\"adjustmentQuotient\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasIssuancePerSecond\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"minGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxGasIssuancePerBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITaikoInbox.BatchMetadata\",\"components\":[{\"name\":\"infoHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"batchId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForcedInclusionProcessed\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIForcedInclusionStore.ForcedInclusion\",\"components\":[{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"feeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"createdAtBatchId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobByteOffset\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobByteSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlobByteOffset\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlobByteSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlobHash\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlobHashesSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlockTxs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBlocks\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OldestForcedInclusionDue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RESOLVER_NOT_FOUND\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// TaikoWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use TaikoWrapperMetaData.ABI instead.
var TaikoWrapperABI = TaikoWrapperMetaData.ABI

// TaikoWrapper is an auto generated Go binding around an Ethereum contract.
type TaikoWrapper struct {
	TaikoWrapperCaller     // Read-only binding to the contract
	TaikoWrapperTransactor // Write-only binding to the contract
	TaikoWrapperFilterer   // Log filterer for contract events
}

// TaikoWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type TaikoWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TaikoWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TaikoWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TaikoWrapperSession struct {
	Contract     *TaikoWrapper     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaikoWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TaikoWrapperCallerSession struct {
	Contract *TaikoWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TaikoWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TaikoWrapperTransactorSession struct {
	Contract     *TaikoWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TaikoWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type TaikoWrapperRaw struct {
	Contract *TaikoWrapper // Generic contract binding to access the raw methods on
}

// TaikoWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TaikoWrapperCallerRaw struct {
	Contract *TaikoWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// TaikoWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TaikoWrapperTransactorRaw struct {
	Contract *TaikoWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTaikoWrapper creates a new instance of TaikoWrapper, bound to a specific deployed contract.
func NewTaikoWrapper(address common.Address, backend bind.ContractBackend) (*TaikoWrapper, error) {
	contract, err := bindTaikoWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapper{TaikoWrapperCaller: TaikoWrapperCaller{contract: contract}, TaikoWrapperTransactor: TaikoWrapperTransactor{contract: contract}, TaikoWrapperFilterer: TaikoWrapperFilterer{contract: contract}}, nil
}

// NewTaikoWrapperCaller creates a new read-only instance of TaikoWrapper, bound to a specific deployed contract.
func NewTaikoWrapperCaller(address common.Address, caller bind.ContractCaller) (*TaikoWrapperCaller, error) {
	contract, err := bindTaikoWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperCaller{contract: contract}, nil
}

// NewTaikoWrapperTransactor creates a new write-only instance of TaikoWrapper, bound to a specific deployed contract.
func NewTaikoWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*TaikoWrapperTransactor, error) {
	contract, err := bindTaikoWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperTransactor{contract: contract}, nil
}

// NewTaikoWrapperFilterer creates a new log filterer instance of TaikoWrapper, bound to a specific deployed contract.
func NewTaikoWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*TaikoWrapperFilterer, error) {
	contract, err := bindTaikoWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperFilterer{contract: contract}, nil
}

// bindTaikoWrapper binds a generic wrapper to an already deployed contract.
func bindTaikoWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TaikoWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoWrapper *TaikoWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoWrapper.Contract.TaikoWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoWrapper *TaikoWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.TaikoWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoWrapper *TaikoWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.TaikoWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoWrapper *TaikoWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoWrapper *TaikoWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoWrapper *TaikoWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.contract.Transact(opts, method, params...)
}

// MINTXSPERFORCEDINCLUSION is a free data retrieval call binding the contract method 0x936ee868.
//
// Solidity: function MIN_TXS_PER_FORCED_INCLUSION() view returns(uint16)
func (_TaikoWrapper *TaikoWrapperCaller) MINTXSPERFORCEDINCLUSION(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "MIN_TXS_PER_FORCED_INCLUSION")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MINTXSPERFORCEDINCLUSION is a free data retrieval call binding the contract method 0x936ee868.
//
// Solidity: function MIN_TXS_PER_FORCED_INCLUSION() view returns(uint16)
func (_TaikoWrapper *TaikoWrapperSession) MINTXSPERFORCEDINCLUSION() (uint16, error) {
	return _TaikoWrapper.Contract.MINTXSPERFORCEDINCLUSION(&_TaikoWrapper.CallOpts)
}

// MINTXSPERFORCEDINCLUSION is a free data retrieval call binding the contract method 0x936ee868.
//
// Solidity: function MIN_TXS_PER_FORCED_INCLUSION() view returns(uint16)
func (_TaikoWrapper *TaikoWrapperCallerSession) MINTXSPERFORCEDINCLUSION() (uint16, error) {
	return _TaikoWrapper.Contract.MINTXSPERFORCEDINCLUSION(&_TaikoWrapper.CallOpts)
}

// ForcedInclusionStore is a free data retrieval call binding the contract method 0xf749e9d4.
//
// Solidity: function forcedInclusionStore() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) ForcedInclusionStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "forcedInclusionStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ForcedInclusionStore is a free data retrieval call binding the contract method 0xf749e9d4.
//
// Solidity: function forcedInclusionStore() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) ForcedInclusionStore() (common.Address, error) {
	return _TaikoWrapper.Contract.ForcedInclusionStore(&_TaikoWrapper.CallOpts)
}

// ForcedInclusionStore is a free data retrieval call binding the contract method 0xf749e9d4.
//
// Solidity: function forcedInclusionStore() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) ForcedInclusionStore() (common.Address, error) {
	return _TaikoWrapper.Contract.ForcedInclusionStore(&_TaikoWrapper.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) Impl() (common.Address, error) {
	return _TaikoWrapper.Contract.Impl(&_TaikoWrapper.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) Impl() (common.Address, error) {
	return _TaikoWrapper.Contract.Impl(&_TaikoWrapper.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoWrapper *TaikoWrapperCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoWrapper *TaikoWrapperSession) InNonReentrant() (bool, error) {
	return _TaikoWrapper.Contract.InNonReentrant(&_TaikoWrapper.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoWrapper *TaikoWrapperCallerSession) InNonReentrant() (bool, error) {
	return _TaikoWrapper.Contract.InNonReentrant(&_TaikoWrapper.CallOpts)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) Inbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "inbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) Inbox() (common.Address, error) {
	return _TaikoWrapper.Contract.Inbox(&_TaikoWrapper.CallOpts)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) Inbox() (common.Address, error) {
	return _TaikoWrapper.Contract.Inbox(&_TaikoWrapper.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) Owner() (common.Address, error) {
	return _TaikoWrapper.Contract.Owner(&_TaikoWrapper.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) Owner() (common.Address, error) {
	return _TaikoWrapper.Contract.Owner(&_TaikoWrapper.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoWrapper *TaikoWrapperCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoWrapper *TaikoWrapperSession) Paused() (bool, error) {
	return _TaikoWrapper.Contract.Paused(&_TaikoWrapper.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoWrapper *TaikoWrapperCallerSession) Paused() (bool, error) {
	return _TaikoWrapper.Contract.Paused(&_TaikoWrapper.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) PendingOwner() (common.Address, error) {
	return _TaikoWrapper.Contract.PendingOwner(&_TaikoWrapper.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) PendingOwner() (common.Address, error) {
	return _TaikoWrapper.Contract.PendingOwner(&_TaikoWrapper.CallOpts)
}

// PreconfRouter is a free data retrieval call binding the contract method 0xf4bb9077.
//
// Solidity: function preconfRouter() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) PreconfRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "preconfRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreconfRouter is a free data retrieval call binding the contract method 0xf4bb9077.
//
// Solidity: function preconfRouter() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) PreconfRouter() (common.Address, error) {
	return _TaikoWrapper.Contract.PreconfRouter(&_TaikoWrapper.CallOpts)
}

// PreconfRouter is a free data retrieval call binding the contract method 0xf4bb9077.
//
// Solidity: function preconfRouter() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) PreconfRouter() (common.Address, error) {
	return _TaikoWrapper.Contract.PreconfRouter(&_TaikoWrapper.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoWrapper *TaikoWrapperCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoWrapper *TaikoWrapperSession) ProxiableUUID() ([32]byte, error) {
	return _TaikoWrapper.Contract.ProxiableUUID(&_TaikoWrapper.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoWrapper *TaikoWrapperCallerSession) ProxiableUUID() ([32]byte, error) {
	return _TaikoWrapper.Contract.ProxiableUUID(&_TaikoWrapper.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoWrapper *TaikoWrapperCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoWrapper.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoWrapper *TaikoWrapperSession) Resolver() (common.Address, error) {
	return _TaikoWrapper.Contract.Resolver(&_TaikoWrapper.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoWrapper *TaikoWrapperCallerSession) Resolver() (common.Address, error) {
	return _TaikoWrapper.Contract.Resolver(&_TaikoWrapper.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoWrapper *TaikoWrapperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoWrapper *TaikoWrapperSession) AcceptOwnership() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.AcceptOwnership(&_TaikoWrapper.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.AcceptOwnership(&_TaikoWrapper.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_TaikoWrapper *TaikoWrapperTransactor) Init(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "init", _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_TaikoWrapper *TaikoWrapperSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Init(&_TaikoWrapper.TransactOpts, _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Init(&_TaikoWrapper.TransactOpts, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoWrapper *TaikoWrapperTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoWrapper *TaikoWrapperSession) Pause() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Pause(&_TaikoWrapper.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) Pause() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Pause(&_TaikoWrapper.TransactOpts)
}

// ProposeBatch is a paid mutator transaction binding the contract method 0x47faad14.
//
// Solidity: function proposeBatch(bytes _params, bytes _txList) returns((bytes32,(uint16,uint8,bytes32[])[],bytes32[],bytes32,address,uint64,uint32,uint32,uint32,uint64,uint64,uint64,bytes32,(uint8,uint8,uint32,uint64,uint32)), (bytes32,address,uint64,uint64))
func (_TaikoWrapper *TaikoWrapperTransactor) ProposeBatch(opts *bind.TransactOpts, _params []byte, _txList []byte) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "proposeBatch", _params, _txList)
}

// ProposeBatch is a paid mutator transaction binding the contract method 0x47faad14.
//
// Solidity: function proposeBatch(bytes _params, bytes _txList) returns((bytes32,(uint16,uint8,bytes32[])[],bytes32[],bytes32,address,uint64,uint32,uint32,uint32,uint64,uint64,uint64,bytes32,(uint8,uint8,uint32,uint64,uint32)), (bytes32,address,uint64,uint64))
func (_TaikoWrapper *TaikoWrapperSession) ProposeBatch(_params []byte, _txList []byte) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.ProposeBatch(&_TaikoWrapper.TransactOpts, _params, _txList)
}

// ProposeBatch is a paid mutator transaction binding the contract method 0x47faad14.
//
// Solidity: function proposeBatch(bytes _params, bytes _txList) returns((bytes32,(uint16,uint8,bytes32[])[],bytes32[],bytes32,address,uint64,uint32,uint32,uint32,uint64,uint64,uint64,bytes32,(uint8,uint8,uint32,uint64,uint32)), (bytes32,address,uint64,uint64))
func (_TaikoWrapper *TaikoWrapperTransactorSession) ProposeBatch(_params []byte, _txList []byte) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.ProposeBatch(&_TaikoWrapper.TransactOpts, _params, _txList)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoWrapper *TaikoWrapperTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoWrapper *TaikoWrapperSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.RenounceOwnership(&_TaikoWrapper.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.RenounceOwnership(&_TaikoWrapper.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoWrapper *TaikoWrapperTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoWrapper *TaikoWrapperSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.TransferOwnership(&_TaikoWrapper.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.TransferOwnership(&_TaikoWrapper.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoWrapper *TaikoWrapperTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoWrapper *TaikoWrapperSession) Unpause() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Unpause(&_TaikoWrapper.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) Unpause() (*types.Transaction, error) {
	return _TaikoWrapper.Contract.Unpause(&_TaikoWrapper.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoWrapper *TaikoWrapperTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoWrapper *TaikoWrapperSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.UpgradeTo(&_TaikoWrapper.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.UpgradeTo(&_TaikoWrapper.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoWrapper *TaikoWrapperTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoWrapper.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoWrapper *TaikoWrapperSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.UpgradeToAndCall(&_TaikoWrapper.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoWrapper *TaikoWrapperTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoWrapper.Contract.UpgradeToAndCall(&_TaikoWrapper.TransactOpts, newImplementation, data)
}

// TaikoWrapperAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TaikoWrapper contract.
type TaikoWrapperAdminChangedIterator struct {
	Event *TaikoWrapperAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperAdminChanged represents a AdminChanged event raised by the TaikoWrapper contract.
type TaikoWrapperAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TaikoWrapperAdminChangedIterator, error) {

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperAdminChangedIterator{contract: _TaikoWrapper.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TaikoWrapperAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperAdminChanged)
				if err := _TaikoWrapper.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseAdminChanged(log types.Log) (*TaikoWrapperAdminChanged, error) {
	event := new(TaikoWrapperAdminChanged)
	if err := _TaikoWrapper.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the TaikoWrapper contract.
type TaikoWrapperBeaconUpgradedIterator struct {
	Event *TaikoWrapperBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperBeaconUpgraded represents a BeaconUpgraded event raised by the TaikoWrapper contract.
type TaikoWrapperBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*TaikoWrapperBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperBeaconUpgradedIterator{contract: _TaikoWrapper.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *TaikoWrapperBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperBeaconUpgraded)
				if err := _TaikoWrapper.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseBeaconUpgraded(log types.Log) (*TaikoWrapperBeaconUpgraded, error) {
	event := new(TaikoWrapperBeaconUpgraded)
	if err := _TaikoWrapper.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperForcedInclusionProcessedIterator is returned from FilterForcedInclusionProcessed and is used to iterate over the raw logs and unpacked data for ForcedInclusionProcessed events raised by the TaikoWrapper contract.
type TaikoWrapperForcedInclusionProcessedIterator struct {
	Event *TaikoWrapperForcedInclusionProcessed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperForcedInclusionProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperForcedInclusionProcessed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperForcedInclusionProcessed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperForcedInclusionProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperForcedInclusionProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperForcedInclusionProcessed represents a ForcedInclusionProcessed event raised by the TaikoWrapper contract.
type TaikoWrapperForcedInclusionProcessed struct {
	Arg0 IForcedInclusionStoreForcedInclusion
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterForcedInclusionProcessed is a free log retrieval operation binding the contract event 0x16736f36093e9ebaadd610bd1fc9e0653ff56964ebeeccb0b17a0b51a0b4c6fc.
//
// Solidity: event ForcedInclusionProcessed((bytes32,uint64,uint64,uint32,uint32) arg0)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterForcedInclusionProcessed(opts *bind.FilterOpts) (*TaikoWrapperForcedInclusionProcessedIterator, error) {

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "ForcedInclusionProcessed")
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperForcedInclusionProcessedIterator{contract: _TaikoWrapper.contract, event: "ForcedInclusionProcessed", logs: logs, sub: sub}, nil
}

// WatchForcedInclusionProcessed is a free log subscription operation binding the contract event 0x16736f36093e9ebaadd610bd1fc9e0653ff56964ebeeccb0b17a0b51a0b4c6fc.
//
// Solidity: event ForcedInclusionProcessed((bytes32,uint64,uint64,uint32,uint32) arg0)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchForcedInclusionProcessed(opts *bind.WatchOpts, sink chan<- *TaikoWrapperForcedInclusionProcessed) (event.Subscription, error) {

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "ForcedInclusionProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperForcedInclusionProcessed)
				if err := _TaikoWrapper.contract.UnpackLog(event, "ForcedInclusionProcessed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseForcedInclusionProcessed is a log parse operation binding the contract event 0x16736f36093e9ebaadd610bd1fc9e0653ff56964ebeeccb0b17a0b51a0b4c6fc.
//
// Solidity: event ForcedInclusionProcessed((bytes32,uint64,uint64,uint32,uint32) arg0)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseForcedInclusionProcessed(log types.Log) (*TaikoWrapperForcedInclusionProcessed, error) {
	event := new(TaikoWrapperForcedInclusionProcessed)
	if err := _TaikoWrapper.contract.UnpackLog(event, "ForcedInclusionProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TaikoWrapper contract.
type TaikoWrapperInitializedIterator struct {
	Event *TaikoWrapperInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperInitialized represents a Initialized event raised by the TaikoWrapper contract.
type TaikoWrapperInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterInitialized(opts *bind.FilterOpts) (*TaikoWrapperInitializedIterator, error) {

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperInitializedIterator{contract: _TaikoWrapper.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TaikoWrapperInitialized) (event.Subscription, error) {

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperInitialized)
				if err := _TaikoWrapper.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseInitialized(log types.Log) (*TaikoWrapperInitialized, error) {
	event := new(TaikoWrapperInitialized)
	if err := _TaikoWrapper.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the TaikoWrapper contract.
type TaikoWrapperOwnershipTransferStartedIterator struct {
	Event *TaikoWrapperOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the TaikoWrapper contract.
type TaikoWrapperOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaikoWrapperOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperOwnershipTransferStartedIterator{contract: _TaikoWrapper.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *TaikoWrapperOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperOwnershipTransferStarted)
				if err := _TaikoWrapper.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseOwnershipTransferStarted(log types.Log) (*TaikoWrapperOwnershipTransferStarted, error) {
	event := new(TaikoWrapperOwnershipTransferStarted)
	if err := _TaikoWrapper.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TaikoWrapper contract.
type TaikoWrapperOwnershipTransferredIterator struct {
	Event *TaikoWrapperOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperOwnershipTransferred represents a OwnershipTransferred event raised by the TaikoWrapper contract.
type TaikoWrapperOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaikoWrapperOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperOwnershipTransferredIterator{contract: _TaikoWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TaikoWrapperOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperOwnershipTransferred)
				if err := _TaikoWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*TaikoWrapperOwnershipTransferred, error) {
	event := new(TaikoWrapperOwnershipTransferred)
	if err := _TaikoWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the TaikoWrapper contract.
type TaikoWrapperPausedIterator struct {
	Event *TaikoWrapperPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperPaused represents a Paused event raised by the TaikoWrapper contract.
type TaikoWrapperPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterPaused(opts *bind.FilterOpts) (*TaikoWrapperPausedIterator, error) {

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperPausedIterator{contract: _TaikoWrapper.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *TaikoWrapperPaused) (event.Subscription, error) {

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperPaused)
				if err := _TaikoWrapper.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) ParsePaused(log types.Log) (*TaikoWrapperPaused, error) {
	event := new(TaikoWrapperPaused)
	if err := _TaikoWrapper.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the TaikoWrapper contract.
type TaikoWrapperUnpausedIterator struct {
	Event *TaikoWrapperUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperUnpaused represents a Unpaused event raised by the TaikoWrapper contract.
type TaikoWrapperUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterUnpaused(opts *bind.FilterOpts) (*TaikoWrapperUnpausedIterator, error) {

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperUnpausedIterator{contract: _TaikoWrapper.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *TaikoWrapperUnpaused) (event.Subscription, error) {

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperUnpaused)
				if err := _TaikoWrapper.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseUnpaused(log types.Log) (*TaikoWrapperUnpaused, error) {
	event := new(TaikoWrapperUnpaused)
	if err := _TaikoWrapper.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoWrapperUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TaikoWrapper contract.
type TaikoWrapperUpgradedIterator struct {
	Event *TaikoWrapperUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaikoWrapperUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoWrapperUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaikoWrapperUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaikoWrapperUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoWrapperUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoWrapperUpgraded represents a Upgraded event raised by the TaikoWrapper contract.
type TaikoWrapperUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TaikoWrapper *TaikoWrapperFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TaikoWrapperUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TaikoWrapper.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TaikoWrapperUpgradedIterator{contract: _TaikoWrapper.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TaikoWrapper *TaikoWrapperFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TaikoWrapperUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TaikoWrapper.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoWrapperUpgraded)
				if err := _TaikoWrapper.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TaikoWrapper *TaikoWrapperFilterer) ParseUpgraded(log types.Log) (*TaikoWrapperUpgraded, error) {
	event := new(TaikoWrapperUpgraded)
	if err := _TaikoWrapper.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
