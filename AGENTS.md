# AGENTS.md - Go SaltStack Client Library

## Build, Test, and Lint Commands

### Standard Go Commands

```bash
# Build the library
go build ./...

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test function
go test -v -run TestFunctionName ./path/to/package

# Format code
go fmt ./...

# Vet the code
go vet ./...

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build
```

### Note

This project currently has no test files. When adding tests, follow these patterns:

- Use table-driven tests for multiple test cases
- Test both success and error paths
- Use `t.Run()` for subtests when appropriate

## Code Style Guidelines

### Imports

- Group imports: standard library first, third-party packages second
- Separate groups with blank lines
- No blank imports (unused imports)
- Example:

  ```go
  import (
      "context"
      "encoding/json"

      "github.com/tidwall/gjson"
  )
  ```

### Naming Conventions

- **Exported types**: PascalCase (e.g., `Client`, `Minion`, `Job`)
- **Exported functions**: PascalCase (e.g., `NewClient`, `Login`, `ListMinions`)
- **Exported constants**: PascalCase (e.g., `LocalClient`, `Glob`)
- **Private fields/variables**: camelCase (e.g., `token`, `skipVerify`)
- **Interface types**: PascalCase (e.g., `ClientOption`, `RunOption`)

### Functional Options Pattern (Heavily Used)

This library uses functional options extensively for configuration:

- Option function types: `ClientOption`, `RunOption`
- Option functions named with `With` prefix: `WithEndpoint`, `WithInsecure`, `WithGlobTarget`
- Options accept a pointer to the config struct and modify it
- Example:

  ```go
  type ClientOption func(*Client)

  func WithTimeout(timeout time.Duration) ClientOption {
      return func(opt *Client) {
          opt.timeout = timeout
      }
  }
  ```

### Error Handling

- Always return errors, never panic (except in `NewClient` init functions)
- Use `fmt.Errorf` with `%w` for error wrapping when appropriate
- Provide context in error messages
- Early return pattern: check error immediately, return if present
- Example:

  ```go
  if err != nil {
      return nil, fmt.Errorf("failed to %s: %w", action, err)
  }
  ```

### JSON/HTTP Response Handling

- Use `encoding/json` for JSON marshaling/unmarshaling
- Use `github.com/tidwall/gjson` for parsing JSON when you need to access specific fields
- Define response structs for API responses with JSON tags
- Check for response existence before parsing (e.g., `gjson.Exists()`)
- Example:

  ```go
  if v := gjson.GetBytes(data, "return.0.data.success"); !v.Bool() {
      return nil, fmt.Errorf("failed to %s", action)
  }
  ```

### Context Usage

- All client methods accept `context.Context` as first parameter
- Pass context to HTTP requests using `http.NewRequestWithContext`
- Use `context.TODO()` when no specific context is available

### HTTP Client Pattern

- Client struct contains `*http.Client`, `endpoint`, `token`
- Generic `doRequest`, `get`, `post` methods for common operations
- Set standard headers: `Accept`, `Content-Type`, `X-Requested-With`, `X-Auth-Token`
- Handle token authentication in header
- Use cookie jar for session management
- Example:

  ```go
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")
  if c.token != "" {
      req.Header.Set("X-Auth-Token", c.token)
  }
  ```

### Documentation

- Exported functions should have doc comments with SaltStack API reference URLs
- Example:

  ```go
  // Login authenticates with SaltStack API
  // https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#login
  func (c *Client) Login(ctx context.Context) error { ... }
  ```

### Time Handling

- Custom `saltTime` type for unmarshaling SaltStack's time format
- Format: `"2006, Jan 02 15:04:05.000000"`

### Type Definitions

- Use `interface{}` for flexible JSON fields where necessary
- Prefer specific types when structure is known
- Use map types for dynamic data: `map[string]interface{}`, `map[string]string`

### Constants

- Client types: `LocalClient`, `RunnerClient`, `WheelClient`, etc.
- Target types: `Glob`, `Pcre`, `List`, `Grain`, `Pillar`, etc.
- Define constants as string types with meaningful names

### Package Organization

- Each file focuses on a specific feature (e.g., `auth.go`, `jobs.go`, `minions.go`)
- Client initialization in `client.go` and `client_options.go`
- Command/target types in `command_clients.go`
- Shared utilities in `utils.go`

## Testing Guidelines (When Adding Tests)

- Use `context.TODO()` or `context.Background()` in tests
- Mock HTTP responses for testing client methods
- Test error paths (invalid responses, network errors, auth failures)
- Use table-driven tests for multiple scenarios
- Keep tests independent and don't rely on external SaltStack instances

## Dependencies

- `github.com/tidwall/gjson` - Fast JSON parsing for specific field access
- `encoding/json` - Standard JSON encoding/decoding
- `net/http` - HTTP client and server
