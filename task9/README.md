# Task Management API - Testing Suite

This project includes a comprehensive testing suite using Testify framework following Clean Architecture principles.

## Test Structure

```
tests/
├── mocks/                          # Testify mocks for repositories
│   ├── mock_task_repository.go
│   └── mock_user_repository.go
├── infrastructure/                 # Infrastructure layer tests
│   ├── password_service_test.go
│   └── jwt_service_test.go
├── usecases/                       # Use case layer tests
│   ├── task_usecases_test.go
│   └── user_usecases_test.go
├── middleware/                     # Middleware tests
│   └── auth_middleware_test.go
├── controllers/                    # Controller tests
│   └── controller_test.go
├── routers/                        # Router tests
│   └── router_test.go
└── repositories_integration/       # Integration tests
    ├── task_repository_integration_test.go
    └── user_repository_integration_test.go
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Run Tests with Detailed Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Specific Test Package

```bash
# Unit tests only
go test ./tests/infrastructure/...
go test ./tests/usecases/...

# Integration tests (requires MongoDB)
go test ./tests/repositories_integration/...

# Middleware tests
go test ./tests/middleware/...

# Controller tests
go test ./tests/controllers/...

# Router tests
go test ./tests/routers/...
```

### Run Specific Test Function

```bash
go test -v ./tests/infrastructure -run TestBcryptHasher_Hash
```

## Test Coverage

The test suite aims for 70%+ coverage across use cases and services. To view coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

To view HTML coverage report:

```bash
go tool cover -html=coverage.out
```

## Test Categories

### Unit Tests

- **Infrastructure Tests**: Test password hashing and JWT token generation/validation
- **Use Case Tests**: Test business logic with mocked repositories
- **Middleware Tests**: Test authentication and authorization middleware

### Integration Tests

- **Repository Tests**: Test MongoDB operations with a test database
- Requires MongoDB to be running (uses `task_manager_test` database)

### Controller & Router Tests

- **Controller Tests**: Test HTTP handlers with mocked use cases
- **Router Tests**: Test complete request/response flow including authentication

## Prerequisites for Integration Tests

Integration tests require MongoDB to be running:

```bash
# Set MongoDB URI (optional, defaults to localhost:27017)
export MONGODB_URI="mongodb://localhost:27017"
```

Integration tests use a separate test database (`task_manager_test`) and clean up after themselves.

## Mock Usage

The test suite uses Testify mocks for:

- `TaskRepository` - Mocked in use case tests
- `UserRepository` - Mocked in use case tests

Mocks allow testing business logic in isolation without database dependencies.

## Example Test Output

```
=== RUN   TestBcryptHasher_Hash
--- PASS: TestBcryptHasher_Hash (0.00s)
=== RUN   TestBcryptHasher_Compare
--- PASS: TestBcryptHasher_Compare (0.00s)
=== RUN   TestAuthUseCase_Register
--- PASS: TestAuthUseCase_Register (0.00s)
...
PASS
coverage: 75.3% of statements
```

## Test Best Practices

1. **Isolation**: Each test is independent and doesn't rely on other tests
2. **Cleanup**: Integration tests clean up test data after execution
3. **Mocking**: External dependencies are mocked in unit tests
4. **Assertions**: Use Testify assertions for clear test failures
5. **Coverage**: Aim for high coverage of business logic and critical paths

## Troubleshooting

### Integration Tests Failing

- Ensure MongoDB is running
- Check `MONGODB_URI` environment variable
- Verify network connectivity to MongoDB

### Mock Assertions Failing

- Ensure all expected mock calls are made
- Check mock setup matches actual method calls
- Verify mock expectations are set before test execution
