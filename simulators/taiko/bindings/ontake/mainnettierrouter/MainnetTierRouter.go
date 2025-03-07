// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mainnettierrouter

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

// MainnetTierRouterMetaData contains all meta data concerning the MainnetTierRouter contract.
var MainnetTierRouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_daoFallbackProposer\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BOND_UNIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DAO_FALLBACK_PROPOSER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinTier\",\"inputs\":[{\"name\":\"_proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rand\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProvider\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTier\",\"inputs\":[{\"name\":\"_tierId\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structITierProvider.Tier\",\"components\":[{\"name\":\"verifierName\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validityBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"contestBond\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"cooldownWindow\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"provingWindow\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxBlocksToVerifyPerProof\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTierIds\",\"inputs\":[],\"outputs\":[{\"name\":\"tiers_\",\"type\":\"uint16[]\",\"internalType\":\"uint16[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"TIER_NOT_FOUND\",\"inputs\":[]}]",
}

// MainnetTierRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use MainnetTierRouterMetaData.ABI instead.
var MainnetTierRouterABI = MainnetTierRouterMetaData.ABI

// MainnetTierRouter is an auto generated Go binding around an Ethereum contract.
type MainnetTierRouter struct {
	MainnetTierRouterCaller     // Read-only binding to the contract
	MainnetTierRouterTransactor // Write-only binding to the contract
	MainnetTierRouterFilterer   // Log filterer for contract events
}

// MainnetTierRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type MainnetTierRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainnetTierRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MainnetTierRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainnetTierRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MainnetTierRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainnetTierRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MainnetTierRouterSession struct {
	Contract     *MainnetTierRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MainnetTierRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MainnetTierRouterCallerSession struct {
	Contract *MainnetTierRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MainnetTierRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MainnetTierRouterTransactorSession struct {
	Contract     *MainnetTierRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MainnetTierRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type MainnetTierRouterRaw struct {
	Contract *MainnetTierRouter // Generic contract binding to access the raw methods on
}

// MainnetTierRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MainnetTierRouterCallerRaw struct {
	Contract *MainnetTierRouterCaller // Generic read-only contract binding to access the raw methods on
}

// MainnetTierRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MainnetTierRouterTransactorRaw struct {
	Contract *MainnetTierRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMainnetTierRouter creates a new instance of MainnetTierRouter, bound to a specific deployed contract.
func NewMainnetTierRouter(address common.Address, backend bind.ContractBackend) (*MainnetTierRouter, error) {
	contract, err := bindMainnetTierRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MainnetTierRouter{MainnetTierRouterCaller: MainnetTierRouterCaller{contract: contract}, MainnetTierRouterTransactor: MainnetTierRouterTransactor{contract: contract}, MainnetTierRouterFilterer: MainnetTierRouterFilterer{contract: contract}}, nil
}

// NewMainnetTierRouterCaller creates a new read-only instance of MainnetTierRouter, bound to a specific deployed contract.
func NewMainnetTierRouterCaller(address common.Address, caller bind.ContractCaller) (*MainnetTierRouterCaller, error) {
	contract, err := bindMainnetTierRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MainnetTierRouterCaller{contract: contract}, nil
}

// NewMainnetTierRouterTransactor creates a new write-only instance of MainnetTierRouter, bound to a specific deployed contract.
func NewMainnetTierRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*MainnetTierRouterTransactor, error) {
	contract, err := bindMainnetTierRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MainnetTierRouterTransactor{contract: contract}, nil
}

// NewMainnetTierRouterFilterer creates a new log filterer instance of MainnetTierRouter, bound to a specific deployed contract.
func NewMainnetTierRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*MainnetTierRouterFilterer, error) {
	contract, err := bindMainnetTierRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MainnetTierRouterFilterer{contract: contract}, nil
}

// bindMainnetTierRouter binds a generic wrapper to an already deployed contract.
func bindMainnetTierRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MainnetTierRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MainnetTierRouter *MainnetTierRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MainnetTierRouter.Contract.MainnetTierRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MainnetTierRouter *MainnetTierRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MainnetTierRouter.Contract.MainnetTierRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MainnetTierRouter *MainnetTierRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MainnetTierRouter.Contract.MainnetTierRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MainnetTierRouter *MainnetTierRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MainnetTierRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MainnetTierRouter *MainnetTierRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MainnetTierRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MainnetTierRouter *MainnetTierRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MainnetTierRouter.Contract.contract.Transact(opts, method, params...)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_MainnetTierRouter *MainnetTierRouterCaller) BONDUNIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "BOND_UNIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_MainnetTierRouter *MainnetTierRouterSession) BONDUNIT() (*big.Int, error) {
	return _MainnetTierRouter.Contract.BONDUNIT(&_MainnetTierRouter.CallOpts)
}

// BONDUNIT is a free data retrieval call binding the contract method 0x8165fd26.
//
// Solidity: function BOND_UNIT() view returns(uint96)
func (_MainnetTierRouter *MainnetTierRouterCallerSession) BONDUNIT() (*big.Int, error) {
	return _MainnetTierRouter.Contract.BONDUNIT(&_MainnetTierRouter.CallOpts)
}

// DAOFALLBACKPROPOSER is a free data retrieval call binding the contract method 0xbf62514d.
//
// Solidity: function DAO_FALLBACK_PROPOSER() view returns(address)
func (_MainnetTierRouter *MainnetTierRouterCaller) DAOFALLBACKPROPOSER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "DAO_FALLBACK_PROPOSER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DAOFALLBACKPROPOSER is a free data retrieval call binding the contract method 0xbf62514d.
//
// Solidity: function DAO_FALLBACK_PROPOSER() view returns(address)
func (_MainnetTierRouter *MainnetTierRouterSession) DAOFALLBACKPROPOSER() (common.Address, error) {
	return _MainnetTierRouter.Contract.DAOFALLBACKPROPOSER(&_MainnetTierRouter.CallOpts)
}

// DAOFALLBACKPROPOSER is a free data retrieval call binding the contract method 0xbf62514d.
//
// Solidity: function DAO_FALLBACK_PROPOSER() view returns(address)
func (_MainnetTierRouter *MainnetTierRouterCallerSession) DAOFALLBACKPROPOSER() (common.Address, error) {
	return _MainnetTierRouter.Contract.DAOFALLBACKPROPOSER(&_MainnetTierRouter.CallOpts)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address _proposer, uint256 _rand) view returns(uint16)
func (_MainnetTierRouter *MainnetTierRouterCaller) GetMinTier(opts *bind.CallOpts, _proposer common.Address, _rand *big.Int) (uint16, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "getMinTier", _proposer, _rand)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address _proposer, uint256 _rand) view returns(uint16)
func (_MainnetTierRouter *MainnetTierRouterSession) GetMinTier(_proposer common.Address, _rand *big.Int) (uint16, error) {
	return _MainnetTierRouter.Contract.GetMinTier(&_MainnetTierRouter.CallOpts, _proposer, _rand)
}

// GetMinTier is a free data retrieval call binding the contract method 0x52c5c56b.
//
// Solidity: function getMinTier(address _proposer, uint256 _rand) view returns(uint16)
func (_MainnetTierRouter *MainnetTierRouterCallerSession) GetMinTier(_proposer common.Address, _rand *big.Int) (uint16, error) {
	return _MainnetTierRouter.Contract.GetMinTier(&_MainnetTierRouter.CallOpts, _proposer, _rand)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_MainnetTierRouter *MainnetTierRouterCaller) GetProvider(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "getProvider", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_MainnetTierRouter *MainnetTierRouterSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _MainnetTierRouter.Contract.GetProvider(&_MainnetTierRouter.CallOpts, arg0)
}

// GetProvider is a free data retrieval call binding the contract method 0x5c42d079.
//
// Solidity: function getProvider(uint256 ) view returns(address)
func (_MainnetTierRouter *MainnetTierRouterCallerSession) GetProvider(arg0 *big.Int) (common.Address, error) {
	return _MainnetTierRouter.Contract.GetProvider(&_MainnetTierRouter.CallOpts, arg0)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_MainnetTierRouter *MainnetTierRouterCaller) GetTier(opts *bind.CallOpts, _tierId uint16) (ITierProviderTier, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "getTier", _tierId)

	if err != nil {
		return *new(ITierProviderTier), err
	}

	out0 := *abi.ConvertType(out[0], new(ITierProviderTier)).(*ITierProviderTier)

	return out0, err

}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_MainnetTierRouter *MainnetTierRouterSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _MainnetTierRouter.Contract.GetTier(&_MainnetTierRouter.CallOpts, _tierId)
}

// GetTier is a free data retrieval call binding the contract method 0x576c3de7.
//
// Solidity: function getTier(uint16 _tierId) pure returns((bytes32,uint96,uint96,uint24,uint16,uint8))
func (_MainnetTierRouter *MainnetTierRouterCallerSession) GetTier(_tierId uint16) (ITierProviderTier, error) {
	return _MainnetTierRouter.Contract.GetTier(&_MainnetTierRouter.CallOpts, _tierId)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_MainnetTierRouter *MainnetTierRouterCaller) GetTierIds(opts *bind.CallOpts) ([]uint16, error) {
	var out []interface{}
	err := _MainnetTierRouter.contract.Call(opts, &out, "getTierIds")

	if err != nil {
		return *new([]uint16), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint16)).(*[]uint16)

	return out0, err

}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_MainnetTierRouter *MainnetTierRouterSession) GetTierIds() ([]uint16, error) {
	return _MainnetTierRouter.Contract.GetTierIds(&_MainnetTierRouter.CallOpts)
}

// GetTierIds is a free data retrieval call binding the contract method 0xd8cde1c6.
//
// Solidity: function getTierIds() pure returns(uint16[] tiers_)
func (_MainnetTierRouter *MainnetTierRouterCallerSession) GetTierIds() ([]uint16, error) {
	return _MainnetTierRouter.Contract.GetTierIds(&_MainnetTierRouter.CallOpts)
}
