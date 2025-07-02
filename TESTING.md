# Run all tests

make test

# Run with verbose output

make test-verbose

# Generate coverage report

make test-coverage

# Test specific packages

make test-utils
make test-middleware
make test-usecase
make test-handler

# Full CI pipeline

make ci
