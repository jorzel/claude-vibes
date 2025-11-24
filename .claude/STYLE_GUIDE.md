# Coding Style Guide

## Go Style

### Follow Standard Go Conventions
- Run `gofmt` on all code
- Use `golangci-lint` for linting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Naming
```go
// Exported names: PascalCase
type UserService struct {}
func NewUserService() *UserService {}

// Unexported names: camelCase
type userRepository struct {}
func calculateTotal() int {}

// Acronyms: all caps or all lowercase
type HTTPAPI struct {}  // exported
type httpClient struct {} // unexported

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30 * time.Second

// Interface names: -er suffix when possible
type Reader interface {}
type UserCreator interface {}
```

### Error Handling
```go
// Always check errors immediately
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use custom error types for domain errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Use errors.Is and errors.As for checking
if errors.Is(err, ErrNotFound) {
    // handle not found
}

var valErr *ValidationError
if errors.As(err, &valErr) {
    // handle validation error
}
```

### Function Design
```go
// Keep functions short and focused
// Accept interfaces, return structs
func ProcessOrder(repo OrderRepository, id string) (*Order, error) {
    // ...
}

// Use functional options for complex constructors
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = d
    }
}

func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        timeout: 30 * time.Second, // default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### Context Usage
```go
// Always pass context as first parameter
func FetchUser(ctx context.Context, id string) (*User, error) {
    // Use context for cancellation and deadlines
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case result := <-ch:
        return result, nil
    }
}

// Don't store context in structs (except rare cases)
```

### Concurrency
```go
// Use channels for communication
// Use sync primitives for synchronization
// Always use sync.WaitGroup or errgroup for goroutine coordination

import "golang.org/x/sync/errgroup"

func ProcessBatch(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item // capture loop variable
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }

    return g.Wait()
}
```

### Testing
```go
// Table-driven tests
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        total    float64
        expected float64
        wantErr  bool
    }{
        {
            name:     "applies 10% discount for orders over 100",
            total:    150.0,
            expected: 135.0,
            wantErr:  false,
        },
        {
            name:     "returns error for negative total",
            total:    -10.0,
            expected: 0,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateDiscount(tt.total)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}

// Use testify for assertions
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result := doSomething()
    assert.Equal(t, expected, result)
    assert.NoError(t, err)
}
```

#### Mocks
Use https://github.com/uber-go/mock gen to generate mocks for interfaces when needed.
---

## Python Style

### Follow PEP 8 and Modern Python Practices
- Use `black` for formatting
- Use `ruff` or `pylint` for linting
- Use `mypy` for type checking
- Minimum Python 3.11+
- use virtual environments (venv, poetry, etc.)

### Naming
```python
# Functions and variables: snake_case
def calculate_total_price(items: list[Item]) -> float:
    user_count = len(items)

# Classes: PascalCase
class UserRepository:
    pass

# Constants: SCREAMING_SNAKE_CASE
MAX_RETRY_ATTEMPTS = 3
DEFAULT_TIMEOUT = 30

# Private: single leading underscore
def _internal_helper():
    pass

class MyClass:
    def __init__(self):
        self._private_var = 0
```

### Type Hints (Always)
```python
from typing import Optional, Protocol
from collections.abc import Sequence

# Always use type hints for function signatures
def fetch_user(user_id: int) -> dict[str, Any]:
    pass

# Use None for optional returns
def find_user(email: str) -> Optional[User]:
    pass

# Use Protocol for structural typing
class Repository(Protocol):
    def save(self, entity: Entity) -> None: ...
    def find(self, id: str) -> Optional[Entity]: ...

# Use modern syntax (Python 3.10+)
def process_items(items: list[str] | None = None) -> dict[str, int]:
    pass
```

### Error Handling
```python
# Use custom exceptions for domain errors
class ValidationError(Exception):
    def __init__(self, field: str, message: str):
        self.field = field
        self.message = message
        super().__init__(f"Validation error on {field}: {message}")

class NotFoundError(Exception):
    pass

# Be explicit about what you catch
try:
    result = risky_operation()
except ValidationError as e:
    logger.error("Validation failed", extra={"field": e.field})
    raise
except Exception as e:
    logger.exception("Unexpected error")
    raise

# Use context managers for resources
with open("file.txt") as f:
    content = f.read()
```

### Function Design
```python
# Keep functions short and focused
# Use dataclasses for structured data
from dataclasses import dataclass

@dataclass
class CreateOrderRequest:
    user_id: str
    items: list[OrderItem]

@dataclass
class CreateOrderResponse:
    order_id: str
    status: str
    total: float

# Use dependency injection
def create_order(
    request: CreateOrderRequest,
    repository: OrderRepository,
    logger: logging.Logger,
) -> CreateOrderResponse:
    # ...
```

### Async/Await
```python
import asyncio

# Use async/await for I/O bound operations
async def fetch_user(user_id: str) -> User:
    async with httpx.AsyncClient() as client:
        response = await client.get(f"/users/{user_id}")
        return User.parse(response.json())

# Gather multiple async operations
async def fetch_all_users(user_ids: list[str]) -> list[User]:
    tasks = [fetch_user(uid) for uid in user_ids]
    return await asyncio.gather(*tasks)
```

### Testing
```python
import pytest
from unittest.mock import Mock, patch

# Use pytest fixtures
@pytest.fixture
def user_repository():
    return InMemoryUserRepository()

@pytest.fixture
def sample_user():
    return User(id="123", name="John Doe", email="john@example.com")

# Descriptive test names
def test_calculate_discount_applies_10_percent_for_orders_over_100():
    result = calculate_discount(total=150.0)
    assert result.discount == 15.0
    assert result.final_total == 135.0

def test_calculate_discount_raises_validation_error_for_negative_total():
    with pytest.raises(ValidationError) as exc_info:
        calculate_discount(total=-10.0)
    assert exc_info.value.field == "total"

# Parametrize for multiple cases
@pytest.mark.parametrize("total,expected_discount,expected_final", [
    (150.0, 15.0, 135.0),
    (100.0, 10.0, 90.0),
    (50.0, 0.0, 50.0),
])
def test_calculate_discount(total, expected_discount, expected_final):
    result = calculate_discount(total)
    assert result.discount == expected_discount
    assert result.final_total == expected_final

# Use mocks sparingly, prefer fakes
def test_create_order_saves_to_repository(user_repository):
    service = OrderService(repository=user_repository)
    order = service.create_order(user_id="123", items=[])

    saved_order = user_repository.find(order.id)
    assert saved_order is not None
    assert saved_order.user_id == "123"
```

---

## Cross-Language Principles

### Comments
```go
// Go: Use godoc style comments
// Package domain provides core business logic.
package domain

// UserService handles user-related operations.
type UserService struct {}

// CreateUser creates a new user in the system.
// It returns an error if the email is already taken.
func (s *UserService) CreateUser(ctx context.Context, email string) (*User, error) {
    // Implementation
}
```
```python
# Python: Use docstrings
"""
Module for user-related operations.
"""

class UserService:
    """Service for managing user operations."""

    def create_user(self, email: str) -> User:
        """
        Create a new user in the system.

        Args:
            email: User's email address

        Returns:
            Created user instance

        Raises:
            ValidationError: If email is invalid
            DuplicateError: If email already exists
        """
        pass
```

### What NOT to Comment
- Obvious code that explains itself
- Commented-out code (delete it)
- TODO without context or owner

### What TO Comment
- Why a decision was made
- Complex algorithms or business rules
- Non-obvious performance optimizations
- Known limitations or edge cases