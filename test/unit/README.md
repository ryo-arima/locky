# Unit Tests

This directory contains unit tests for the locky project with 100% code coverage goal.

## Prerequisites

- Go 1.22 or later
- No external dependencies (database, Redis, AWS) required for unit tests
- **All test data must be placed in `test/unit/testdata/`**
- Unit tests use mocks and test data in `testdata/`
- Integration tests requiring external services are in `test/e2e/`

## Principles

1. **Test Data Centralization**: All test data (JSON, YAML, text files) must be in `test/unit/testdata/`
2. **No External Dependencies**: Unit tests do not connect to real databases, Redis, or AWS
3. **Use Mocks**: Use mock implementations from `test/unit/mock/` for external dependencies
4. **Use Test Utilities**: Use helper functions from `test/unit/internal/testutil/` for loading test data

## Directory Structure

```
test/unit/
├── cmd/                    # Tests for cmd/ packages
│   └── client/
│       └── main.go
├── internal/              # Internal test utilities
│   └── testutil/         # Test data loading helpers
│       ├── helper.go
│       └── helper_test.go
├── mock/                  # Mock implementations
│   ├── server/
│   │   ├── controller/   # Mock controllers
│   │   ├── usecase/      # Mock usecases
│   │   └── repository/   # Mock repositories
│   └── client/           # Mock clients
├── pkg/                   # Tests for pkg/ packages
│   ├── client/
│   ├── code/
│   ├── config/
│   ├── entity/
│   ├── logger/
│   └── server/
└── testdata/             # Test data files
    ├── config/
    ├── entity/
    ├── request/
    ├── response/
    └── casbin/
```

## Test Coverage Status

| Package | Coverage | Status |
|---------|----------|--------|
| pkg/code | 100% | ✅ Complete |
| pkg/logger | 100% | ✅ Complete |
| pkg/config | 0% | 🚧 In Progress |
| pkg/entity | 0% | 📋 Planned |
| pkg/server/controller | 0% | 📋 Planned |
| pkg/server/usecase | 0% | 📋 Planned |
| pkg/server/repository | 0% | 📋 Planned |
| pkg/server/middleware | 0% | 📋 Planned |
| pkg/client | 0% | 📋 Planned |

## Running Tests

### Run all unit tests
```bash
go test ./test/unit/pkg/...
```

### Run tests with coverage
```bash
# Check coverage for specific package
go test -cover -coverpkg=./pkg/code/... ./test/unit/pkg/code/...
go test -cover -coverpkg=./pkg/logger/... ./test/unit/pkg/logger/...
go test -cover -coverpkg=./pkg/config/... ./test/unit/pkg/config/...

# Run coverage check script for all packages
./test/unit/coverage.sh
```

### Generate coverage report
```bash
# Generate coverage data and HTML report in one command
go test -coverpkg=./pkg/... ./test/unit/pkg/... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
# or
xdg-open coverage.html  # Linux
```

### Run specific package tests
```bash
go test -v ./test/unit/pkg/code/...
go test -v ./test/unit/pkg/logger/...
go test -v ./test/unit/pkg/config/...
```

## Using Test Data

### Loading JSON Data

```go
import "github.com/ryo-arima/locky/test/unit/internal/testutil"

// Load JSON into a struct
var user model.User
err := testutil.LoadJSONFile("entity/user.json", &user)
```

### Loading YAML Data

```go
// Load YAML as bytes
yamlData, err := testutil.LoadYAMLFile("config/app.yaml")
```

### Loading Text Files

```go
// Load any file as text
content, err := testutil.LoadTextFile("config/app_invalid.yaml")
```

## Writing New Tests

### Test File Naming
- Test files should end with `_test.go`
- Place tests in the same package structure as the code being tested
- Example: `pkg/code/mcode.go` → `test/unit/pkg/code/mcode_test.go`

### Using Mocks
Mocks are generated using `go.uber.org/mock` and placed in `test/unit/mock/` directory:

```go
import "github.com/ryo-arima/locky/test/unit/mock/server/repository"

ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockRepo := repository.NewMockUserRepository(ctrl)
mockRepo.EXPECT().GetUser(gomock.Any()).Return(expectedUser, nil)
```

### Test Data Guidelines

1. **Organize by Type**: Place test data in appropriate testdata subdirectories
2. **Naming Convention**: Use descriptive names (e.g., `user_invalid.json`, `config_minimal.yaml`)
3. **Format**: Use JSON for entities, YAML for configs, CSV for Casbin
4. **Documentation**: Add comments explaining test scenarios
5. **Maintenance**: Keep minimal and focused on specific scenarios

## Test Coverage Goals

- **Target**: 100% code coverage for all packages
- **Strategy**: 
  - Unit tests for individual functions
  - Integration tests in `test/e2e/`
  - Mock external dependencies
  - Test edge cases and error paths

## Contributing

When adding new features:
1. Write tests first (TDD approach)
2. Ensure coverage doesn't decrease
3. Add mock implementations if needed
4. Update this README with coverage status
5. Add test data to `testdata/` if required

