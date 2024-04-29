package contracts

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
)

const (
	Erc20AbiString string = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function"}]`
)

// Global container for the parsed ABI above
var erc20Abi *abi.ABI

// ==================
// === Interfaces ===
// ==================

// Simple binding for ERC20 tokens
type IErc20Token interface {
	// The address of the token
	Address() common.Address

	// Get the full name of the token
	Name() string

	// Get the token tracker symbol
	Symbol() string

	// Get the number of decimal places the token uses
	Decimals() uint8

	// The token balance of the given address
	BalanceOf(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address)

	// Transfer tokens to a different address
	Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error)
}

// ===============
// === Structs ===
// ===============

// Simple binding for ERC20 tokens
type Erc20Contract struct {
	name     string
	symbol   string
	decimals uint8
	contract *eth.Contract
	txMgr    *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a contract wrapper for the ERC20 at the given address
func NewErc20Contract(address common.Address, client eth.IExecutionClient, queryMgr *eth.QueryManager, txMgr *eth.TransactionManager, opts *bind.CallOpts) (*Erc20Contract, error) {
	// Parse the ABI
	if erc20Abi == nil {
		abiParsed, err := abi.JSON(strings.NewReader(Erc20AbiString))
		if err != nil {
			return nil, fmt.Errorf("error parsing ERC20 ABI: %w", err)
		}
		erc20Abi = &abiParsed
	}

	// Create contract
	contract := &eth.Contract{
		ContractImpl: bind.NewBoundContract(address, *erc20Abi, client, client, client),
		Address:      address,
		ABI:          erc20Abi,
	}

	// Create the wrapper
	wrapper := &Erc20Contract{
		contract: contract,
		txMgr:    txMgr,
	}

	// Get the details
	err := queryMgr.Query(func(mc *batch.MultiCaller) error {
		eth.AddCallToMulticaller(mc, contract, &wrapper.name, "name")
		eth.AddCallToMulticaller(mc, contract, &wrapper.symbol, "symbol")
		eth.AddCallToMulticaller(mc, contract, &wrapper.decimals, "decimals")
		return nil
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting ERC-20 details of token %s: %w", address.Hex(), err)
	}

	return wrapper, nil
}

// =============
// === Calls ===
// =============

// The address of the token
func (c *Erc20Contract) Address() common.Address {
	return c.contract.Address
}

// Get the full name of the token
func (c *Erc20Contract) Name() string {
	return c.name
}

// Get the token tracker symbol
func (c *Erc20Contract) Symbol() string {
	return c.symbol
}

// Get the number of decimal places the token uses
func (c *Erc20Contract) Decimals() uint8 {
	return c.decimals
}

// Get the token balance for an address
func (c *Erc20Contract) BalanceOf(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	eth.AddCallToMulticaller(mc, c.contract, balance_Out, "balanceOf", address)
}

// ====================
// === Transactions ===
// ====================

// Get info for transferring the ERC20 to another address
func (c *Erc20Contract) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.contract, "transfer", opts, to, amount)
}
