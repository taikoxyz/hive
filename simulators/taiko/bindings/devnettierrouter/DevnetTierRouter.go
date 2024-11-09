// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package devnettierrouter

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

// DevnetTierRouterMetaData contains all meta data concerning the DevnetTierRouter contract.
var DevnetTierRouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BOND_UNIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinTier\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getProvider\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTier\",\"inputs\":[{\"name\":\"_tierId\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITierProvider.Tier\",\"components\":[{\"name\":\"verifierName\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validityBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"contestBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"cooldownWindow\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"provingWindow\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxBlocksToVerifyPerProof\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTierIds\",\"inputs\":[],\"outputs\":[{\"name\":\"tiers_\",\"type\":\"uint16[]\",\"internalType\":\"uint16[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"TIER_NOT_FOUND\",\"inputs\":[]}]",
}

// DevnetTierRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use DevnetTierRouterMetaData.ABI instead.
var DevnetTierRouterABI = DevnetTierRouterMetaData.ABI

// DevnetTierRouter is an auto generated Go binding around an Ethereum contract.
type DevnetTierRouter struct {
	DevnetTierRouterCaller     // Read-only binding to the contract
	DevnetTierRouterTransactor // Write-only binding to the contract
	DevnetTierRouterFilterer   // Log filterer for contract events
}

// DevnetTierRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type DevnetTierRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DevnetTierRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DevnetTierRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevnetTierRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DevnetTierRouterSession struct {
	Contract     *DevnetTierRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DevnetTierRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DevnetTierRouterCallerSession struct {
	Contract *DevnetTierRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DevnetTierRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DevnetTierRouterTransactorSession struct {
	Contract     *DevnetTierRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DevnetTierRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type DevnetTierRouterRaw struct {
	Contract *DevnetTierRouter // Generic contract binding to access the raw methods on
}

// DevnetTierRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DevnetTierRouterCallerRaw struct {
	Contract *DevnetTierRouterCaller // Generic read-only contract binding to access the raw methods on
}

// DevnetTierRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DevnetTierRouterTransactorRaw struct {
	Contract *DevnetTierRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDevnetTierRouter creates a new instance of DevnetTierRouter, bound to a specific deployed contract.
func NewDevnetTierRouter(address common.Address, backend bind.ContractBackend) (*DevnetTierRouter, error) {
	contract, err := bindDevnetTierRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DevnetTierRouter{DevnetTierRouterCaller: DevnetTierRouterCaller{contract: contract}, DevnetTierRouterTransactor: DevnetTierRouterTransactor{contract: contract}, DevnetTierRouterFilterer: DevnetTierRouterFilterer{contract: contract}}, nil
}

// NewDevnetTierRouterCaller creates a new read-only instance of DevnetTierRouter, bound to a specific deployed contract.
func NewDevnetTierRouterCaller(address common.Address, caller bind.ContractCaller) (*DevnetTierRouterCaller, error) {
	contract, err := bindDevnetTierRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DevnetTierRouterCaller{contract: contract}, nil
}

// NewDevnetTierRouterTransactor creates a new write-only instance of DevnetTierRouter, bound to a specific deployed contract.
func NewDevnetTierRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*DevnetTierRouterTransactor, error) {
	contract, err := bindDevnetTierRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DevnetTierRouterTransactor{contract: contract}, nil
}

// NewDevnetTierRouterFilterer creates a new log filterer instance of DevnetTierRouter, bound to a specific deployed contract.
func NewDevnetTierRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*DevnetTierRouterFilterer, error) {
	contract, err := bindDevnetTierRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DevnetTierRouterFilterer{contract: contract}, nil
}

// bindDevnetTierRouter binds a generic wrapper to an already deployed contract.
func bindDevnetTierRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DevnetTierRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevnetTierRouter *DevnetTierRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DevnetTierRouter.Contract.DevnetTierRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevnetTierRouter *DevnetTierRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevnetTierRouter.Contract.DevnetTierRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevnetTierRouter *DevnetTierRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevnetTierRouter.Contract.DevnetTierRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevnetTierRouter *DevnetTierRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DevnetTierRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevnetTierRouter *DevnetTierRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevnetTierRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevnetTierRouter *DevnetTierRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevnetTierRouter.Contract.contract.Transact(opts, method, params...)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierRouter *DevnetTierRouterCaller) BONDUNIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DevnetTierRouter.contract.Call(opts, &out, "BOND_UNIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierRouter *DevnetTierRouterSession) BONDUNIT() (*big.Int, error) {
	return _DevnetTierRouter.Contract.BONDUNIT(&_DevnetTierRouter.CallOpts)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_DevnetTierRouter *DevnetTierRouterCallerSession) BONDUNIT() (*big.Int, error) {
	return _DevnetTierRouter.Contract.BONDUNIT(&_DevnetTierRouter.CallOpts)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierRouter *DevnetTierRouterCaller) GetMinTier(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (uint16, error) {
	var out []interface{}
	err := _DevnetTierRouter.contract.Call(opts, &out, "getMinTier", arg0, arg1)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierRouter *DevnetTierRouterSession) GetMinTier(arg0 common.Address, arg1 *big.Int) (uint16, error) {
	return _DevnetTierRouter.Contract.GetMinTier(&_DevnetTierRouter.CallOpts, arg0, arg1)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address , uint256 ) pure returns(uint16)
func (_DevnetTierRouter *DevnetTierRouterCallerSession) GetMinTier(arg0 common.Address, arg1 *big.Int) (uint16, error) {
	return _DevnetTierRouter.Contract.GetMinTier(&_DevnetTierRouter.CallOpts, arg0, arg1)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierRouter *DevnetTierRouterCaller) GetProvider(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DevnetTierRouter.contract.Call(opts, &out, "getProvider", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierRouter *DevnetTierRouterSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _DevnetTierRouter.Contract.GetProvider(&_DevnetTierRouter.CallOpts, arg0)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_DevnetTierRouter *DevnetTierRouterCallerSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _DevnetTierRouter.Contract.GetProvider(&_DevnetTierRouter.CallOpts, arg0)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierRouter *DevnetTierRouterCaller) GetTier(opts *bind.CallOpts, _tierId uint16) (ITierProviderTier, error) {
	var out []interface{}
	err := _DevnetTierRouter.contract.Call(opts, &out, "getTier", _tierId)

	if err != nil {
		return *new(ITierProviderTier), err
	}

	out0 := *abi.ConvertType(out[0], new(ITierProviderTier)).(*ITierProviderTier)

	return out0, err

}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierRouter *DevnetTierRouterSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _DevnetTierRouter.Contract.GetTier(&_DevnetTierRouter.CallOpts, _tierId)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_DevnetTierRouter *DevnetTierRouterCallerSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _DevnetTierRouter.Contract.GetTier(&_DevnetTierRouter.CallOpts, _tierId)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierRouter *DevnetTierRouterCaller) GetTierIds(opts *bind.CallOpts) ([]uint16, error) {
	var out []interface{}
	err := _DevnetTierRouter.contract.Call(opts, &out, "getTierIds")

	if err != nil {
		return *new([]uint16), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint16)).(*[]uint16)

	return out0, err

}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierRouter *DevnetTierRouterSession) GetTierIds() ([]uint16, error) {
	return _DevnetTierRouter.Contract.GetTierIds(&_DevnetTierRouter.CallOpts)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_DevnetTierRouter *DevnetTierRouterCallerSession) GetTierIds() ([]uint16, error) {
	return _DevnetTierRouter.Contract.GetTierIds(&_DevnetTierRouter.CallOpts)
}
