// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package taikoanchor

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

// LibSharedDataBaseFeeConfig is an auto generated low-level Go binding around an user-defined struct.
type LibSharedDataBaseFeeConfig struct {
	AdjustmentQuotient     uint8
	SharingPctg            uint8
	GasIssuancePerSecond   uint32
	MinGasExcess           uint64
	MaxGasIssuancePerBlock uint32
}

// TaikoAnchorMetaData contains all meta data concerning the TaikoAnchor contract.
var TaikoAnchorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_resolver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_signalService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pacayaForkHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"GOLDEN_TOUCH_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"adjustExcess\",\"inputs\":[{\"name\":\"_currGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_currGasTarget\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_newGasTarget\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"newGasExcess_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"_l1BlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_l1StateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_l1BlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorV2\",\"inputs\":[{\"name\":\"_anchorBlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_anchorStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_baseFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structLibSharedData.BaseFeeConfig\",\"components\":[{\"name\":\"adjustmentQuotient\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasIssuancePerSecond\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"minGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxGasIssuancePerBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorV3\",\"inputs\":[{\"name\":\"_anchorBlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_anchorStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_baseFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structLibSharedData.BaseFeeConfig\",\"components\":[{\"name\":\"adjustmentQuotient\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasIssuancePerSecond\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"minGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxGasIssuancePerBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"_signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"calculateBaseFee\",\"inputs\":[{\"name\":\"_baseFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structLibSharedData.BaseFeeConfig\",\"components\":[{\"name\":\"adjustmentQuotient\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasIssuancePerSecond\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"minGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxGasIssuancePerBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"_blocktime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_parentGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"basefee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"parentGasExcess_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getBasefee\",\"inputs\":[{\"name\":\"_anchorBlockId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"basefee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"parentGasExcess_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getBasefeeV2\",\"inputs\":[{\"name\":\"_parentGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_blockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_baseFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structLibSharedData.BaseFeeConfig\",\"components\":[{\"name\":\"adjustmentQuotient\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasIssuancePerSecond\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"minGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxGasIssuancePerBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[{\"name\":\"basefee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newGasTarget_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newGasExcess_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockHash\",\"inputs\":[{\"name\":\"_blockId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l1ChainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_initialGasExcess\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isOnL1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"l1ChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastSyncedBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pacayaForkHeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parentGasExcess\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parentGasTarget\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parentTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"publicInputHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"signalService\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISignalService\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"skipFeeCheck\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Anchored\",\"inputs\":[{\"name\":\"parentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"parentGasExcess\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP1559Update\",\"inputs\":[{\"name\":\"oldGasTarget\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newGasTarget\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"oldGasExcess\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newGasExcess\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"basefee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ETH_TRANSFER_FAILED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_BASEFEE_MISMATCH\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_DEPRECATED_METHOD\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_FORK_ERROR\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_INVALID_L1_CHAIN_ID\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_INVALID_L2_CHAIN_ID\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_INVALID_SENDER\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_PUBLIC_INPUT_HASH_MISMATCH\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L2_TOO_LATE\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RESOLVER_NOT_FOUND\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// TaikoAnchorABI is the input ABI used to generate the binding from.
// Deprecated: Use TaikoAnchorMetaData.ABI instead.
var TaikoAnchorABI = TaikoAnchorMetaData.ABI

// TaikoAnchor is an auto generated Go binding around an Ethereum contract.
type TaikoAnchor struct {
	TaikoAnchorCaller     // Read-only binding to the contract
	TaikoAnchorTransactor // Write-only binding to the contract
	TaikoAnchorFilterer   // Log filterer for contract events
}

// TaikoAnchorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TaikoAnchorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoAnchorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TaikoAnchorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoAnchorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TaikoAnchorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoAnchorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TaikoAnchorSession struct {
	Contract     *TaikoAnchor      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaikoAnchorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TaikoAnchorCallerSession struct {
	Contract *TaikoAnchorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TaikoAnchorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TaikoAnchorTransactorSession struct {
	Contract     *TaikoAnchorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TaikoAnchorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TaikoAnchorRaw struct {
	Contract *TaikoAnchor // Generic contract binding to access the raw methods on
}

// TaikoAnchorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TaikoAnchorCallerRaw struct {
	Contract *TaikoAnchorCaller // Generic read-only contract binding to access the raw methods on
}

// TaikoAnchorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TaikoAnchorTransactorRaw struct {
	Contract *TaikoAnchorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTaikoAnchor creates a new instance of TaikoAnchor, bound to a specific deployed contract.
func NewTaikoAnchor(address common.Address, backend bind.ContractBackend) (*TaikoAnchor, error) {
	contract, err := bindTaikoAnchor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchor{TaikoAnchorCaller: TaikoAnchorCaller{contract: contract}, TaikoAnchorTransactor: TaikoAnchorTransactor{contract: contract}, TaikoAnchorFilterer: TaikoAnchorFilterer{contract: contract}}, nil
}

// NewTaikoAnchorCaller creates a new read-only instance of TaikoAnchor, bound to a specific deployed contract.
func NewTaikoAnchorCaller(address common.Address, caller bind.ContractCaller) (*TaikoAnchorCaller, error) {
	contract, err := bindTaikoAnchor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorCaller{contract: contract}, nil
}

// NewTaikoAnchorTransactor creates a new write-only instance of TaikoAnchor, bound to a specific deployed contract.
func NewTaikoAnchorTransactor(address common.Address, transactor bind.ContractTransactor) (*TaikoAnchorTransactor, error) {
	contract, err := bindTaikoAnchor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorTransactor{contract: contract}, nil
}

// NewTaikoAnchorFilterer creates a new log filterer instance of TaikoAnchor, bound to a specific deployed contract.
func NewTaikoAnchorFilterer(address common.Address, filterer bind.ContractFilterer) (*TaikoAnchorFilterer, error) {
	contract, err := bindTaikoAnchor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorFilterer{contract: contract}, nil
}

// bindTaikoAnchor binds a generic wrapper to an already deployed contract.
func bindTaikoAnchor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TaikoAnchorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoAnchor *TaikoAnchorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoAnchor.Contract.TaikoAnchorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoAnchor *TaikoAnchorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.TaikoAnchorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoAnchor *TaikoAnchorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.TaikoAnchorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoAnchor *TaikoAnchorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoAnchor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoAnchor *TaikoAnchorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoAnchor *TaikoAnchorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.contract.Transact(opts, method, params...)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) GOLDENTOUCHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "GOLDEN_TOUCH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _TaikoAnchor.Contract.GOLDENTOUCHADDRESS(&_TaikoAnchor.CallOpts)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _TaikoAnchor.Contract.GOLDENTOUCHADDRESS(&_TaikoAnchor.CallOpts)
}

// AdjustExcess is a free data retrieval call binding the contract method 0x136dc4a8.
//
// Solidity: function adjustExcess(uint64 _currGasExcess, uint64 _currGasTarget, uint64 _newGasTarget) pure returns(uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorCaller) AdjustExcess(opts *bind.CallOpts, _currGasExcess uint64, _currGasTarget uint64, _newGasTarget uint64) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "adjustExcess", _currGasExcess, _currGasTarget, _newGasTarget)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// AdjustExcess is a free data retrieval call binding the contract method 0x136dc4a8.
//
// Solidity: function adjustExcess(uint64 _currGasExcess, uint64 _currGasTarget, uint64 _newGasTarget) pure returns(uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorSession) AdjustExcess(_currGasExcess uint64, _currGasTarget uint64, _newGasTarget uint64) (uint64, error) {
	return _TaikoAnchor.Contract.AdjustExcess(&_TaikoAnchor.CallOpts, _currGasExcess, _currGasTarget, _newGasTarget)
}

// AdjustExcess is a free data retrieval call binding the contract method 0x136dc4a8.
//
// Solidity: function adjustExcess(uint64 _currGasExcess, uint64 _currGasTarget, uint64 _newGasTarget) pure returns(uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorCallerSession) AdjustExcess(_currGasExcess uint64, _currGasTarget uint64, _newGasTarget uint64) (uint64, error) {
	return _TaikoAnchor.Contract.AdjustExcess(&_TaikoAnchor.CallOpts, _currGasExcess, _currGasTarget, _newGasTarget)
}

// CalculateBaseFee is a free data retrieval call binding the contract method 0xe902461a.
//
// Solidity: function calculateBaseFee((uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, uint64 _blocktime, uint64 _parentGasExcess, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorCaller) CalculateBaseFee(opts *bind.CallOpts, _baseFeeConfig LibSharedDataBaseFeeConfig, _blocktime uint64, _parentGasExcess uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "calculateBaseFee", _baseFeeConfig, _blocktime, _parentGasExcess, _parentGasUsed)

	outstruct := new(struct {
		Basefee         *big.Int
		ParentGasExcess uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Basefee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ParentGasExcess = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// CalculateBaseFee is a free data retrieval call binding the contract method 0xe902461a.
//
// Solidity: function calculateBaseFee((uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, uint64 _blocktime, uint64 _parentGasExcess, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorSession) CalculateBaseFee(_baseFeeConfig LibSharedDataBaseFeeConfig, _blocktime uint64, _parentGasExcess uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.CalculateBaseFee(&_TaikoAnchor.CallOpts, _baseFeeConfig, _blocktime, _parentGasExcess, _parentGasUsed)
}

// CalculateBaseFee is a free data retrieval call binding the contract method 0xe902461a.
//
// Solidity: function calculateBaseFee((uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, uint64 _blocktime, uint64 _parentGasExcess, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorCallerSession) CalculateBaseFee(_baseFeeConfig LibSharedDataBaseFeeConfig, _blocktime uint64, _parentGasExcess uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.CalculateBaseFee(&_TaikoAnchor.CallOpts, _baseFeeConfig, _blocktime, _parentGasExcess, _parentGasUsed)
}

// GetBasefee is a free data retrieval call binding the contract method 0xa7e022d1.
//
// Solidity: function getBasefee(uint64 _anchorBlockId, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorCaller) GetBasefee(opts *bind.CallOpts, _anchorBlockId uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "getBasefee", _anchorBlockId, _parentGasUsed)

	outstruct := new(struct {
		Basefee         *big.Int
		ParentGasExcess uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Basefee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ParentGasExcess = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetBasefee is a free data retrieval call binding the contract method 0xa7e022d1.
//
// Solidity: function getBasefee(uint64 _anchorBlockId, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorSession) GetBasefee(_anchorBlockId uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.GetBasefee(&_TaikoAnchor.CallOpts, _anchorBlockId, _parentGasUsed)
}

// GetBasefee is a free data retrieval call binding the contract method 0xa7e022d1.
//
// Solidity: function getBasefee(uint64 _anchorBlockId, uint32 _parentGasUsed) pure returns(uint256 basefee_, uint64 parentGasExcess_)
func (_TaikoAnchor *TaikoAnchorCallerSession) GetBasefee(_anchorBlockId uint64, _parentGasUsed uint32) (struct {
	Basefee         *big.Int
	ParentGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.GetBasefee(&_TaikoAnchor.CallOpts, _anchorBlockId, _parentGasUsed)
}

// GetBasefeeV2 is a free data retrieval call binding the contract method 0x893f5460.
//
// Solidity: function getBasefeeV2(uint32 _parentGasUsed, uint64 _blockTimestamp, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) view returns(uint256 basefee_, uint64 newGasTarget_, uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorCaller) GetBasefeeV2(opts *bind.CallOpts, _parentGasUsed uint32, _blockTimestamp uint64, _baseFeeConfig LibSharedDataBaseFeeConfig) (struct {
	Basefee      *big.Int
	NewGasTarget uint64
	NewGasExcess uint64
}, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "getBasefeeV2", _parentGasUsed, _blockTimestamp, _baseFeeConfig)

	outstruct := new(struct {
		Basefee      *big.Int
		NewGasTarget uint64
		NewGasExcess uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Basefee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NewGasTarget = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.NewGasExcess = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetBasefeeV2 is a free data retrieval call binding the contract method 0x893f5460.
//
// Solidity: function getBasefeeV2(uint32 _parentGasUsed, uint64 _blockTimestamp, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) view returns(uint256 basefee_, uint64 newGasTarget_, uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorSession) GetBasefeeV2(_parentGasUsed uint32, _blockTimestamp uint64, _baseFeeConfig LibSharedDataBaseFeeConfig) (struct {
	Basefee      *big.Int
	NewGasTarget uint64
	NewGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.GetBasefeeV2(&_TaikoAnchor.CallOpts, _parentGasUsed, _blockTimestamp, _baseFeeConfig)
}

// GetBasefeeV2 is a free data retrieval call binding the contract method 0x893f5460.
//
// Solidity: function getBasefeeV2(uint32 _parentGasUsed, uint64 _blockTimestamp, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) view returns(uint256 basefee_, uint64 newGasTarget_, uint64 newGasExcess_)
func (_TaikoAnchor *TaikoAnchorCallerSession) GetBasefeeV2(_parentGasUsed uint32, _blockTimestamp uint64, _baseFeeConfig LibSharedDataBaseFeeConfig) (struct {
	Basefee      *big.Int
	NewGasTarget uint64
	NewGasExcess uint64
}, error) {
	return _TaikoAnchor.Contract.GetBasefeeV2(&_TaikoAnchor.CallOpts, _parentGasUsed, _blockTimestamp, _baseFeeConfig)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 _blockId) view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCaller) GetBlockHash(opts *bind.CallOpts, _blockId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "getBlockHash", _blockId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 _blockId) view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorSession) GetBlockHash(_blockId *big.Int) ([32]byte, error) {
	return _TaikoAnchor.Contract.GetBlockHash(&_TaikoAnchor.CallOpts, _blockId)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 _blockId) view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCallerSession) GetBlockHash(_blockId *big.Int) ([32]byte, error) {
	return _TaikoAnchor.Contract.GetBlockHash(&_TaikoAnchor.CallOpts, _blockId)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) Impl() (common.Address, error) {
	return _TaikoAnchor.Contract.Impl(&_TaikoAnchor.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) Impl() (common.Address, error) {
	return _TaikoAnchor.Contract.Impl(&_TaikoAnchor.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoAnchor *TaikoAnchorCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoAnchor *TaikoAnchorSession) InNonReentrant() (bool, error) {
	return _TaikoAnchor.Contract.InNonReentrant(&_TaikoAnchor.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_TaikoAnchor *TaikoAnchorCallerSession) InNonReentrant() (bool, error) {
	return _TaikoAnchor.Contract.InNonReentrant(&_TaikoAnchor.CallOpts)
}

// IsOnL1 is a free data retrieval call binding the contract method 0xa4b23554.
//
// Solidity: function isOnL1() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorCaller) IsOnL1(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "isOnL1")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOnL1 is a free data retrieval call binding the contract method 0xa4b23554.
//
// Solidity: function isOnL1() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorSession) IsOnL1() (bool, error) {
	return _TaikoAnchor.Contract.IsOnL1(&_TaikoAnchor.CallOpts)
}

// IsOnL1 is a free data retrieval call binding the contract method 0xa4b23554.
//
// Solidity: function isOnL1() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorCallerSession) IsOnL1() (bool, error) {
	return _TaikoAnchor.Contract.IsOnL1(&_TaikoAnchor.CallOpts)
}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) L1ChainId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "l1ChainId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) L1ChainId() (uint64, error) {
	return _TaikoAnchor.Contract.L1ChainId(&_TaikoAnchor.CallOpts)
}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) L1ChainId() (uint64, error) {
	return _TaikoAnchor.Contract.L1ChainId(&_TaikoAnchor.CallOpts)
}

// LastSyncedBlock is a free data retrieval call binding the contract method 0x33d5ac9b.
//
// Solidity: function lastSyncedBlock() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) LastSyncedBlock(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "lastSyncedBlock")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastSyncedBlock is a free data retrieval call binding the contract method 0x33d5ac9b.
//
// Solidity: function lastSyncedBlock() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) LastSyncedBlock() (uint64, error) {
	return _TaikoAnchor.Contract.LastSyncedBlock(&_TaikoAnchor.CallOpts)
}

// LastSyncedBlock is a free data retrieval call binding the contract method 0x33d5ac9b.
//
// Solidity: function lastSyncedBlock() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) LastSyncedBlock() (uint64, error) {
	return _TaikoAnchor.Contract.LastSyncedBlock(&_TaikoAnchor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) Owner() (common.Address, error) {
	return _TaikoAnchor.Contract.Owner(&_TaikoAnchor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) Owner() (common.Address, error) {
	return _TaikoAnchor.Contract.Owner(&_TaikoAnchor.CallOpts)
}

// PacayaForkHeight is a free data retrieval call binding the contract method 0xba9f41e8.
//
// Solidity: function pacayaForkHeight() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) PacayaForkHeight(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "pacayaForkHeight")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PacayaForkHeight is a free data retrieval call binding the contract method 0xba9f41e8.
//
// Solidity: function pacayaForkHeight() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) PacayaForkHeight() (uint64, error) {
	return _TaikoAnchor.Contract.PacayaForkHeight(&_TaikoAnchor.CallOpts)
}

// PacayaForkHeight is a free data retrieval call binding the contract method 0xba9f41e8.
//
// Solidity: function pacayaForkHeight() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) PacayaForkHeight() (uint64, error) {
	return _TaikoAnchor.Contract.PacayaForkHeight(&_TaikoAnchor.CallOpts)
}

// ParentGasExcess is a free data retrieval call binding the contract method 0xb8c7b30c.
//
// Solidity: function parentGasExcess() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) ParentGasExcess(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "parentGasExcess")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ParentGasExcess is a free data retrieval call binding the contract method 0xb8c7b30c.
//
// Solidity: function parentGasExcess() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) ParentGasExcess() (uint64, error) {
	return _TaikoAnchor.Contract.ParentGasExcess(&_TaikoAnchor.CallOpts)
}

// ParentGasExcess is a free data retrieval call binding the contract method 0xb8c7b30c.
//
// Solidity: function parentGasExcess() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) ParentGasExcess() (uint64, error) {
	return _TaikoAnchor.Contract.ParentGasExcess(&_TaikoAnchor.CallOpts)
}

// ParentGasTarget is a free data retrieval call binding the contract method 0xa7137c0f.
//
// Solidity: function parentGasTarget() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) ParentGasTarget(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "parentGasTarget")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ParentGasTarget is a free data retrieval call binding the contract method 0xa7137c0f.
//
// Solidity: function parentGasTarget() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) ParentGasTarget() (uint64, error) {
	return _TaikoAnchor.Contract.ParentGasTarget(&_TaikoAnchor.CallOpts)
}

// ParentGasTarget is a free data retrieval call binding the contract method 0xa7137c0f.
//
// Solidity: function parentGasTarget() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) ParentGasTarget() (uint64, error) {
	return _TaikoAnchor.Contract.ParentGasTarget(&_TaikoAnchor.CallOpts)
}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCaller) ParentTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "parentTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorSession) ParentTimestamp() (uint64, error) {
	return _TaikoAnchor.Contract.ParentTimestamp(&_TaikoAnchor.CallOpts)
}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoAnchor *TaikoAnchorCallerSession) ParentTimestamp() (uint64, error) {
	return _TaikoAnchor.Contract.ParentTimestamp(&_TaikoAnchor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoAnchor *TaikoAnchorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoAnchor *TaikoAnchorSession) Paused() (bool, error) {
	return _TaikoAnchor.Contract.Paused(&_TaikoAnchor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TaikoAnchor *TaikoAnchorCallerSession) Paused() (bool, error) {
	return _TaikoAnchor.Contract.Paused(&_TaikoAnchor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) PendingOwner() (common.Address, error) {
	return _TaikoAnchor.Contract.PendingOwner(&_TaikoAnchor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) PendingOwner() (common.Address, error) {
	return _TaikoAnchor.Contract.PendingOwner(&_TaikoAnchor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorSession) ProxiableUUID() ([32]byte, error) {
	return _TaikoAnchor.Contract.ProxiableUUID(&_TaikoAnchor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _TaikoAnchor.Contract.ProxiableUUID(&_TaikoAnchor.CallOpts)
}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCaller) PublicInputHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "publicInputHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorSession) PublicInputHash() ([32]byte, error) {
	return _TaikoAnchor.Contract.PublicInputHash(&_TaikoAnchor.CallOpts)
}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoAnchor *TaikoAnchorCallerSession) PublicInputHash() ([32]byte, error) {
	return _TaikoAnchor.Contract.PublicInputHash(&_TaikoAnchor.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) Resolver() (common.Address, error) {
	return _TaikoAnchor.Contract.Resolver(&_TaikoAnchor.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) Resolver() (common.Address, error) {
	return _TaikoAnchor.Contract.Resolver(&_TaikoAnchor.CallOpts)
}

// SignalService is a free data retrieval call binding the contract method 0x62d09453.
//
// Solidity: function signalService() view returns(address)
func (_TaikoAnchor *TaikoAnchorCaller) SignalService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "signalService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SignalService is a free data retrieval call binding the contract method 0x62d09453.
//
// Solidity: function signalService() view returns(address)
func (_TaikoAnchor *TaikoAnchorSession) SignalService() (common.Address, error) {
	return _TaikoAnchor.Contract.SignalService(&_TaikoAnchor.CallOpts)
}

// SignalService is a free data retrieval call binding the contract method 0x62d09453.
//
// Solidity: function signalService() view returns(address)
func (_TaikoAnchor *TaikoAnchorCallerSession) SignalService() (common.Address, error) {
	return _TaikoAnchor.Contract.SignalService(&_TaikoAnchor.CallOpts)
}

// SkipFeeCheck is a free data retrieval call binding the contract method 0x2f980473.
//
// Solidity: function skipFeeCheck() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorCaller) SkipFeeCheck(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TaikoAnchor.contract.Call(opts, &out, "skipFeeCheck")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SkipFeeCheck is a free data retrieval call binding the contract method 0x2f980473.
//
// Solidity: function skipFeeCheck() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorSession) SkipFeeCheck() (bool, error) {
	return _TaikoAnchor.Contract.SkipFeeCheck(&_TaikoAnchor.CallOpts)
}

// SkipFeeCheck is a free data retrieval call binding the contract method 0x2f980473.
//
// Solidity: function skipFeeCheck() pure returns(bool)
func (_TaikoAnchor *TaikoAnchorCallerSession) SkipFeeCheck() (bool, error) {
	return _TaikoAnchor.Contract.SkipFeeCheck(&_TaikoAnchor.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoAnchor *TaikoAnchorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoAnchor *TaikoAnchorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AcceptOwnership(&_TaikoAnchor.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AcceptOwnership(&_TaikoAnchor.TransactOpts)
}

// Anchor is a paid mutator transaction binding the contract method 0xda69d3db.
//
// Solidity: function anchor(bytes32 _l1BlockHash, bytes32 _l1StateRoot, uint64 _l1BlockId, uint32 _parentGasUsed) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) Anchor(opts *bind.TransactOpts, _l1BlockHash [32]byte, _l1StateRoot [32]byte, _l1BlockId uint64, _parentGasUsed uint32) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "anchor", _l1BlockHash, _l1StateRoot, _l1BlockId, _parentGasUsed)
}

// Anchor is a paid mutator transaction binding the contract method 0xda69d3db.
//
// Solidity: function anchor(bytes32 _l1BlockHash, bytes32 _l1StateRoot, uint64 _l1BlockId, uint32 _parentGasUsed) returns()
func (_TaikoAnchor *TaikoAnchorSession) Anchor(_l1BlockHash [32]byte, _l1StateRoot [32]byte, _l1BlockId uint64, _parentGasUsed uint32) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Anchor(&_TaikoAnchor.TransactOpts, _l1BlockHash, _l1StateRoot, _l1BlockId, _parentGasUsed)
}

// Anchor is a paid mutator transaction binding the contract method 0xda69d3db.
//
// Solidity: function anchor(bytes32 _l1BlockHash, bytes32 _l1StateRoot, uint64 _l1BlockId, uint32 _parentGasUsed) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) Anchor(_l1BlockHash [32]byte, _l1StateRoot [32]byte, _l1BlockId uint64, _parentGasUsed uint32) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Anchor(&_TaikoAnchor.TransactOpts, _l1BlockHash, _l1StateRoot, _l1BlockId, _parentGasUsed)
}

// AnchorV2 is a paid mutator transaction binding the contract method 0xfd85eb2d.
//
// Solidity: function anchorV2(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) AnchorV2(opts *bind.TransactOpts, _anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "anchorV2", _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig)
}

// AnchorV2 is a paid mutator transaction binding the contract method 0xfd85eb2d.
//
// Solidity: function anchorV2(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) returns()
func (_TaikoAnchor *TaikoAnchorSession) AnchorV2(_anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AnchorV2(&_TaikoAnchor.TransactOpts, _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig)
}

// AnchorV2 is a paid mutator transaction binding the contract method 0xfd85eb2d.
//
// Solidity: function anchorV2(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) AnchorV2(_anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AnchorV2(&_TaikoAnchor.TransactOpts, _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig)
}

// AnchorV3 is a paid mutator transaction binding the contract method 0x48080a45.
//
// Solidity: function anchorV3(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, bytes32[] _signalSlots) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) AnchorV3(opts *bind.TransactOpts, _anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "anchorV3", _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig, _signalSlots)
}

// AnchorV3 is a paid mutator transaction binding the contract method 0x48080a45.
//
// Solidity: function anchorV3(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, bytes32[] _signalSlots) returns()
func (_TaikoAnchor *TaikoAnchorSession) AnchorV3(_anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AnchorV3(&_TaikoAnchor.TransactOpts, _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig, _signalSlots)
}

// AnchorV3 is a paid mutator transaction binding the contract method 0x48080a45.
//
// Solidity: function anchorV3(uint64 _anchorBlockId, bytes32 _anchorStateRoot, uint32 _parentGasUsed, (uint8,uint8,uint32,uint64,uint32) _baseFeeConfig, bytes32[] _signalSlots) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) AnchorV3(_anchorBlockId uint64, _anchorStateRoot [32]byte, _parentGasUsed uint32, _baseFeeConfig LibSharedDataBaseFeeConfig, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.AnchorV3(&_TaikoAnchor.TransactOpts, _anchorBlockId, _anchorStateRoot, _parentGasUsed, _baseFeeConfig, _signalSlots)
}

// Init is a paid mutator transaction binding the contract method 0xb310e9e9.
//
// Solidity: function init(address _owner, uint64 _l1ChainId, uint64 _initialGasExcess) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) Init(opts *bind.TransactOpts, _owner common.Address, _l1ChainId uint64, _initialGasExcess uint64) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "init", _owner, _l1ChainId, _initialGasExcess)
}

// Init is a paid mutator transaction binding the contract method 0xb310e9e9.
//
// Solidity: function init(address _owner, uint64 _l1ChainId, uint64 _initialGasExcess) returns()
func (_TaikoAnchor *TaikoAnchorSession) Init(_owner common.Address, _l1ChainId uint64, _initialGasExcess uint64) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Init(&_TaikoAnchor.TransactOpts, _owner, _l1ChainId, _initialGasExcess)
}

// Init is a paid mutator transaction binding the contract method 0xb310e9e9.
//
// Solidity: function init(address _owner, uint64 _l1ChainId, uint64 _initialGasExcess) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) Init(_owner common.Address, _l1ChainId uint64, _initialGasExcess uint64) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Init(&_TaikoAnchor.TransactOpts, _owner, _l1ChainId, _initialGasExcess)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoAnchor *TaikoAnchorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoAnchor *TaikoAnchorSession) Pause() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Pause(&_TaikoAnchor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) Pause() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Pause(&_TaikoAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoAnchor *TaikoAnchorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoAnchor *TaikoAnchorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.RenounceOwnership(&_TaikoAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.RenounceOwnership(&_TaikoAnchor.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoAnchor *TaikoAnchorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.TransferOwnership(&_TaikoAnchor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.TransferOwnership(&_TaikoAnchor.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoAnchor *TaikoAnchorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoAnchor *TaikoAnchorSession) Unpause() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Unpause(&_TaikoAnchor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) Unpause() (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Unpause(&_TaikoAnchor.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoAnchor *TaikoAnchorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.UpgradeTo(&_TaikoAnchor.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.UpgradeTo(&_TaikoAnchor.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoAnchor *TaikoAnchorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoAnchor *TaikoAnchorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.UpgradeToAndCall(&_TaikoAnchor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.UpgradeToAndCall(&_TaikoAnchor.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_TaikoAnchor *TaikoAnchorTransactor) Withdraw(opts *bind.TransactOpts, _token common.Address, _to common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.contract.Transact(opts, "withdraw", _token, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_TaikoAnchor *TaikoAnchorSession) Withdraw(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Withdraw(&_TaikoAnchor.TransactOpts, _token, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_TaikoAnchor *TaikoAnchorTransactorSession) Withdraw(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _TaikoAnchor.Contract.Withdraw(&_TaikoAnchor.TransactOpts, _token, _to)
}

// TaikoAnchorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TaikoAnchor contract.
type TaikoAnchorAdminChangedIterator struct {
	Event *TaikoAnchorAdminChanged // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorAdminChanged)
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
		it.Event = new(TaikoAnchorAdminChanged)
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
func (it *TaikoAnchorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorAdminChanged represents a AdminChanged event raised by the TaikoAnchor contract.
type TaikoAnchorAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TaikoAnchorAdminChangedIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorAdminChangedIterator{contract: _TaikoAnchor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TaikoAnchorAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorAdminChanged)
				if err := _TaikoAnchor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseAdminChanged(log types.Log) (*TaikoAnchorAdminChanged, error) {
	event := new(TaikoAnchorAdminChanged)
	if err := _TaikoAnchor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorAnchoredIterator is returned from FilterAnchored and is used to iterate over the raw logs and unpacked data for Anchored events raised by the TaikoAnchor contract.
type TaikoAnchorAnchoredIterator struct {
	Event *TaikoAnchorAnchored // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorAnchored)
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
		it.Event = new(TaikoAnchorAnchored)
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
func (it *TaikoAnchorAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorAnchored represents a Anchored event raised by the TaikoAnchor contract.
type TaikoAnchorAnchored struct {
	ParentHash      [32]byte
	ParentGasExcess uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAnchored is a free log retrieval operation binding the contract event 0x41c3f410f5c8ac36bb46b1dccef0de0f964087c9e688795fa02ecfa2c20b3fe4.
//
// Solidity: event Anchored(bytes32 parentHash, uint64 parentGasExcess)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterAnchored(opts *bind.FilterOpts) (*TaikoAnchorAnchoredIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "Anchored")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorAnchoredIterator{contract: _TaikoAnchor.contract, event: "Anchored", logs: logs, sub: sub}, nil
}

// WatchAnchored is a free log subscription operation binding the contract event 0x41c3f410f5c8ac36bb46b1dccef0de0f964087c9e688795fa02ecfa2c20b3fe4.
//
// Solidity: event Anchored(bytes32 parentHash, uint64 parentGasExcess)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchAnchored(opts *bind.WatchOpts, sink chan<- *TaikoAnchorAnchored) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "Anchored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorAnchored)
				if err := _TaikoAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
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

// ParseAnchored is a log parse operation binding the contract event 0x41c3f410f5c8ac36bb46b1dccef0de0f964087c9e688795fa02ecfa2c20b3fe4.
//
// Solidity: event Anchored(bytes32 parentHash, uint64 parentGasExcess)
func (_TaikoAnchor *TaikoAnchorFilterer) ParseAnchored(log types.Log) (*TaikoAnchorAnchored, error) {
	event := new(TaikoAnchorAnchored)
	if err := _TaikoAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the TaikoAnchor contract.
type TaikoAnchorBeaconUpgradedIterator struct {
	Event *TaikoAnchorBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorBeaconUpgraded)
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
		it.Event = new(TaikoAnchorBeaconUpgraded)
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
func (it *TaikoAnchorBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorBeaconUpgraded represents a BeaconUpgraded event raised by the TaikoAnchor contract.
type TaikoAnchorBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*TaikoAnchorBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorBeaconUpgradedIterator{contract: _TaikoAnchor.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *TaikoAnchorBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorBeaconUpgraded)
				if err := _TaikoAnchor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseBeaconUpgraded(log types.Log) (*TaikoAnchorBeaconUpgraded, error) {
	event := new(TaikoAnchorBeaconUpgraded)
	if err := _TaikoAnchor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorEIP1559UpdateIterator is returned from FilterEIP1559Update and is used to iterate over the raw logs and unpacked data for EIP1559Update events raised by the TaikoAnchor contract.
type TaikoAnchorEIP1559UpdateIterator struct {
	Event *TaikoAnchorEIP1559Update // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorEIP1559UpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorEIP1559Update)
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
		it.Event = new(TaikoAnchorEIP1559Update)
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
func (it *TaikoAnchorEIP1559UpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorEIP1559UpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorEIP1559Update represents a EIP1559Update event raised by the TaikoAnchor contract.
type TaikoAnchorEIP1559Update struct {
	OldGasTarget uint64
	NewGasTarget uint64
	OldGasExcess uint64
	NewGasExcess uint64
	Basefee      *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEIP1559Update is a free log retrieval operation binding the contract event 0x781ae5c2215806150d5c71a4ed5336e5dc3ad32aef04fc0f626a6ee0c2f8d1c8.
//
// Solidity: event EIP1559Update(uint64 oldGasTarget, uint64 newGasTarget, uint64 oldGasExcess, uint64 newGasExcess, uint256 basefee)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterEIP1559Update(opts *bind.FilterOpts) (*TaikoAnchorEIP1559UpdateIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "EIP1559Update")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorEIP1559UpdateIterator{contract: _TaikoAnchor.contract, event: "EIP1559Update", logs: logs, sub: sub}, nil
}

// WatchEIP1559Update is a free log subscription operation binding the contract event 0x781ae5c2215806150d5c71a4ed5336e5dc3ad32aef04fc0f626a6ee0c2f8d1c8.
//
// Solidity: event EIP1559Update(uint64 oldGasTarget, uint64 newGasTarget, uint64 oldGasExcess, uint64 newGasExcess, uint256 basefee)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchEIP1559Update(opts *bind.WatchOpts, sink chan<- *TaikoAnchorEIP1559Update) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "EIP1559Update")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorEIP1559Update)
				if err := _TaikoAnchor.contract.UnpackLog(event, "EIP1559Update", log); err != nil {
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

// ParseEIP1559Update is a log parse operation binding the contract event 0x781ae5c2215806150d5c71a4ed5336e5dc3ad32aef04fc0f626a6ee0c2f8d1c8.
//
// Solidity: event EIP1559Update(uint64 oldGasTarget, uint64 newGasTarget, uint64 oldGasExcess, uint64 newGasExcess, uint256 basefee)
func (_TaikoAnchor *TaikoAnchorFilterer) ParseEIP1559Update(log types.Log) (*TaikoAnchorEIP1559Update, error) {
	event := new(TaikoAnchorEIP1559Update)
	if err := _TaikoAnchor.contract.UnpackLog(event, "EIP1559Update", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TaikoAnchor contract.
type TaikoAnchorInitializedIterator struct {
	Event *TaikoAnchorInitialized // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorInitialized)
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
		it.Event = new(TaikoAnchorInitialized)
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
func (it *TaikoAnchorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorInitialized represents a Initialized event raised by the TaikoAnchor contract.
type TaikoAnchorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterInitialized(opts *bind.FilterOpts) (*TaikoAnchorInitializedIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorInitializedIterator{contract: _TaikoAnchor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TaikoAnchorInitialized) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorInitialized)
				if err := _TaikoAnchor.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseInitialized(log types.Log) (*TaikoAnchorInitialized, error) {
	event := new(TaikoAnchorInitialized)
	if err := _TaikoAnchor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the TaikoAnchor contract.
type TaikoAnchorOwnershipTransferStartedIterator struct {
	Event *TaikoAnchorOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorOwnershipTransferStarted)
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
		it.Event = new(TaikoAnchorOwnershipTransferStarted)
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
func (it *TaikoAnchorOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the TaikoAnchor contract.
type TaikoAnchorOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaikoAnchorOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorOwnershipTransferStartedIterator{contract: _TaikoAnchor.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *TaikoAnchorOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorOwnershipTransferStarted)
				if err := _TaikoAnchor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseOwnershipTransferStarted(log types.Log) (*TaikoAnchorOwnershipTransferStarted, error) {
	event := new(TaikoAnchorOwnershipTransferStarted)
	if err := _TaikoAnchor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TaikoAnchor contract.
type TaikoAnchorOwnershipTransferredIterator struct {
	Event *TaikoAnchorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorOwnershipTransferred)
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
		it.Event = new(TaikoAnchorOwnershipTransferred)
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
func (it *TaikoAnchorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorOwnershipTransferred represents a OwnershipTransferred event raised by the TaikoAnchor contract.
type TaikoAnchorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaikoAnchorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorOwnershipTransferredIterator{contract: _TaikoAnchor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TaikoAnchorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorOwnershipTransferred)
				if err := _TaikoAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseOwnershipTransferred(log types.Log) (*TaikoAnchorOwnershipTransferred, error) {
	event := new(TaikoAnchorOwnershipTransferred)
	if err := _TaikoAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the TaikoAnchor contract.
type TaikoAnchorPausedIterator struct {
	Event *TaikoAnchorPaused // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorPaused)
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
		it.Event = new(TaikoAnchorPaused)
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
func (it *TaikoAnchorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorPaused represents a Paused event raised by the TaikoAnchor contract.
type TaikoAnchorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterPaused(opts *bind.FilterOpts) (*TaikoAnchorPausedIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorPausedIterator{contract: _TaikoAnchor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *TaikoAnchorPaused) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorPaused)
				if err := _TaikoAnchor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParsePaused(log types.Log) (*TaikoAnchorPaused, error) {
	event := new(TaikoAnchorPaused)
	if err := _TaikoAnchor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the TaikoAnchor contract.
type TaikoAnchorUnpausedIterator struct {
	Event *TaikoAnchorUnpaused // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorUnpaused)
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
		it.Event = new(TaikoAnchorUnpaused)
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
func (it *TaikoAnchorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorUnpaused represents a Unpaused event raised by the TaikoAnchor contract.
type TaikoAnchorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*TaikoAnchorUnpausedIterator, error) {

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorUnpausedIterator{contract: _TaikoAnchor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *TaikoAnchorUnpaused) (event.Subscription, error) {

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorUnpaused)
				if err := _TaikoAnchor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseUnpaused(log types.Log) (*TaikoAnchorUnpaused, error) {
	event := new(TaikoAnchorUnpaused)
	if err := _TaikoAnchor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoAnchorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TaikoAnchor contract.
type TaikoAnchorUpgradedIterator struct {
	Event *TaikoAnchorUpgraded // Event containing the contract specifics and raw log

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
func (it *TaikoAnchorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoAnchorUpgraded)
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
		it.Event = new(TaikoAnchorUpgraded)
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
func (it *TaikoAnchorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoAnchorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoAnchorUpgraded represents a Upgraded event raised by the TaikoAnchor contract.
type TaikoAnchorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TaikoAnchor *TaikoAnchorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TaikoAnchorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TaikoAnchor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TaikoAnchorUpgradedIterator{contract: _TaikoAnchor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TaikoAnchor *TaikoAnchorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TaikoAnchorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TaikoAnchor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoAnchorUpgraded)
				if err := _TaikoAnchor.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_TaikoAnchor *TaikoAnchorFilterer) ParseUpgraded(log types.Log) (*TaikoAnchorUpgraded, error) {
	event := new(TaikoAnchorUpgraded)
	if err := _TaikoAnchor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
