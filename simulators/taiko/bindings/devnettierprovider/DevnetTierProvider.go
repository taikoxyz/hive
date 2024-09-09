// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package devnettierprovider

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

// ITierProviderTier is an auto generated low-level Go binding around an user-defined struct.
type ITierProviderTier struct {
	VerifierName              [32]byte
	ValidityBond              *big.Int
	ContestBond               *big.Int
	CooldownWindow            *big.Int
	ProvingWindow             uint16
	MaxBlocksToVerifyPerProof uint8
}

// DevnetTierProviderMetaData contains all meta data concerning the DevnetTierProvider contract.
var DevnetTierProviderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BOND_UNIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"GRACE_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinTier\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getProvider\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTier\",\"inputs\":[{\"name\":\"_tierId\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITierProvider.Tier\",\"components\":[{\"name\":\"verifierName\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validityBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"contestBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"cooldownWindow\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"provingWindow\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxBlocksToVerifyPerProof\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTierIds\",\"inputs\":[],\"outputs\":[{\"name\":\"tiers_\",\"type\":\"uint16[]\",\"internalType\":\"uint16[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"TIER_NOT_FOUND\",\"inputs\":[]}]",
}

// DevnetTierProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use DevnetTierProviderMetaData.ABI instead.
var DevnetTierProviderABI = DevnetTierProviderMetaData.ABI

// DevnetTierProvider is an auto generated Go binding around an Ethereum contract.
type DevnetTierProvider struct {
	DevnetTierProviderCaller     // Read-only binding to the contract
	DevnetTierProviderTransactor // Write-only binding to the contract
	DevnetTierProviderFilterer   // Log filterer for contract events
}

// DevnetTierProviderCaller is an auto generated read-only Go binding around an Ethereum contract.
type DevnetTierProviderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierProviderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DevnetTierProviderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierProviderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DevnetTierProviderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierProviderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DevnetTierProviderSession struct {
	Contract     *DevnetTierProvider // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DevnetTierProviderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DevnetTierProviderCallerSession struct {
	Contract *DevnetTierProviderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DevnetTierProviderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DevnetTierProviderTransactorSession struct {
	Contract     *DevnetTierProviderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DevnetTierProviderRaw is an auto generated low-level Go binding around an Ethereum contract.
type DevnetTierProviderRaw struct {
	Contract *DevnetTierProvider // Generic contract binding to access the raw methods on
}

// DevnetTierProviderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DevnetTierProviderCallerRaw struct {
	Contract *DevnetTierProviderCaller // Generic read-only contract binding to access the raw methods on
}

// DevnetTierProviderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DevnetTierProviderTransactorRaw struct {
	Contract *DevnetTierProviderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDevnetTierProvider creates a new instance of DevnetTierProvider, bound to a specific deployed contract.
func NewDevnetTierProvider(address common.Address, backend bind.ContractBackend) (*DevnetTierProvider, error) {
	contract, err := bindDevnetTierProvider(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DevnetTierProvider{DevnetTierProviderCaller: DevnetTierProviderCaller{contract: contract}, DevnetTierProviderTransactor: DevnetTierProviderTransactor{contract: contract}, DevnetTierProviderFilterer: DevnetTierProviderFilterer{contract: contract}}, nil
}

// NewDevnetTierProviderCaller creates a new read-only instance of DevnetTierProvider, bound to a specific deployed contract.
func NewDevnetTierProviderCaller(address common.Address, caller bind.ContractCaller) (*DevnetTierProviderCaller, error) {
	contract, err := bindDevnetTierProvider(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DevnetTierProviderCaller{contract: contract}, nil
}

// NewDevnetTierProviderTransactor creates a new write-only instance of DevnetTierProvider, bound to a specific deployed contract.
func NewDevnetTierProviderTransactor(address common.Address, transactor bind.ContractTransactor) (*DevnetTierProviderTransactor, error) {
	contract, err := bindDevnetTierProvider(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DevnetTierProviderTransactor{contract: contract}, nil
}

// NewDevnetTierProviderFilterer creates a new log filterer instance of DevnetTierProvider, bound to a specific deployed contract.
func NewDevnetTierProviderFilterer(address common.Address, filterer bind.ContractFilterer) (*DevnetTierProviderFilterer, error) {
	contract, err := bindDevnetTierProvider(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DevnetTierProviderFilterer{contract: contract}, nil
}

// bindDevnetTierProvider binds a generic wrapper to an already deployed contract.
func bindDevnetTierProvider(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DevnetTierProviderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevnetTierProvider *DevnetTierProviderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DevnetTierProvider.Contract.DevnetTierProviderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevnetTierProvider *DevnetTierProviderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevnetTierProvider.Contract.DevnetTierProviderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevnetTierProvider *DevnetTierProviderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevnetTierProvider.Contract.DevnetTierProviderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevnetTierProvider *DevnetTierProviderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DevnetTierProvider.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevnetTierProvider *DevnetTierProviderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevnetTierProvider.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevnetTierProvider *DevnetTierProviderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevnetTierProvider.Contract.contract.Transact(opts, method, params...)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierProvider *DevnetTierProviderCaller) BONDUNIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "BOND_UNIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierProvider *DevnetTierProviderSession) BONDUNIT() (*big.Int, error) {
	return _DevnetTierProvider.Contract.BONDUNIT(&_DevnetTierProvider.CallOpts)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierProvider *DevnetTierProviderCallerSession) BONDUNIT() (*big.Int, error) {
	return _DevnetTierProvider.Contract.BONDUNIT(&_DevnetTierProvider.CallOpts)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderCaller) GRACEPERIOD(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "GRACE_PERIOD")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderSession) GRACEPERIOD() (uint16, error) {
	return _DevnetTierProvider.Contract.GRACEPERIOD(&_DevnetTierProvider.CallOpts)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderCallerSession) GRACEPERIOD() (uint16, error) {
	return _DevnetTierProvider.Contract.GRACEPERIOD(&_DevnetTierProvider.CallOpts)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderCaller) GetMinTier(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (uint16, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "getMinTier", arg0, arg1)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderSession) GetMinTier(arg0 common.Address, arg1 *big.Int) (uint16, error) {
	return _DevnetTierProvider.Contract.GetMinTier(&_DevnetTierProvider.CallOpts, arg0, arg1)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierProvider *DevnetTierProviderCallerSession) GetMinTier(arg0 common.Address, arg1 *big.Int) (uint16, error) {
	return _DevnetTierProvider.Contract.GetMinTier(&_DevnetTierProvider.CallOpts, arg0, arg1)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierProvider *DevnetTierProviderCaller) GetProvider(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "getProvider", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierProvider *DevnetTierProviderSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _DevnetTierProvider.Contract.GetProvider(&_DevnetTierProvider.CallOpts, arg0)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierProvider *DevnetTierProviderCallerSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _DevnetTierProvider.Contract.GetProvider(&_DevnetTierProvider.CallOpts, arg0)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierProvider *DevnetTierProviderCaller) GetTier(opts *bind.CallOpts, _tierId uint16) (ITierProviderTier, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "getTier", _tierId)

	if err != nil {
		return *new(ITierProviderTier), err
	}

	out0 := *abi.ConvertType(out[0], new(ITierProviderTier)).(*ITierProviderTier)

	return out0, err

}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierProvider *DevnetTierProviderSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _DevnetTierProvider.Contract.GetTier(&_DevnetTierProvider.CallOpts, _tierId)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierProvider *DevnetTierProviderCallerSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _DevnetTierProvider.Contract.GetTier(&_DevnetTierProvider.CallOpts, _tierId)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierProvider *DevnetTierProviderCaller) GetTierIds(opts *bind.CallOpts) ([]uint16, error) {
	var out []interface{}
	err := _DevnetTierProvider.contract.Call(opts, &out, "getTierIds")

	if err != nil {
		return *new([]uint16), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint16)).(*[]uint16)

	return out0, err

}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierProvider *DevnetTierProviderSession) GetTierIds() ([]uint16, error) {
	return _DevnetTierProvider.Contract.GetTierIds(&_DevnetTierProvider.CallOpts)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierProvider *DevnetTierProviderCallerSession) GetTierIds() ([]uint16, error) {
	return _DevnetTierProvider.Contract.GetTierIds(&_DevnetTierProvider.CallOpts)
}
