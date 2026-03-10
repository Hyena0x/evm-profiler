package analyzer

import (
	"context"
	"fmt"
	"os"
	"time"

	"evm-profiler/internal/fetcher"
	"evm-profiler/internal/model"
	"evm-profiler/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

// Analyzer coordinates data parallel fetching using errgroup and constructs the Profile.
type Analyzer struct {
	rpcFetcher       *fetcher.RPCFetcher
	etherscanFetcher *fetcher.EtherscanFetcher
	goplusFetcher    *fetcher.GoPlusFetcher
}

// NewAnalyzer creates an analyzer bound with individual feature fetchers.
func NewAnalyzer(rpc *fetcher.RPCFetcher, etherscan *fetcher.EtherscanFetcher, goplus *fetcher.GoPlusFetcher) *Analyzer {
	return &Analyzer{
		rpcFetcher:       rpc,
		etherscanFetcher: etherscan,
		goplusFetcher:    goplus,
	}
}

// Analyze triggers a concurrent burst of network checks.
func (a *Analyzer) Analyze(ctx context.Context, address string) (*model.AddressProfile, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid EVM address: %s", address)
	}

	// Make sure we output the EIP-55 Checksum address natively
	addr := common.HexToAddress(address).Hex()

	profile := &model.AddressProfile{
		Address:     addr,
		LastUpdated: time.Now(),
	}

	// errgroup helps bound all goroutines nicely; returns the first non-nil err if it happens.
	g, groupCtx := errgroup.WithContext(ctx)

	// Sub-task: Basic RPC data
	g.Go(func() error {
		balance, isContract, err := a.rpcFetcher.GetBasicInfo(groupCtx, addr)
		if err != nil {
			return fmt.Errorf("Failed to fetch account basic info via RPC: %w", err)
		}
		profile.Balance = balance
		profile.BalanceEther = utils.WeiToEther(balance)
		profile.IsContract = isContract
		return nil
	})

	// Sub-task: Etherscan First Tx
	g.Go(func() error {
		funder, hash, err := a.etherscanFetcher.GetFirstTx(groupCtx, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ [Warning] Failed to fetch Etherscan First Tx: %v\n", err)
			return nil
		}
		profile.Funder = funder
		profile.FirstTxHash = hash
		return nil
	})

	// Sub-task: Etherscan Recent Tx Volume (Behavior Profiling)
	g.Go(func() error {
		count, err := a.etherscanFetcher.GetRecentTxs(groupCtx, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ [Warning] Failed to fetch Etherscan Tx Volume: %v\n", err)
			return nil
		}
		profile.TransactionCount = count

		// Simple behavioral label heuristic
		switch {
		case count >= 100:
			profile.ActivityLabel = "Highly Active User / Bot"
		case count > 10:
			profile.ActivityLabel = "Active User"
		case count > 0:
			profile.ActivityLabel = "Occasional User"
		default:
			profile.ActivityLabel = "Inactive / New"
		}
		return nil
	})

	// Sub-task: GoPlus Security Audit
	g.Go(func() error {
		risks, err := a.goplusFetcher.GetSecurityRisks(groupCtx, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ [Warning] Failed to fetch GoPlus Security Risks: %v\n", err)
			return nil
		}
		profile.RiskLabels = risks
		return nil
	})

	// Blocks until all goroutines resolve or error out
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("analysis encountered error: %w", err)
	}

	return profile, nil
}
