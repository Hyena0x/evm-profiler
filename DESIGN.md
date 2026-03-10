# System Design Document: evm-profiler

## 1. Project Overview
`evm-profiler` is a high-performance, stateless Golang CLI tool designed to generate deep behavioral and security profiles for any EVM address. It aggregates data from multiple sources (RPC, Etherscan, GoPlus) concurrently to provide a comprehensive identity report in sub-second response times.

## 2. Tech Stack
- **Language**: Golang 1.25+
- **CLI Framework**: `spf13/cobra`
- **Web3 SDK**: `ethereum/go-ethereum` (for RPC interaction and address validation)
- **HTTP Client**: `go-resty/resty/v2`
- **Concurrency**: `golang.org/x/sync/errgroup` (for managed task orchestration)
- **Terminal UI**: `fatih/color`, `olekukonko/tablewriter`
- **Configuration**: `spf13/viper`

## 3. Architecture Design
The project follows a modular, clean-architecture approach to ensure extensibility, maintainability, and testability.

### Directory Structure
- `cmd/`: CLI command definitions and Cobra entry points.
- `internal/fetcher/`: Data acquisition layer (RPC, Etherscan, GoPlus) with built-in rate-limiting and retry logic.
- `internal/analyzer/`: Core orchestration engine that triggers fetchers concurrently and performs data cleansing/labeling.
- `internal/model/`: Shared domain models and data structures (e.g., `AddressProfile`).
- `internal/utils/`: High-precision math (Wei to Ether) and generic helper functions.
- `internal/printer/`: Formatting logic for terminal output and visual reports.

## 4. Key Analysis Dimensions
1. **Base Identity (RPC)**: 
   - Checksum validation.
   - Code detection (`eth_getCode`) to distinguish between EOA and Smart Contracts.
   - High-precision native token balance retrieval.
2. **Fund Origin (Etherscan)**: 
   - Tracing the initial transaction to identify the "Funder" address (the source of initial gas).
3. **Behavioral Profile (Etherscan)**: 
   - Statistical analysis of transaction history.
   - Automatic tagging based on activity (e.g., Whale, Active Bot, Early Adopter, DEX User).
4. **Security Audit (GoPlus API)**: 
   - Real-time risk assessment including blacklists, phishing involvement, and mixer interactions.

## 5. Technical Constraints & Performance Safeguards
- **Precision Safety**: **Absolute prohibition of `float64` for raw EVM amounts.** All amounts are handled via `math/big.Int` (Wei) and converted to `math/big.Float` (Ether) only at the presentation layer with 4-decimal precision.
- **Concurrency Orchestration**: Every network request must execute in a parallel goroutine managed by `errgroup`. Strict context-based timeouts (default 10s) are enforced across all IO operations.
- **Rate-Limiting & Resilience**: Integrated retry mechanisms for public/free APIs (e.g., Etherscan 5 QPS limit) to ensure stability under load.
- **Error Propagation**: Internal packages must return errors and never invoke `log.Fatal` or `os.Exit`. Final error handling and user notification are reserved for the CLI entry layer.
- **Statelessness**: No local database or persistent state is used other than basic configuration; analysis is calculated in real-time.