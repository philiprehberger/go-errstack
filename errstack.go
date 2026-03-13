// Package errstack provides error wrapping with stack traces for Go.
package errstack

import (
	"fmt"
	"runtime"
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
