# go-errstack

[![CI](https://github.com/philiprehberger/go-errstack/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-errstack/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-errstack.svg)](https://pkg.go.dev/github.com/philiprehberger/go-errstack)
[![License](https://img.shields.io/github/license/philiprehberger/go-errstack)](LICENSE)

Error wrapping with stack traces for Go.

## Installation

```bash
go get github.com/philiprehberger/go-errstack
```

## Usage

```go
import "github.com/philiprehberger/go-errstack"
```

### Wrapping an existing error

```go
file, err := os.Open("config.json")
if err != nil {
    return errstack.Wrap(err)
}
```

### Wrapping with a message

```go
user, err := db.FindUser(id)
if err != nil {
    return errstack.Wrapf(err, "failed to find user %d", id)
}
```

### Creating a new error

```go
if name == "" {
    return errstack.New("name must not be empty")
}
```

### Extracting the stack trace

```go
if err != nil {
    frames := errstack.Stack(err)
    for _, f := range frames {
        fmt.Println(f) // main.handleRequest (/app/server.go:42)
    }
}
```

### Compatible with errors.Is and errors.As

```go
var ErrNotFound = errors.New("not found")

err := errstack.Wrap(ErrNotFound)
errors.Is(err, ErrNotFound) // true
```

## API

| Function / Type | Description |
|-----------------|-------------|
| `Frame` | A single stack frame with Function, File, and Line fields |
| `Frame.String()` | Formats the frame as "Function (File:Line)" |
| `Wrap(err)` | Wraps an error with a stack trace; returns nil if err is nil |
| `Wrapf(err, fmt, args...)` | Wraps an error with a formatted message and stack trace |
| `New(msg)` | Creates a new error with a stack trace |
| `Newf(fmt, args...)` | Creates a new formatted error with a stack trace |
| `Stack(err)` | Extracts stack frames from an error; returns nil if none found |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
