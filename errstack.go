// Package errstack provides error wrapping with stack traces for Go.
package errstack

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Frame represents a single stack frame.
type Frame struct {
	Function string
	File     string
	Line     int
}

// String formats the frame as "Function (File:Line)".
func (f Frame) String() string {
	return fmt.Sprintf("%s (%s:%d)", f.Function, f.File, f.Line)
}

// stackError is an error that carries a stack trace and an optional cause.
type stackError struct {
	msg    string
	cause  error
	frames []Frame
}

// Error returns the error message. If a cause is present, the message is
// formatted as "msg: cause". If msg is empty, the cause's message is returned
// directly.
func (e *stackError) Error() string {
	if e.cause != nil {
		if e.msg != "" {
			return e.msg + ": " + e.cause.Error()
		}
		return e.cause.Error()
	}
	return e.msg
}

// Unwrap returns the underlying cause of the error, enabling compatibility
// with errors.Is, errors.As, and errors.Unwrap.
func (e *stackError) Unwrap() error {
	return e.cause
}

// Wrap wraps err with a stack trace captured at the call site. Returns nil if
// err is nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return &stackError{
		cause:  err,
		frames: capture(2),
	}
}

// Wrapf wraps err with a formatted message and a stack trace captured at the
// call site. Returns nil if err is nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &stackError{
		msg:    fmt.Sprintf(format, args...),
		cause:  err,
		frames: capture(2),
	}
}

// New creates a new error with the given message and a stack trace captured at
// the call site.
func New(msg string) error {
	return &stackError{
		msg:    msg,
		frames: capture(2),
	}
}

// Newf creates a new error with a formatted message and a stack trace captured
// at the call site.
func Newf(format string, args ...any) error {
	return &stackError{
		msg:    fmt.Sprintf(format, args...),
		frames: capture(2),
	}
}

// Stack extracts stack frames from an error. It walks the Unwrap chain looking
// for an error that carries stack information. Returns nil if no stack frames
// are found.
func Stack(err error) []Frame {
	for err != nil {
		if se, ok := err.(*stackError); ok {
			return se.frames
		}
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return nil
		}
		err = unwrapper.Unwrap()
	}
	return nil
}

// Caller returns a single stack frame at the given skip depth. Skip 0 refers
// to the caller of Caller itself.
func Caller(skip int) Frame {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return Frame{}
	}
	frame, _ := runtime.CallersFrames(pcs[:1]).Next()
	return Frame{
		Function: frame.Function,
		File:     frame.File,
		Line:     frame.Line,
	}
}

// annotatedError is an error that carries a key-value annotation.
type annotatedError struct {
	err error
	key string
	val any
}

// Error returns the underlying error message.
func (a *annotatedError) Error() string {
	return a.err.Error()
}

// Unwrap returns the underlying error.
func (a *annotatedError) Unwrap() error {
	return a.err
}

// WithValue wraps err with a key-value annotation. The original error and its
// chain are preserved. Returns nil if err is nil.
func WithValue(err error, key string, val any) error {
	if err == nil {
		return nil
	}
	return &annotatedError{
		err: err,
		key: key,
		val: val,
	}
}

// Value extracts an annotation value from the error chain by key. It walks the
// chain using errors.As looking for an annotatedError with a matching key.
// Returns the value and true if found, or nil and false otherwise.
func Value(err error, key string) (any, bool) {
	var ae *annotatedError
	for {
		if !errors.As(err, &ae) {
			return nil, false
		}
		if ae.key == key {
			return ae.val, true
		}
		err = ae.err
	}
}

// StackString returns a formatted multi-line stack trace string from an error.
// Each frame is rendered as two lines: the function name, then the file and
// line indented with a tab. Returns an empty string if no stack frames are
// found.
func StackString(err error) string {
	frames := Stack(err)
	if len(frames) == 0 {
		return ""
	}
	var b strings.Builder
	for i, f := range frames {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(f.Function)
		b.WriteByte('\n')
		fmt.Fprintf(&b, "\t%s:%d", f.File, f.Line)
	}
	return b.String()
}

// TrimAbove removes frames from packages above the specified package in the
// stack. It finds the first frame whose Function contains pkg and returns that
// frame and all frames below it (i.e., frames from the first occurrence
// onward).
func TrimAbove(frames []Frame, pkg string) []Frame {
	for i, f := range frames {
		if strings.Contains(f.Function, pkg) {
			return frames[i:]
		}
	}
	return nil
}

// TrimBelow removes frames from packages below the specified package in the
// stack. It finds the last frame whose Function contains pkg and returns all
// frames up to and including it.
func TrimBelow(frames []Frame, pkg string) []Frame {
	last := -1
	for i, f := range frames {
		if strings.Contains(f.Function, pkg) {
			last = i
		}
	}
	if last < 0 {
		return nil
	}
	return frames[:last+1]
}

// capture collects stack frames, skipping the given number of callers to
// exclude internal frames.
func capture(skip int) []Frame {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	if n == 0 {
		return nil
	}

	frames := make([]Frame, 0, n)
	iter := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := iter.Next()
		frames = append(frames, Frame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}
	return frames
}
