# Testing Specification

## Testing Philosophy
- **Test behavior, not implementation**
- **Fast feedback**: Unit tests <5s, integration <30s
- **Clear failures**: Test names explain what broke
- **Realistic data**: No "foo", "bar", "test" values

## Go Testing Standards

### Test Organization
```go
// _test.go files in same package
package service

import "testing"

// Table-driven tests for multiple scenarios
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
        errType error
    }{
        {
            name:    "creates user with valid email",
            email:   "user@example.com",
            wantErr: false,
        },
        {
            name:    "returns error for invalid email",
            email:   "invalid",
            wantErr: true,
            errType: ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := NewUserService(newMockRepo())
            user, err := svc.CreateUser(context.Background(), tt.email)

            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if tt.wantErr && !errors.Is(err, tt.errType) {
                t.Errorf("expected error type %v, got %v", tt.errType, err)
            }

            if !tt.wantErr && user.Email != tt.email {
                t.Errorf("expected email %v, got %v", tt.email, user.Email)
            }
        })
    }
}
```

### Integration Tests
```go
// Place in tests/ directory or use build tags
// +build integration

package tests

func TestCreateOrder_Integration(t *testing.T) {
    // Setup real database
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewPostgresOrderRepository(db)
    svc := service.NewOrderService(repo)

    // Test actual behavior
    order, err := svc.CreateOrder(context.Background(), &service.CreateOrderRequest{
        UserID: "user_123",
        Items:  []service.OrderItem{{ProductID: "prod_1", Quantity: 2}},
    })

    require.NoError(t, err)
    assert.NotEmpty(t, order.ID)

    // Verify in database
    found, err := repo.FindByID(context.Background(), order.ID)
    require.NoError(t, err)
    assert.Equal(t, order.ID, found.ID)
}
```

### Test Helpers
```go
// Use test helpers to reduce boilerplate
func newTestUserService(t *testing.T) *UserService {
    t.Helper()
    return NewUserService(newMockRepo())
}

// Use testify require for setup that must succeed
import "github.com/stretchr/testify/require"

func TestSomething(t *testing.T) {
    db := setupDB(t)
    require.NotNil(t, db, "database setup must succeed")

    // Continue with test
}
```

---

## Python Testing Standards

### Test Organization
```python
# tests/ directory mirrors src/ structure
# tests/test_user_service.py

import pytest
from myapp.service import UserService
from myapp.repository import InMemoryUserRepository

class TestUserService:
    """Tests for UserService."""

    @pytest.fixture
    def repository(self):
        return InMemoryUserRepository()

    @pytest.fixture
    def service(self, repository):
        return UserService(repository=repository)

    def test_create_user_with_valid_email(self, service):
        user = service.create_user(email="user@example.com")

        assert user.email == "user@example.com"
        assert user.id is not None

    def test_create_user_raises_validation_error_for_invalid_email(self, service):
        with pytest.raises(ValidationError) as exc_info:
            service.create_user(email="invalid")

        assert exc_info.value.field == "email"

    @pytest.mark.parametrize("email,should_succeed", [
        ("valid@example.com", True),
        ("another@test.com", True),
        ("invalid", False),
        ("", False),
    ])
    def test_create_user_with_various_emails(self, service, email, should_succeed):
        if should_succeed:
            user = service.create_user(email=email)
            assert user.email == email
        else:
            with pytest.raises(ValidationError):
                service.create_user(email=email)
```

### Async Tests
```python
import pytest

@pytest.mark.asyncio
async def test_fetch_user_returns_user_data():
    service = AsyncUserService()
    user = await service.fetch_user(user_id="123")

    assert user.id == "123"
    assert user.name is not None
```

### Integration Tests
```python
import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

@pytest.fixture(scope="session")
def db_engine():
    """Create test database engine."""
    engine = create_engine("postgresql://test:test@localhost/testdb")
    Base.metadata.create_all(engine)
    yield engine
    Base.metadata.drop_all(engine)

@pytest.fixture
def db_session(db_engine):
    """Create a new session for each test."""
    connection = db_engine.connect()
    transaction = connection.begin()
    session = Session(bind=connection)

    yield session

    session.close()
    transaction.rollback()
    connection.close()

@pytest.mark.integration
def test_create_order_persists_to_database(db_session):
    repository = PostgresOrderRepository(session=db_session)
    service = OrderService(repository=repository)

    order = service.create_order(
        user_id="user_123",
        items=[OrderItem(product_id="prod_1", quantity=2)]
    )

    assert order.id is not None

    # Verify in database
    found = repository.find_by_id(order.id)
    assert found is not None
    assert found.user_id == "user_123"
```

---

## Test Coverage Requirements
- **Critical paths**: 100% (auth, payments, data integrity)
- **Business logic**: 90%+
- **API handlers/routes**: 80%+
- **Utilities**: 80%+

## Mocking Strategy
- **Go**: Use interfaces, implement test doubles
- **Python**: Prefer fakes, use mocks sparingly
- **Both**: Only mock external boundaries (APIs, databases for unit tests)

## Test Data Factories

### Go
```go
func NewTestUser(overrides ...func(*User)) *User {
    u := &User{
        ID:    uuid.New().String(),
        Email: "test@example.com",
        Name:  "Test User",
    }
    for _, override := range overrides {
        override(u)
    }
    return u
}

// Usage
user := NewTestUser(func(u *User) {
    u.Email = "custom@example.com"
})
```

### Python
```python
from factory import Factory, Faker

class UserFactory(Factory):
    class Meta:
        model = User

    id = Faker("uuid4")
    email = Faker("email")
    name = Faker("name")

# Usage
user = UserFactory.create(email="custom@example.com")
```

---

## Running Tests

### Makefile
- create a Makefile with common test commands
    - `make test`
    - `make unit-test`
    - `make integration-test`

### Go
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...

# Run integration tests only
go test -tags=integration ./tests/...

# Run specific test
go test -run TestUserService_CreateUser ./internal/service
```

### Python
```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=src --cov-report=html

# Run integration tests only
pytest -m integration

# Run specific test
pytest tests/test_user_service.py::TestUserService::test_create_user

# Run with output
pytest -v -s
```