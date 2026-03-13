package errstack

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestWrapCapturesStack(t *testing.T) {
	err := Wrap(errors.New("base"))
	frames := Stack(err)
	if len(frames) == 0 {
		t.Fatal("expected stack frames, got none")
	}
	f := frames[0]
	if f.File == "" {
		t.Error("expected non-empty File in first frame")
	}
	if f.Line <= 0 {
		t.Errorf("expected Line > 0, got %d", f.Line)
	}
	if !strings.Contains(f.Function, "TestWrapCapturesStack") {
		t.Errorf("expected Function to contain test name, got %q", f.Function)
	}
}

func TestWrapNilReturnsNil(t *testing.T) {
	if err := Wrap(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestWrapfNilReturnsNil(t *testing.T) {
	if err := Wrapf(nil, "should not matter"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestWrapfAddsMessage(t *testing.T) {
	base := errors.New("base error")
	err := Wrapf(base, "context %d", 42)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	want := "context 42: base error"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestNewCreatesErrorWithStack(t *testing.T) {
	err := New("something went wrong")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Error() != "something went wrong" {
		t.Errorf("got %q, want %q", err.Error(), "something went wrong")
	}
	frames := Stack(err)
	if len(frames) == 0 {
		t.Fatal("expected stack frames from New")
	}
	if !strings.Contains(frames[0].Function, "TestNewCreatesErrorWithStack") {
		t.Errorf("expected first frame to contain test function name, got %q", frames[0].Function)
	}
}

func TestNewfCreatesFormattedError(t *testing.T) {
	err := Newf("failed with code %d: %s", 404, "not found")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	want := "failed with code 404: not found"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
	frames := Stack(err)
	if len(frames) == 0 {
		t.Fatal("expected stack frames from Newf")
	}
}

func TestStackExtractsFrames(t *testing.T) {
	base := errors.New("base")
	wrapped := Wrap(base)
	frames := Stack(wrapped)
	if frames == nil {
		t.Fatal("expected frames, got nil")
	}
	if len(frames) == 0 {
		t.Fatal("expected at least one frame")
	}
}

func TestStackReturnsNilForNonStackErrors(t *testing.T) {
	err := errors.New("plain error")
	frames := Stack(err)
	if frames != nil {
		t.Errorf("expected nil frames for plain error, got %v", frames)
	}
}

func TestStackReturnsNilForNil(t *testing.T) {
	frames := Stack(nil)
	if frames != nil {
		t.Errorf("expected nil frames for nil error, got %v", frames)
	}
}

// sentinel is a custom error type for testing errors.Is and errors.As.
type sentinel struct {
	code int
}

func (s *sentinel) Error() string {
	return fmt.Sprintf("sentinel(%d)", s.code)
}

func TestErrorsIsWorksThroughWrap(t *testing.T) {
	base := errors.New("target")
	wrapped := Wrap(base)
	if !errors.Is(wrapped, base) {
		t.Error("errors.Is should find base through Wrap")
	}
}

func TestErrorsIsWorksThroughWrapf(t *testing.T) {
	base := errors.New("target")
	wrapped := Wrapf(base, "adding context")
	if !errors.Is(wrapped, base) {
		t.Error("errors.Is should find base through Wrapf")
	}
}

func TestErrorsAsWorksThroughWrap(t *testing.T) {
	base := &sentinel{code: 42}
	wrapped := Wrap(base)
	var target *sentinel
	if !errors.As(wrapped, &target) {
		t.Fatal("errors.As should find sentinel through Wrap")
	}
	if target.code != 42 {
		t.Errorf("expected code 42, got %d", target.code)
	}
}

func TestErrorsAsWorksThroughWrapf(t *testing.T) {
	base := &sentinel{code: 99}
	wrapped := Wrapf(base, "wrapped")
	var target *sentinel
	if !errors.As(wrapped, &target) {
		t.Fatal("errors.As should find sentinel through Wrapf")
	}
	if target.code != 99 {
		t.Errorf("expected code 99, got %d", target.code)
	}
}

func TestUnwrapChain(t *testing.T) {
	base := errors.New("root cause")
	wrapped := Wrap(base)
	unwrapped := errors.Unwrap(wrapped)
	if unwrapped != base {
		t.Errorf("expected Unwrap to return base error, got %v", unwrapped)
	}
}

func TestErrorMessageWithCauseNoMsg(t *testing.T) {
	base := errors.New("underlying")
	err := Wrap(base)
	if err.Error() != "underlying" {
		t.Errorf("got %q, want %q", err.Error(), "underlying")
	}
}

func TestErrorMessageWithCauseAndMsg(t *testing.T) {
	base := errors.New("underlying")
	err := Wrapf(base, "context")
	want := "context: underlying"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestErrorMessageNoCause(t *testing.T) {
	err := New("standalone")
	if err.Error() != "standalone" {
		t.Errorf("got %q, want %q", err.Error(), "standalone")
	}
}

func TestFrameString(t *testing.T) {
	f := Frame{
		Function: "main.doStuff",
		File:     "/app/main.go",
		Line:     42,
	}
	want := "main.doStuff (/app/main.go:42)"
	if f.String() != want {
		t.Errorf("got %q, want %q", f.String(), want)
	}
}

func TestStackWalksUnwrapChain(t *testing.T) {
	base := errors.New("base")
	inner := Wrap(base)
	// Wrap the stack error inside a standard fmt.Errorf wrapper
	outer := fmt.Errorf("outer: %w", inner)
	frames := Stack(outer)
	if frames == nil {
		t.Fatal("expected Stack to find frames through fmt.Errorf wrapper")
	}
	if len(frames) == 0 {
		t.Error("expected at least one frame")
	}
}
