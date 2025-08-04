# Test Strategy and Standards

## Testing Philosophy
- **Approach:** Test-Driven Development (TDD) with comprehensive test automation
- **Coverage Goals:** Minimum 80% code coverage for unit tests, 70% for integration tests
- **Test Pyramid:** 70% unit tests, 20% integration tests, 10% end-to-end tests

## Test Types and Organization

### Unit Tests
- **Framework:** Go native testing framework with Testify v1.8+
- **File Convention:** `*_test.go` files alongside source code
- **Mocking Library:** Testify/mock and GoMock for interface mocking
- **Coverage Requirement:** Minimum 80% line coverage

### Integration Tests
- **Scope:** Database interactions, NATS messaging, HTTP API endpoints
- **Location:** `tests/integration/` directory
- **Test Infrastructure:** Testcontainers PostgreSQL for realistic database testing

### End-to-End Tests
- **Framework:** Playwright v1.40+ for web UI testing
- **Scope:** Critical user journeys across frontend and backend
- **Environment:** Dedicated test environment with clean data
