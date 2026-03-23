# Changelog

## 0.2.0

- Add `Caller` function to capture a single stack frame at a given depth
- Add `WithValue` and `Value` for attaching key-value annotations to errors
- Add `StackString` for formatted multi-line stack trace output
- Add `TrimAbove` and `TrimBelow` for filtering stack frames by package

## 0.1.2

- Consolidate README badges onto single line

## 0.1.1

- Add badges and Development section to README

## 0.1.0

- Initial release
- `Wrap` and `Wrapf` for adding stack traces to existing errors
- `New` and `Newf` for creating errors with stack traces
- `Stack` for extracting frames from wrapped errors
- Compatible with `errors.Is`, `errors.As`, and `errors.Unwrap`
