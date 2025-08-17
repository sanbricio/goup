# goup - Go Dependency Updater

A smart and user-friendly tool to update Go module dependencies with enhanced features, beautiful console output, and a well-tested, modular architecture.

## Features

- 🎯 **Smart Updates**: Updates only direct dependencies by default (with option for all)
- 🔍 **Selective Updates**: Choose exactly which dependencies to update interactively
- 🎨 **Beautiful Output**: Colored console output with progress indicators
- 🏃 **Dry Run Mode**: Preview what would be updated without making changes
- 🤝 **Interactive Mode**: Ask for confirmation before updating
- 📝 **Verbose Logging**: Detailed output when needed
- ✅ **Error Handling**: Robust error handling and reporting
- 🚀 **Fast**: Efficient dependency parsing and updating

## Installation

### Option 1: Using Make (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd goup

# Build and install using make
make install
```

### Option 2: Manual Build
```bash
# Clone the repository
git clone <repository-url>
cd goup

# Build manually
go build -o goup ./cmd/goup
mv goup ~/go/bin/
export PATH=$PATH:$HOME/go/bin
```

### Option 3: Development Setup
```bash
# Clone the repository
git clone <repository-url>
cd goup

# Install dependencies and run tests
make dev

# Install to local bin
make install
```

## Usage

### Basic Usage
```bash
# Update all direct dependencies
goup

# Show what would be updated (dry run)
goup --dry-run

# Update with confirmation prompt
goup --interactive

# Update all dependencies (including indirect)
goup --all
```

### Selective Updates
```bash
# Interactively select which dependencies to update
goup --select

# Select from all dependencies (including indirect)
goup --select --all

# Preview selective updates without making changes
goup --select --dry-run

# Combine with other options
goup --select --verbose
```

### Advanced Options
```bash
# Show detailed output during updates
goup --verbose

# Disable colored output
goup --no-color

# Combine multiple options
goup --interactive --verbose --all
```

### Command Line Options

| Flag | Description |
|------|-------------|
| `--select` | Interactively select which dependencies to update |
| `--dry-run` | Show what would be updated without making changes |
| `--interactive` | Ask for confirmation before updating |
| `--verbose` | Show detailed output during the update process |
| `--no-color` | Disable colored console output |
| `--all` | Update indirect dependencies as well as direct ones |
| `--help` | Show help message |

## Examples

### Example 1: Safe Preview
```bash
$ goup --dry-run
=== goup - Go Dependency Updater ===

[INFO] Found 3 direct dependencies to update:

  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/golang-jwt/jwt/v4@v4.5.0
  3. golang.org/x/crypto@v0.14.0

[WARNING] Dry run mode - no actual updates will be performed
```

### Example 2: Interactive Update
```bash
$ goup --interactive
=== goup - Go Dependency Updater ===

[INFO] Found 3 direct dependencies to update:

  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/golang-jwt/jwt/v4@v4.5.0
  3. golang.org/x/crypto@v0.14.0

Do you want to proceed with the update? (y/N): y

[INFO] Updating dependencies...

[PROGRESS] Updating github.com/gin-gonic/gin... (1/3)
[SUCCESS] ✓ Updated github.com/gin-gonic/gin
[PROGRESS] Updating github.com/golang-jwt/jwt/v4... (2/3)
[SUCCESS] ✓ Updated github.com/golang-jwt/jwt/v4
[PROGRESS] Updating golang.org/x/crypto... (3/3)
[SUCCESS] ✓ Updated golang.org/x/crypto

[SUCCESS] All dependencies updated successfully!

[INFO] Running go mod tidy...
[SUCCESS] ✓ go mod tidy completed

[SUCCESS] Dependency update completed!
```

### Example 3: Selective Updates
```bash
$ goup --select
=== goup - Go Dependency Updater ===

[INFO] Found 5 direct dependencies available for update:

  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/golang-jwt/jwt/v4@v4.5.0
  3. github.com/stretchr/testify@v1.8.4
  4. golang.org/x/crypto@v0.14.0
  5. golang.org/x/time@v0.3.0

[INFO] Selection options:
  • Enter numbers (e.g., 1,3,5 or 1-3 or 1,3-5)
  • Enter 'all' to select all dependencies
  • Enter package names or patterns (e.g., 'github.com/gin*')
  • Press Enter without input to cancel

Select dependencies to update: 1,3-4

[SUCCESS] Selected 3 dependencies:
  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/stretchr/testify@v1.8.4
  3. golang.org/x/crypto@v0.14.0

Proceed with these selected dependencies? (y/N): y

[INFO] Updating dependencies...
# ... continues with selected updates only ...
```

### Example 4: Pattern-based Selection
```bash
$ goup --select --all
=== goup - Go Dependency Updater ===

[INFO] Found 8 all dependencies available for update:

  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/golang-jwt/jwt/v4@v4.5.0
  3. github.com/stretchr/testify@v1.8.4
  4. golang.org/x/crypto@v0.14.0 (indirect)
  5. golang.org/x/net@v0.16.0 (indirect)
  6. golang.org/x/sys@v0.13.0 (indirect)
  7. golang.org/x/time@v0.3.0
  8. golang.org/x/text@v0.13.0 (indirect)

Select dependencies to update: golang.org/x/*

[SUCCESS] Selected 4 dependencies:
  1. golang.org/x/crypto@v0.14.0 (indirect)
  2. golang.org/x/net@v0.16.0 (indirect)
  3. golang.org/x/sys@v0.13.0 (indirect)
  4. golang.org/x/text@v0.13.0 (indirect)

Proceed with these selected dependencies? (y/N): y
# ... updates only golang.org/x/ packages ...
```

## Selection Syntax

When using `--select`, you can choose dependencies using various formats:

### By Numbers
- `1` - Select dependency #1
- `1,3,5` - Select dependencies #1, #3, and #5
- `1-3` - Select dependencies #1 through #3
- `1,3-5,7` - Select #1, #3 through #5, and #7

### By Names/Patterns
- `github.com/gin-gonic/gin` - Select specific package
- `github.com/gin*` - Select all packages starting with "github.com/gin"
- `*crypto*` - Select all packages containing "crypto"
- `golang.org/x/*` - Select all golang.org/x/ packages

### Special Keywords
- `all` - Select all available dependencies

## Project Structure

```
goup/
├── cmd/goup/              # Application entry point
│   └── main.go
├── internal/              # Private application code
│   ├── app/              # Main application logic
│   │   ├── app.go
│   │   └── app_test.go
│   ├── config/           # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── dependency/       # Dependency management
│   │   ├── interfaces.go
│   │   ├── manager.go
│   │   └── manager_test.go
│   ├── selector/         # Interactive selection
│   │   ├── interfaces.go
│   │   ├── selector.go
│   │   └── selector_test.go
│   ├── ui/              # User interface
│   │   ├── interfaces.go
│   │   ├── console.go
│   │   └── console_test.go
│   └── updater/         # Dependency updating
│       ├── interfaces.go
│       ├── updater.go
│       └── updater_test.go
├── Makefile             # Build automation
├── go.mod
├── go.sum
└── README.md
```

## Development

### Prerequisites
- Go 1.24 or later
- Make (optional, but recommended)
- golangci-lint (for linting)

### Building
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Quick build without tests
make quick
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run benchmarks
make bench
```

### Code Quality
```bash
# Format code
make fmt

# Lint code
make lint

# Vet code
make vet

# Run all quality checks
make ci
```

### Development Workflow
```bash
# Complete development workflow
make dev

# This runs: clean, deps, fmt, vet, test, build
```

## Requirements

- Go 1.24 or later
- A valid `go.mod` file in the current directory

## Error Handling

- **Missing go.mod**: Clear error message if no `go.mod` file is found
- **Parse Errors**: Helpful error messages for malformed `go.mod` files
- **Invalid Selections**: Friendly error messages for invalid selection syntax
- **Update Failures**: Continues updating other dependencies if one fails
- **Network Issues**: Proper error reporting for network-related failures

## How It Works

1. **Parse go.mod**: Reads and parses the `go.mod` file in the current directory
2. **Filter Dependencies**: Identifies direct dependencies (or all if `--all` flag is used)
3. **Selection Interface**: In selective mode, presents an interactive selection interface
4. **Display Plan**: Shows what will be updated with colored, formatted output
5. **Update**: Runs `go get -u` for each selected dependency
6. **Tidy**: Runs `go mod tidy` to clean up the module file

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Run the test suite**: `make ci`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Code Standards
- Follow Go conventions and best practices
- Write tests for new functionality
- Ensure all tests pass: `make test`
- Run linting: `make lint`
- Format code: `make fmt`
- Maintain >90% test coverage

### Commit Messages
Use conventional commit messages:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test additions/changes
- `refactor:` for code refactoring

## Future Enhancements

- [ ] Configuration file support (.gouprc, goup.yaml)
- [ ] Update scheduling and automation
- [ ] Integration with CI/CD pipelines (GitHub Actions, GitLab CI)
- [ ] Backup and rollback functionality
- [ ] Vulnerability scanning integration
- [ ] Custom update strategies (major, minor, patch)
- [ ] Progress bars for large updates
- [ ] Dependency graph visualization

## License

This project is open source. Feel free to use and modify as needed.

## Support

- 📚 **Documentation**: Check this README and code comments
- 🐛 **Bug Reports**: Open an issue on GitHub
- 💡 **Feature Requests**: Open an issue with the enhancement label
- ❓ **Questions**: Start a discussion on GitHub

---