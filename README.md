# 🛒 EMMSA Price Tracker

A Go application for tracking and storing daily agricultural market prices from EMMSA. The application fetches price data and can optionally store it in Pantry (a free JSON storage service) for historical tracking and analysis.

> **Note**: This is a CLI tool designed to be run periodically (e.g., via cron) to collect price data.

## ✨ Features

- 🚀 Fetch daily agricultural prices from EMMSA's public data
- 💾 Store price data locally or in Pantry cloud storage
- 📅 Automatic date-based basket naming (`prices_YYYY_MM_DD`)
- 🔄 Simple CLI interface with flexible options
- 🧪 Comprehensive test coverage (71%+)
- 🔒 Environment-based configuration
- 📊 JSON output for easy parsing and integration

## Prerequisites

- Go 1.16 or higher
- Git
- (Optional) Pantry account for cloud storage (sign up at [getpantry.cloud](https://getpantry.cloud/))

## 🚀 Quick Start

1. **Install Go** (1.16 or higher) if you haven't already
2. **Clone and build**:
   ```bash
   git clone https://github.com/your-username/price-tracker.git
   cd price-tracker
   go build -o price-tracker ./cmd/price-tracker
   ```
3. **Run a test fetch**:
   ```bash
   ./price-tracker -output prices_today.json
   ```

## 🔧 Installation

### Prerequisites
- Go 1.16 or higher
- Git
- (Optional) [Pantry account](https://getpantry.cloud/) for cloud storage

### Build from Source
```bash
git clone https://github.com/your-username/price-tracker.git
cd price-tracker
go build -o price-tracker ./cmd/price-tracker
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# Required for Pantry integration
PANTRY_API_KEY=your_pantry_api_key  # From Pantry dashboard

# Optional: HTTP client settings
HTTP_TIMEOUT=30  # Timeout in seconds

# Optional: Log level (debug, info, warn, error)
LOG_LEVEL=info
```

### Configuration Notes
- The application will work without a `.env` file if using local file output only
- All environment variables have sensible defaults

## 💻 Usage

### Basic Commands

```bash
# Fetch today's prices and print to console
./price-tracker

# Save to a JSON file
./price-tracker -output prices.json

# Fetch for a specific date (YYYY-MM-DD)
./price-tracker -date 2025-06-18

# Store in Pantry
./price-tracker -pantry

# Combine options
./price-tracker -pantry -date 2025-06-18 -output prices_20250618.json
```

### Command Line Options

```
  -date string
        Date in YYYY-MM-DD format (default: today)
  -output string
        Output file path (default: stdout)
  -pantry
        Store data in Pantry
  -v    Show version
```

### Output Format

Example JSON output:

```json
[
  {
    "product": "PAPA BLANCA",
    "variety": "COMUN",
    "market": "Mercado Mayorista de Santa Anita",
    "min_price": 2.5,
    "max_price": 3.0,
    "unit": "kg",
    "date": "2025-06-18T00:00:00Z"
  },
  ...
]
```

## 💾 Data Storage

### Local Storage

By default, the application outputs price data to stdout or a specified file in JSON format. The data includes timestamps and is structured for easy parsing.

### ☁️ Pantry Integration

[Pantry](https://getpantry.cloud/) is a free JSON storage service. Each day's prices are stored in a separate basket named `prices_YYYY_MM_DD`.

#### Managing Pantry Data

Use the included `pantry-cli` tool:

```bash
# Build the CLI
go build -o pantry-cli ./cmd/pantry-cli

# List all baskets
./pantry-cli list

# Get basket contents
./pantry-cli get prices_2025_06_18

# Delete a basket
./pantry-cli delete prices_2025_06_18
```

#### Pantry CLI Commands

```
Available commands:
  list      List all baskets
  get       Get basket contents
  delete    Delete a basket
  help      Show help
```

## 🛠 Development

### Building

```bash
# Build main application
go build -o price-tracker ./cmd/price-tracker

# Build pantry CLI
go build -o pantry-cli ./cmd/pantry-cli
```

## 🧪 Testing & Code Quality

### Running Tests

Run the test suite:

```bash
go test -v ./...
```

### Linting

This project uses `golangci-lint` for code quality checks. The linter configuration is in `.golangci.yml`.

#### Running Linters Locally

1. Install `golangci-lint` (if not already installed):
   ```bash
   # Using Homebrew (macOS/Linux)
   brew install golangci-lint
   
   # Or using go install
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. Run the linter:
   ```bash
   # Run all linters
   make lint
   
   # Or run directly
   golangci-lint run
   ```

#### GitHub Integration

Linting is automatically performed on every push and pull request via GitHub Actions. The workflow file is located at `.github/workflows/lint.yml`.

A status badge is available to show the current linting status:

```markdown
[![Lint Status](https://github.com/your-username/price-tracker/actions/workflows/lint.yml/badge.svg)](https://github.com/your-username/price-tracker/actions/workflows/lint.yml)
```

#### Common Linting Issues

- Formatting: Run `gofmt -w .` to automatically format your code
- Imports: Run `goimports -w .` to organize imports
- To see all issues: `golangci-lint run -v`
- To fix auto-fixable issues: `golangci-lint run --fix`

## 🛠️ Development

### Code Quality

This project uses [golangci-lint](https://golangci-lint.run/) for code quality checks. To get started:

```bash
# Install golangci-lint (requires Go 1.16+)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0

# Run linters
golangci-lint run

# Or use the Makefile
make lint
```

### Testing

```bash
# Run all tests
make test

# Run tests with race detector
make test-race

# Generate coverage report
make cover

# Run tests for a specific package
go test -v -cover ./internal/storage/pantry/
```

### Common Tasks

```bash
# Build the application
make build

# Clean build artifacts
make clean

# Install development tools
make tools
```

### Code Structure

```
.
├── cmd/
│   ├── price-tracker/  # Main application
│   └── pantry-cli/     # Pantry management tool
├── internal/
│   ├── api/emmsa/       # EMMSA API client
│   └── storage/pantry/  # Pantry storage client
└── scripts/            # Build and deployment scripts
```

### Testing Notes
- The test suite includes unit tests and integration tests
- Mock servers are used for testing external API calls
- Current test coverage: 71%+
- Run `make cover` to generate a coverage report
- Run `make lint` to check code quality

## 🚨 Troubleshooting

### Common Issues

**EMMSA API Connection Issues**
- Ensure you have an active internet connection
- The EMMSA API might be temporarily unavailable
- Check if the API endpoint has changed

**Pantry Integration**
- Verify your `PANTRY_API_KEY` is set correctly
- Check your pantry dashboard for usage limits
- Ensure the basket naming follows `prices_YYYY_MM_DD` format

**Build Issues**
- Make sure you're using Go 1.16 or higher
- Run `go mod tidy` to ensure all dependencies are properly downloaded

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  Made with ❤️ for agricultural market analysis
</div>
