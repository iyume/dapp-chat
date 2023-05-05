// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chatABI

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

// ChatMetaData contains all meta data concerning the Chat contract.
var ChatMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"increase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101c4806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063209652551461003b57806330f3f0db14610059575b600080fd5b610043610075565b60405161005091906100b2565b60405180910390f35b610073600480360381019061006e91906100fe565b61007e565b005b60008054905090565b8060008082825461008f919061015a565b9250508190555050565b6000819050919050565b6100ac81610099565b82525050565b60006020820190506100c760008301846100a3565b92915050565b600080fd5b6100db81610099565b81146100e657600080fd5b50565b6000813590506100f8816100d2565b92915050565b600060208284031215610114576101136100cd565b5b6000610122848285016100e9565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061016582610099565b915061017083610099565b92508282019050808211156101885761018761012b565b5b9291505056fea2646970667358221220e624749545f5f1f9718af423f0b5695faeea92cda439f59779fbcbbb3aa8785d64736f6c63430008130033",
}

// ChatABI is the input ABI used to generate the binding from.
// Deprecated: Use ChatMetaData.ABI instead.
var ChatABI = ChatMetaData.ABI

// ChatBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChatMetaData.Bin instead.
var ChatBin = ChatMetaData.Bin

// DeployChat deploys a new Ethereum contract, binding an instance of Chat to it.
func DeployChat(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Chat, error) {
	parsed, err := ChatMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChatBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Chat{ChatCaller: ChatCaller{contract: contract}, ChatTransactor: ChatTransactor{contract: contract}, ChatFilterer: ChatFilterer{contract: contract}}, nil
}

// Chat is an auto generated Go binding around an Ethereum contract.
type Chat struct {
	ChatCaller     // Read-only binding to the contract
	ChatTransactor // Write-only binding to the contract
	ChatFilterer   // Log filterer for contract events
}

// ChatCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChatCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChatTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChatFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChatSession struct {
	Contract     *Chat             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChatCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChatCallerSession struct {
	Contract *ChatCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ChatTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChatTransactorSession struct {
	Contract     *ChatTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChatRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChatRaw struct {
	Contract *Chat // Generic contract binding to access the raw methods on
}

// ChatCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChatCallerRaw struct {
	Contract *ChatCaller // Generic read-only contract binding to access the raw methods on
}

// ChatTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChatTransactorRaw struct {
	Contract *ChatTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChat creates a new instance of Chat, bound to a specific deployed contract.
func NewChat(address common.Address, backend bind.ContractBackend) (*Chat, error) {
	contract, err := bindChat(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chat{ChatCaller: ChatCaller{contract: contract}, ChatTransactor: ChatTransactor{contract: contract}, ChatFilterer: ChatFilterer{contract: contract}}, nil
}

// NewChatCaller creates a new read-only instance of Chat, bound to a specific deployed contract.
func NewChatCaller(address common.Address, caller bind.ContractCaller) (*ChatCaller, error) {
	contract, err := bindChat(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChatCaller{contract: contract}, nil
}

// NewChatTransactor creates a new write-only instance of Chat, bound to a specific deployed contract.
func NewChatTransactor(address common.Address, transactor bind.ContractTransactor) (*ChatTransactor, error) {
	contract, err := bindChat(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChatTransactor{contract: contract}, nil
}

// NewChatFilterer creates a new log filterer instance of Chat, bound to a specific deployed contract.
func NewChatFilterer(address common.Address, filterer bind.ContractFilterer) (*ChatFilterer, error) {
	contract, err := bindChat(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChatFilterer{contract: contract}, nil
}

// bindChat binds a generic wrapper to an already deployed contract.
func bindChat(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChatMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chat *ChatRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chat.Contract.ChatCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chat *ChatRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chat.Contract.ChatTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chat *ChatRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chat.Contract.ChatTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chat *ChatCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chat.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chat *ChatTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chat.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chat *ChatTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chat.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Chat *ChatCaller) GetValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chat.contract.Call(opts, &out, "getValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Chat *ChatSession) GetValue() (*big.Int, error) {
	return _Chat.Contract.GetValue(&_Chat.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Chat *ChatCallerSession) GetValue() (*big.Int, error) {
	return _Chat.Contract.GetValue(&_Chat.CallOpts)
}

// Increase is a paid mutator transaction binding the contract method 0x30f3f0db.
//
// Solidity: function increase(uint256 a) returns()
func (_Chat *ChatTransactor) Increase(opts *bind.TransactOpts, a *big.Int) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "increase", a)
}

// Increase is a paid mutator transaction binding the contract method 0x30f3f0db.
//
// Solidity: function increase(uint256 a) returns()
func (_Chat *ChatSession) Increase(a *big.Int) (*types.Transaction, error) {
	return _Chat.Contract.Increase(&_Chat.TransactOpts, a)
}

// Increase is a paid mutator transaction binding the contract method 0x30f3f0db.
//
// Solidity: function increase(uint256 a) returns()
func (_Chat *ChatTransactorSession) Increase(a *big.Int) (*types.Transaction, error) {
	return _Chat.Contract.Increase(&_Chat.TransactOpts, a)
}
