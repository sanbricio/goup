# goup - Go Dependency Updater

A smart and user-friendly tool to update Go module dependencies with enhanced features and beautiful console output.

## Features

- 🎯 **Smart Updates**: Updates only direct dependencies by default (with option for all)
- 🎨 **Beautiful Output**: Colored console output with progress indicators
- 🔍 **Dry Run Mode**: Preview what would be updated without making changes
- 🤝 **Interactive Mode**: Ask for confirmation before updating
- 📝 **Verbose Logging**: Detailed output when needed
- ✅ **Error Handling**: Robust error handling and reporting
- 🚀 **Fast**: Efficient dependency parsing and updating

## Installation

### Option 1: Build and Install Locally
```bash
# Clone or download the project
go build -o goup main.go
mv goup ~/go/bin/
export PATH=$PATH:$HOME/go/bin
```

### Option 2: Direct Installation
```bash
# Make sure ~/go/bin is in your PATH
export PATH=$PATH:$HOME/go/bin

# Build and install in one step
go build -o ~/go/bin/goup main.go
```

### Option 3: Global Installation
```bash
# Install to system-wide location (may require sudo)
go build -o goup main.go
sudo mv goup /usr/local/bin/
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

### Example 3: Update All Dependencies
```bash
$ goup --all --verbose
=== goup - Go Dependency Updater ===

[INFO] Found 5 all dependencies to update:

  1. github.com/gin-gonic/gin@v1.9.1
  2. github.com/golang-jwt/jwt/v4@v4.5.0
  3. golang.org/x/crypto@v0.14.0 (indirect)
  4. golang.org/x/net@v0.16.0 (indirect)
  5. golang.org/x/sys@v0.13.0 (indirect)

[INFO] Updating dependencies...
# ... detailed verbose output ...
```

## Requirements

- Go 1.19 or later
- A valid `go.mod` file in the current directory

## How It Works

1. **Parse go.mod**: Reads and parses the `go.mod` file in the current directory
2. **Filter Dependencies**: Identifies direct dependencies (or all if `--all` flag is used)
3. **Display Plan**: Shows what will be updated with colored, formatted output
4. **Update**: Runs `go get -u` for each dependency
5. **Tidy**: Runs `go mod tidy` to clean up the module file

## Error Handling

- **Missing go.mod**: Clear error message if no `go.mod` file is found
- **Parse Errors**: Helpful error messages for malformed `go.mod` files
- **Update Failures**: Continues updating other dependencies if one fails
- **Network Issues**: Proper error reporting for network-related failures

## Contributing

Feel free to contribute by:
- Reporting bugs
- Suggesting new features
- Submitting pull requests
- Improving documentation

## Future Enhancements

- [ ] Configuration file support
- [ ] Selective dependency updates
- [ ] Update scheduling and automation
- [ ] Integration with CI/CD pipelines
- [ ] Backup and rollback functionality
- [ ] Vulnerability scanning integration
- [ ] Custom update strategies
- [ ] Progress bars for large updates

## License

This project is open source. Feel free to use and modify as needed.