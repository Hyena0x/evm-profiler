# evm-profiler

[![Go Version](https://img.shields.io/github/go-mod/go-version/Hyena0x/evm-profiler)](https://golang.org/)
[![License](https://img.shields.io/github/license/Hyena0x/evm-profiler)](./LICENSE)

A high-performance, stateless CLI tool for deep EVM address profiling.

## 🚀 Overview

`evm-profiler` is a professional-grade Golang CLI tool designed for Web3 researchers, security auditors, and developers. It generates a comprehensive "identity report" of any EVM address by concurrently aggregating data from RPC nodes, Etherscan, and GoPlus security APIs.

## ✨ Key Features

- **🔍 Smart Detection**: Instant identification of EOA vs. Smart Contract with checksum validation.
- **💰 Balance Tracking**: High-precision native token balance retrieval using `math/big`.
- **🧬 Fund Origin**: Automatically traces the address back to its first transaction to identify the initial gas source (Funder).
- **📊 Behavioral Profiling**: Analyzes transaction history to tag addresses (e.g., Whale, Active Bot, Early Adopter).
- **🛡️ Security Audit**: Integrated risk assessment (blacklists, phishing, mixer interactions) via GoPlus API.
- **⚡ High Concurrency**: Orchestrates multiple network requests in parallel using `errgroup` for sub-second response times.
- **🎨 Premium UX**: Beautifully formatted terminal tables with colorized risk indicators.

## 📥 Installation

Ensure you have Go 1.25+ installed:

```bash
go install github.com/Hyena0x/evm-profiler@latest
```

Alternatively, clone and build manually:

```bash
git clone https://github.com/Hyena0x/evm-profiler.git
cd evm-profiler
go build -o evm-profiler
```

## ⚙️ Configuration

Start by setting up your RPC and API keys to avoid rate limits:

```bash
# Set a public Ethereum RPC (no key needed for basic tests)
evm-profiler config set --rpc https://ethereum-rpc.publicnode.com

# Or set a private provider like Alchemy (requires key)
evm-profiler config set --rpc https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY
```

*Configuration is stored locally via Viper.*

## 🛠️ Usage

Analyze any address with a single command (example uses Vitalik's public address):

```bash
evm-profiler analyze 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

### Flags

| Flag | Description |
|------|-------------|
| `--rpc` | Override default RPC URL |
| `--etherscan-key` | Provide Etherscan API key |
| `--network` | Specify network (default: mainnet) |

## 📐 Architecture

This project follows a modular design for maximum extensibility:
- `cmd/`: CLI entry points (Cobra).
- `internal/fetcher/`: External data providers (RPC, Etherscan, GoPlus).
- `internal/analyzer/`: Concurrent orchestration engine.
- `internal/printer/`: Visual report generation.

For deep technical details, check the [DESIGN.md](./DESIGN.md).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
