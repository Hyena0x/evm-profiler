package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// EtherscanFetcher queries Etherscan APIs with retry features to prevent rate-limit dropping.
type EtherscanFetcher struct {
	client *resty.Client
	apiKey string
}

// NewEtherscanFetcher initializes an API client with resty.
func NewEtherscanFetcher(apiKey string) *EtherscanFetcher {
	client := resty.New().
		SetBaseURL("https://api.etherscan.io").
		SetTimeout(10 * time.Second).
		SetRetryCount(3). // Retries when 5 QPS limits are hit
		SetRetryWaitTime(1 * time.Second)
	return &EtherscanFetcher{
		client: client,
		apiKey: apiKey,
	}
}

type txListResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// GetFirstTx fetches the earliest normal transaction to find the funding address (Funder).
func (f *EtherscanFetcher) GetFirstTx(ctx context.Context, address string) (funder, hash string, err error) {
	var resp txListResponse
	_, err = f.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"chainid":    "1",
			"module":     "account",
			"action":     "txlist",
			"address":    address,
			"startblock": "0",
			"endblock":   "99999999",
			"page":       "1",
			"offset":     "1",
			"sort":       "asc",
			"apikey":     f.apiKey,
		}).
		SetResult(&resp).
		Get("/v2/api")

	if err != nil {
		return "", "", err
	}

	// '0' status means error, but message can be 'No transactions found' if it's a completely new address.
	if resp.Status != "1" && resp.Message != "No transactions found" {
		var errMsg string
		if len(resp.Result) > 0 && resp.Result[0] == '"' {
			json.Unmarshal(resp.Result, &errMsg)
		} else {
			errMsg = string(resp.Result)
		}
		return "", "", fmt.Errorf("etherscan err: %s, %s", resp.Message, errMsg)
	}

	if resp.Status == "1" {
		var txs []struct {
			From string `json:"from"`
			Hash string `json:"hash"`
		}
		if err := json.Unmarshal(resp.Result, &txs); err != nil {
			return "", "", fmt.Errorf("etherscan decode err: %w", err)
		}
		if len(txs) > 0 {
			return txs[0].From, txs[0].Hash, nil
		}
	}

	return "", "", nil
}

// GetRecentTxs retrieves up to the last 100 transactions to assess activity frequency.
func (f *EtherscanFetcher) GetRecentTxs(ctx context.Context, address string) (count int, err error) {
	var resp txListResponse
	_, err = f.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"chainid":    "1",
			"module":     "account",
			"action":     "txlist",
			"address":    address,
			"startblock": "0",
			"endblock":   "99999999",
			"page":       "1",
			"offset":     "100",
			"sort":       "desc",
			"apikey":     f.apiKey,
		}).
		SetResult(&resp).
		Get("/v2/api")

	if err != nil {
		return 0, err
	}

	if resp.Status != "1" && resp.Message != "No transactions found" {
		var errMsg string
		if len(resp.Result) > 0 && resp.Result[0] == '"' {
			json.Unmarshal(resp.Result, &errMsg)
		} else {
			errMsg = string(resp.Result)
		}
		return 0, fmt.Errorf("etherscan err: %s, %s", resp.Message, errMsg)
	}

	if resp.Status == "1" {
		var txs []struct {
			From string `json:"from"`
			Hash string `json:"hash"`
		}
		if err := json.Unmarshal(resp.Result, &txs); err != nil {
			return 0, fmt.Errorf("etherscan decode err: %w", err)
		}
		return len(txs), nil
	}

	return 0, nil
}
