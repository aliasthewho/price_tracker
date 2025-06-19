# Contributing to Price Tracker

Thank you for your interest in contributing! We welcome all contributions, from bug reports to new features and documentation improvements.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

## Code of Conduct

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
   ```bash
   git clone git@github.com:your-username/price-tracker.git
   cd price-tracker
   ```
3. Set up your development environment (see [Development](#development))
4. Create a feature branch
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development

### Prerequisites

- Go 1.16 or higher
- Git
- (Optional) Docker for containerized development

### Building

```bash
# Build the application
go build -o price-tracker ./cmd/price-tracker

# Build the pantry CLI
go build -o pantry-cli ./cmd/pantry-cli
```

### Running Tests

```bash
# Run all tests
go test -v -cover ./...

# Run tests with race detector
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Run `gofmt` and `goimports` before committing
- Keep lines under 100 characters
- Document all exported functions, types, and packages

## Documentation

### Go Documentation

- Document all exported functions, types, and packages using GoDoc comments
- Include examples in the documentation when appropriate
- Keep documentation up-to-date with code changes

### Updating README.md

When making significant changes to the project:
1. Update the relevant sections in README.md
2. Ensure all commands and examples are accurate
3. Update the table of contents if needed

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build
2. Update the README.md with details of changes if needed
3. Ensure tests pass and add new tests as appropriate
4. Update documentation as needed
5. Make sure your code lints (run `golangci-lint run`)
6. Submit a pull request with a clear description of the changes

## Reporting Bugs

Please use [GitHub Issues](https://github.com/your-username/price-tracker/issues) to report any bugs. Include:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Any relevant logs or screenshots

## Feature Requests

We welcome feature requests! Please open an issue with:

- A clear description of the feature
- The problem it solves
- Any alternative solutions considered

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
