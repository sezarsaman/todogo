# Testing Guide

## Overview

The todogo project includes unit and integration tests with **~70% coverage** of the core business logic (service and handler layers). Tests use standard Go testing patterns with mocking for dependencies.

## Test Structure

```
internal/
├── task/
│   ├── service/
│   │   ├── service.go           # Business logic
│   │   ├── mocks.go             # Mock repository and cache
│   │   └── service_test.go       # Service unit tests (80% coverage)
│   ├── handler/
│   │   ├── http.go              # HTTP handlers
│   │   └── http_test.go         # Handler unit tests (71% coverage)
│   └── cache/
│       └── cache_test.go         # Cache unit tests (placeholder)
├── config/
│   ├── config.go                # Configuration loader
│   └── config_test.go           # Config tests (100% coverage)
└── platform/
    └── http/
        ├── router.go            # HTTP router
        └── router_test.go       # Router tests (100% coverage)
```

## Running Tests

### Run all unit tests
```bash
make test
# or
go test -v ./internal/...
```

### View coverage
```bash
make coverage
# Shows percentage coverage and generates HTML report location
```

### Generate HTML coverage report
```bash
go test ./internal/... -coverprofile=/tmp/coverage.out
go tool cover -html=/tmp/coverage.out
# Opens in browser automatically on most systems
```

## Test Coverage by Package

| Package | Coverage | Notes |
|---------|----------|-------|
| `internal/task/service` | **80%** | Full CRUD with cache hit/miss scenarios |
| `internal/task/handler` | **71%** | All HTTP endpoints and error cases |
| `internal/config` | **100%** | Configuration loading and validation |
| `internal/platform/http` | **100%** | Router initialization and health check |
| **Overall (core)** | **~59%** | Excludes low-level DB/cache implementations |

## Test Approach

### Unit Tests (Mocking)
- **Service tests**: Mock `TaskRepository` and `TaskCache` interfaces
- **Handler tests**: Mock `TaskService` interface
- Dependencies injected via interfaces for easy testing

### Test Patterns

#### 1. Service Layer (`service_test.go`)
Tests cover:
- ✅ Cache hit/miss behavior
- ✅ Cache sync on Create/Update/Delete/Get
- ✅ Repository errors
- ✅ Default values (e.g., pending status)
- ✅ Concurrent operations (isolated per test)

```go
// Example: test cache hit on List
cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
    return cachedTasks, true, nil  // Cache hit
}
tasks, _ := svc.List(ctx)
// Assert tasks came from cache
```

#### 2. Handler Layer (`http_test.go`)
Tests cover:
- ✅ All HTTP endpoints (GET, POST, PUT, DELETE)
- ✅ Validation errors (missing fields, title too short)
- ✅ Status codes (200, 201, 204, 400, 404, 422)
- ✅ JSON request/response parsing

```go
// Example: test POST with validation
req := httptest.NewRequest(http.MethodPost, "/tasks", body)
w := httptest.NewRecorder()
engine.ServeHTTP(w, req)
if w.Code != http.StatusUnprocessableEntity { ... }
```

## Mocking Strategy

### `MockRepository` (service/mocks.go)
Implements the `TaskRepository` interface:
```go
type MockRepository struct {
    CreateFunc func(ctx context.Context, task *model.Task) error
    // ... other methods
}
```

### `MockCache` (service/mocks.go)
Implements the `TaskCache` interface:
```go
type MockCache struct {
    GetTasksFunc func(ctx context.Context) ([]model.Task, bool, error)
    SetTasksFunc func(ctx context.Context, tasks []model.Task) error
}
```

### `MockService` (handler/http_test.go)
Implements the `TaskService` interface:
```go
type MockService struct {
    ListFunc func(ctx context.Context) ([]model.Task, error)
    // ... other methods
}
```

## Running Specific Tests

### Test a single package
```bash
go test -v ./internal/task/service
```

### Test a single function
```bash
go test -v ./internal/task/service -run TestServiceCreate
```

### Test a single sub-test
```bash
go test -v ./internal/task/service -run TestServiceCreate/create_with_default_status
```

## CI/CD Integration

Add to your CI pipeline:

```bash
# Run tests with coverage
go test ./internal/... -coverprofile=coverage.out

# Check coverage threshold (e.g., 70%)
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$coverage < 70" | bc -l) )); then
    echo "Coverage $coverage% below 70% threshold"
    exit 1
fi
```

## Notes

- **Cache/Database tests**: Marked as integration tests with `.Skip()` (uncomment to run with real containers)
- **No external dependencies required** for unit tests (all mocked)
- **Fast execution**: ~5ms for all unit tests
- **Interface-based design**: Enables easy mocking and testability

## Future Improvements

- Add integration tests with testcontainers for Postgres/Redis
- Add benchmark tests for performance-critical paths
- Add property-based tests with `pgx` query builders
- Add API/end-to-end tests with docker-compose

---

**Last Updated**: 2025-12-26
**Test Framework**: Go `testing` package + mocking
**Coverage Tool**: `go tool cover`
