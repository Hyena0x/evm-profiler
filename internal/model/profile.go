package model

import (
	"math/big"
	"time"
)

// AddressProfile is the core struct representing the fetched profile of an EVM address.
type AddressProfile struct {
	Address          string         // The exact queried address (checksummed)
	IsContract       bool           // Whether the address is an EOA or a Contract (based on eth_getCode)
	Balance          *big.Int       // Raw token balance in Wei
	BalanceEther     *big.Float     // Processed balance in Ether (retaining 4 decimal places for UI)
	Funder           string         // The funding address of the first tx
	FirstTxHash      string         // Hash of the funding tx
	TransactionCount int            // Transaction count in recent history (e.g., last 100 tx)
	ActivityLabel    string         // Profile label. E.g., Whale, Bot, Normal User
	RiskLabels       []string       // Any security risks reported by GoPlus API
	LastUpdated      time.Time      // Time the profile was generated
}
