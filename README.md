# Cocobase Go Client

A powerful Go client for the Cocobase Backend as a Service (BaaS).

## Features

- ✅ Full CRUD operations on documents
- ✅ Advanced query filtering with 12+ operators
- ✅ Boolean logic (AND, OR, named OR groups)
- ✅ Multi-field search
- ✅ Authentication (login, register, user management)
- ✅ Real-time updates via WebSocket
- ✅ Pluggable storage for token persistence
- ✅ Thread-safe operations
- ✅ Context support for cancellation and timeouts
- ✅ Comprehensive error handling

## Installation

```bash
go get github.com/lordace-coder/cocobase-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lordace-coder/cocobase-go/cocobase"
)

func main() {
    // Initialize client
    client := cocobase.NewClient(cocobase.Config{
        APIKey: "your-api-key",
    })

    ctx := context.Background()

    // Create a document
    doc, err := client.CreateDocument(ctx, "users", map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created: %s\n", doc.ID)
}
```

## Advanced Querying

### Basic Operators

```go
// Equality
query := cocobase.NewQuery().Where("status", "active")

// Comparison operators
query = cocobase.NewQuery().
    Filter("age", "gte", 18).
    Filter("age", "lte", 65)

// String operations
query = cocobase.NewQuery().
    Filter("email", "endswith", "gmail.com")

// IN operator
query = cocobase.NewQuery().
    Filter("role", "in", "admin,moderator,support")

// NULL checks
query = cocobase.NewQuery().
    Filter("deletedAt", "isnull", true)
```

### Boolean Logic

```go
// Simple OR
query := cocobase.NewQuery().
    Where("status", "active").
    Or("isPremium", "eq", true).
    Or("isVerified", "eq", true)

// Named OR groups
query = cocobase.NewQuery().
    OrGroup("tier", "isPremium", "eq", true).
    OrGroup("tier", "isVerified", "eq", true).
    OrGroup("location", "country", "eq", "US").
    OrGroup("location", "country", "eq", "UK")

// Multi-field search
query = cocobase.NewQuery().
    MultiFieldOr([]string{"name", "email"}, "contains", "john")
```

### Pagination & Sorting

```go
query := cocobase.NewQuery().
    Where("status", "active").
    Sort("createdAt").
    OrderDesc().
    Limit(50).
    Offset(100)

docs, err := client.ListDocuments(ctx, "users", query)
```

## Authentication

```go
// Register
err := client.Register(ctx, "user@example.com", "password", map[string]interface{}{
    "firstName": "John",
    "lastName":  "Doe",
})

// Login
err = client.Login(ctx, "user@example.com", "password")

// Get current user
user, err := client.GetCurrentUser(ctx)

// Update user
newEmail := "new@example.com"
user, err = client.UpdateUser(ctx, map[string]interface{}{
    "phone": "+1234567890",
}, &newEmail, nil)

// Logout
err = client.Logout()
```

## Real-time Updates

```go
conn, err := client.WatchCollection(ctx, "users", func(event cocobase.Event) {
    fmt.Printf("Event: %s\n", event.Event)
    fmt.Printf("Document: %+v\n", event.Data)
}, "users-watcher")

if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

## Storage Persistence

```go
import "github.com/lordace-coder/cocobase-go/storage"

// Memory storage
store := storage.NewMemoryStorage()

// File storage
store, err := storage.NewFileStorage(".cocobase/storage.json")

// Use with client
client := cocobase.NewClient(cocobase.Config{
    APIKey:  "your-api-key",
    Storage: store,
})
```

## Query Operators

| Operator     | Usage                                 |
| ------------ | ------------------------------------- |
| `eq`         | Equals (default)                      |
| `ne`         | Not equals                            |
| `gt`         | Greater than                          |
| `gte`        | Greater than or equal                 |
| `lt`         | Less than                             |
| `lte`        | Less than or equal                    |
| `contains`   | Contains substring (case-insensitive) |
| `startswith` | Starts with                           |
| `endswith`   | Ends with                             |
| `in`         | In list (comma-separated)             |
| `notin`      | Not in list                           |
| `isnull`     | Is null/not null                      |

## Examples

See the `examples/` directory for complete examples:

- `examples/basic/` - Basic CRUD operations
- `examples/advanced/` - Advanced querying
- `examples/auth/` - Authentication flows
- `examples/realtime/` - WebSocket real-time updates

## Testing

```bash
go test ./tests/...
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
