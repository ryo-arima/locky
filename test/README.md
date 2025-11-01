# Test Runner Usage

This project provides a rich TUI-based test runner using Bubble Tea.

## Directory Structure

The test directory has a simple and intuitive structure:

```
test/
├── client/              # Client-related tests (to be added in future)
├── cmd/                 # Command-line tools such as test runners
│   └── tui_runner.go    # TUI test runner
├── config/              # Configuration-related tests (to be added in future)
├── controller/          # Controller layer tests
│   ├── test_helper.go   # Test helper
│   └── *.go             # Controller test files
├── entity/              # Entity-related tests (to be added in future)
├── repository/          # Repository layer tests
│   ├── test_helper.go   # Test helper
│   └── *.go             # Repository test files
├── results/             # Test results and coverage reports
└── run_tests.sh         # Test execution script
```

## Usage

### Via Shell Script (Recommended)

```bash
# Test entire server
./test/run_tests.sh server

# Test repository layer
./test/run_tests.sh repository

# Test controller layer
./test/run_tests.sh controller

# Test with coverage
./test/run_tests.sh coverage
```

### Via Makefile

```bash
# Test entire server
make test-server

# Test repository layer
make test-repository

# Test controller layer
make test-controller

# Test with coverage
make test-coverage
```

### Direct Execution

```bash
# Run TUI runner directly
go run ./test/cmd/tui_runner.go server
go run ./test/cmd/tui_runner.go repository
go run ./test/cmd/tui_runner.go controller
go run ./test/cmd/tui_runner.go coverage
```

### Standard Go Test Commands

```bash
# Test specific packages
go test ./test/repository/... -v
go test ./test/controller/... -v

# Run all tests
go test ./test/... -v

# Test with coverage
go test ./test/... -cover
```

## Features

- **Rich TUI**: Beautiful terminal interface using Bubble Tea
- **Automatic Mode Switching**: Auto-detection of interactive/non-interactive modes
- **Real-time Display**: Visual confirmation of test execution in real-time
- **Colorful Output**: Color-coded display of success/failure/errors
- **Coverage Reports**: Generate HTML and text format coverage reports
- **Simple Operation**: Exit with 'q' key (interactive mode)

## Mode Switching

The test runner automatically detects the execution environment and switches modes:

### Interactive Mode (TUI Mode)
- When executed directly in terminal
- Beautiful Bubble Tea-based UI
- Real-time display of test progress
- Program control via keyboard operations

### Non-interactive Mode (Standard Output Mode)
Automatically switches to non-interactive mode in the following cases:
- CI environments (when environment variables such as `CI`, `GITHUB_ACTIONS`, `GITLAB_CI` are set)
- Output redirection (`./test/run_tests.sh > output.txt`)
- Via pipe (`./test/run_tests.sh | grep PASS`)
- When `NON_INTERACTIVE=1` environment variable is set

```bash
# Explicitly run in non-interactive mode
NON_INTERACTIVE=1 ./test/run_tests.sh server

# Save output to file (automatically non-interactive)
./test/run_tests.sh coverage > test_results.txt

# Pass to other commands via pipe (automatically non-interactive)
./test/run_tests.sh repository | grep -E "(PASS|FAIL)"
```

## TUI Operations

- **q**: Exit program
- **Ctrl+C**: Force exit

## Generated Files

The following files are generated when running coverage tests:

- `test/results/coverage.html`: HTML coverage report
- `test/results/coverage.txt`: Text coverage report
- `test/results/coverage.out`: Raw coverage data

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Terminal styling
- [go-cmp](https://github.com/google/go-cmp): Test assertions
- [sqlmock](https://github.com/DATA-DOG/go-sqlmock): Database mocking
