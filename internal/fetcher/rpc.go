package fetcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RPCFetcher retrieves basic blockchain states using Ethereum JSON-RPC.
type RPCFetcher struct {
	client *ethclient.Client
}

// NewRPCFetcher initializes the RPC client with the given endpoint.
func NewRPCFetcher(rpcURL string) (*RPCFetcher, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &RPCFetcher{client: client}, nil
}

// GetBasicInfo fetches the raw Wei balance and checks if the given address is a smart contract.
func (f *RPCFetcher) GetBasicInfo(ctx context.Context, address string) (*big.Int, bool, error) {
	addr := common.HexToAddress(address)

	// Fetch Balance (in Wei)
	balance, err := f.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, false, err
	}

	// Fetch Code (if length > 0, it means the address has deployed code -> contract)
	code, err := f.client.CodeAt(ctx, addr, nil)
	if err != nil {
		return nil, false, err
	}

	isContract := len(code) > 0

	return balance, isContract, nil
}
