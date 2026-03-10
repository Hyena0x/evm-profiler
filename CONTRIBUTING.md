# Contributing to evm-profiler

Welcome to `evm-profiler`! I'm excited to have your help. This document provides some basic guidelines for contributing to this project.

## 🛠️ Development Setup

To contribute to `evm-profiler`, you'll need:
- [Go 1.25+](https://golang.org/)
- An Ethereum RPC URL (e.g., [Alchemy](https://alchemy.com/), [Infura](https://infura.io/))
- Optional: [Etherscan API Key](https://etherscan.io/apis)

### Setup Steps
1. Fork the repository.
2. Clone your fork:
   ```bash
   git clone https://github.com/Hyena0x/evm-profiler.git
   cd evm-profiler
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run tests:
   ```bash
   go test ./...
   ```

## 📝 Coding Standards
- **Clean Architecture**: Follow the modular design in `internal/`. Keep logical separation between data fetching, analysis, and printing.
- **Error Handling**: Use professional error handling. Return errors from internal packages and handle them at the CLI level.
- **Precision**: Never use `float64` for EVM amounts. Use `math/big`.
- **Concurrency**: Use `errgroup` for managing parallel tasks.

## 🚀 Pull Request Process
1. Create a descriptive branch (e.g., `feature/add-bnb-support` or `fix/rpc-timeout`).
2. Make your changes and ensure tests pass.
3. Update `README.md` if necessary.
4. Open a Pull Request with a clear description of your changes.

## 📜 Code of Conduct
Please be respectful and helpful in all interactions. I aim to maintain a welcoming environment for everyone.

---

Thank you for making `evm-profiler` better!
