# go-errstack

[![CI](https://github.com/philiprehberger/go-errstack/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-errstack/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-errstack.svg)](https://pkg.go.dev/github.com/philiprehberger/go-errstack)
[![Last updated](https://img.shields.io/github/last-commit/philiprehberger/go-errstack)](https://github.com/philiprehberger/go-errstack/commits/main)

Error wrapping with stack traces for Go

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

### Annotations

Attach key-value metadata to errors without changing the message:

```go
err := errstack.WithValue(err, "request_id", "abc-123")
err = errstack.WithValue(err, "user", "alice")

val, ok := errstack.Value(err, "request_id") // "abc-123", true
```

### Caller

Capture a single stack frame at a given depth:

```go
f := errstack.Caller(0) // frame of the current function
fmt.Println(f)           // main.handleRequest (/app/server.go:18)
```

### Frame filtering

Trim stack frames to focus on relevant packages:

```go
frames := errstack.Stack(err)
frames = errstack.TrimAbove(frames, "myapp/handler") // remove frames above handler
frames = errstack.TrimBelow(frames, "myapp/handler") // remove frames below handler
```

### Formatted stack trace

Get a formatted multi-line stack trace string:

```go
fmt.Println(errstack.StackString(err))
// main.doWork
//     /path/to/file.go:42
// main.main
//     /path/to/file.go:15
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
| `StackString(err)` | Returns a formatted multi-line stack trace string |
| `Caller(skip)` | Returns a single stack frame at the given skip depth |
| `WithValue(err, key, val)` | Wraps an error with a key-value annotation |
| `Value(err, key)` | Extracts an annotation value from the error chain |
| `TrimAbove(frames, pkg)` | Removes frames above the first occurrence of pkg |
| `TrimBelow(frames, pkg)` | Removes frames below the last occurrence of pkg |

## Development

```bash
go test ./...
go vet ./...
```

## Support

If you find this project useful:

⭐ [Star the repo](https://github.com/philiprehberger/go-errstack)

🐛 [Report issues](https://github.com/philiprehberger/go-errstack/issues?q=is%3Aissue+is%3Aopen+label%3Abug)

💡 [Suggest features](https://github.com/philiprehberger/go-errstack/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement)

❤️ [Sponsor development](https://github.com/sponsors/philiprehberger)

🌐 [All Open Source Projects](https://philiprehberger.com/open-source-packages)

💻 [GitHub Profile](https://github.com/philiprehberger)

🔗 [LinkedIn Profile](https://www.linkedin.com/in/philiprehberger)

## License

[MIT](LICENSE)
