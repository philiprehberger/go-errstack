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

// --- Caller tests ---

func TestCallerReturnsCurrentFrame(t *testing.T) {
	f := Caller(0)
	if !strings.Contains(f.Function, "TestCallerReturnsCurrentFrame") {
		t.Errorf("expected Function to contain test name, got %q", f.Function)
	}
	if f.File == "" {
		t.Error("expected non-empty File")
	}
	if f.Line <= 0 {
		t.Errorf("expected Line > 0, got %d", f.Line)
	}
}

func callerHelper() Frame {
	return Caller(1)
}

func TestCallerSkipReturnsParent(t *testing.T) {
	f := callerHelper()
	if !strings.Contains(f.Function, "TestCallerSkipReturnsParent") {
		t.Errorf("expected Function to contain test name, got %q", f.Function)
	}
}

// --- WithValue / Value tests ---

func TestWithValueNilReturnsNil(t *testing.T) {
	if err := WithValue(nil, "key", "val"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestWithValuePreservesMessage(t *testing.T) {
	base := errors.New("original")
	err := WithValue(base, "user", "alice")
	if err.Error() != "original" {
		t.Errorf("got %q, want %q", err.Error(), "original")
	}
}

func TestValueExtractsAnnotation(t *testing.T) {
	base := errors.New("fail")
	err := WithValue(base, "request_id", "abc-123")
	val, ok := Value(err, "request_id")
	if !ok {
		t.Fatal("expected to find request_id annotation")
	}
	if val != "abc-123" {
		t.Errorf("got %v, want %q", val, "abc-123")
	}
}

func TestValueReturnsFalseForMissingKey(t *testing.T) {
	base := errors.New("fail")
	err := WithValue(base, "key", "val")
	_, ok := Value(err, "other")
	if ok {
		t.Error("expected ok to be false for missing key")
	}
}

func TestValueReturnsFalseForPlainError(t *testing.T) {
	err := errors.New("plain")
	_, ok := Value(err, "key")
	if ok {
		t.Error("expected ok to be false for plain error")
	}
}

func TestAnnotationChaining(t *testing.T) {
	base := errors.New("fail")
	err := WithValue(base, "user", "alice")
	err = WithValue(err, "request_id", "req-1")
	err = WithValue(err, "trace_id", "trace-2")

	val, ok := Value(err, "user")
	if !ok || val != "alice" {
		t.Errorf("user: got %v (%v), want alice", val, ok)
	}
	val, ok = Value(err, "request_id")
	if !ok || val != "req-1" {
		t.Errorf("request_id: got %v (%v), want req-1", val, ok)
	}
	val, ok = Value(err, "trace_id")
	if !ok || val != "trace-2" {
		t.Errorf("trace_id: got %v (%v), want trace-2", val, ok)
	}
}

func TestAnnotationOverrideReturnsNewest(t *testing.T) {
	base := errors.New("fail")
	err := WithValue(base, "key", "first")
	err = WithValue(err, "key", "second")
	val, ok := Value(err, "key")
	if !ok {
		t.Fatal("expected to find key")
	}
	if val != "second" {
		t.Errorf("got %v, want %q (should return outermost)", val, "second")
	}
}

func TestAnnotatedErrorUnwrapsCorrectly(t *testing.T) {
	base := errors.New("target")
	err := WithValue(base, "k", "v")
	if !errors.Is(err, base) {
		t.Error("errors.Is should find base through WithValue")
	}
}

// --- StackString tests ---

func TestStackStringFormatsFrames(t *testing.T) {
	err := New("test error")
	s := StackString(err)
	if s == "" {
		t.Fatal("expected non-empty stack string")
	}
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	// First line should be function name
	if !strings.Contains(lines[0], "TestStackStringFormatsFrames") {
		t.Errorf("expected first line to contain test name, got %q", lines[0])
	}
	// Second line should be indented file:line
	if !strings.HasPrefix(lines[1], "\t") {
		t.Errorf("expected second line to start with tab, got %q", lines[1])
	}
	if !strings.Contains(lines[1], ".go:") {
		t.Errorf("expected second line to contain file:line, got %q", lines[1])
	}
}

func TestStackStringReturnsEmptyForPlainError(t *testing.T) {
	err := errors.New("plain")
	s := StackString(err)
	if s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
}

func TestStackStringReturnsEmptyForNil(t *testing.T) {
	s := StackString(nil)
	if s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
}

// --- TrimAbove / TrimBelow tests ---

func TestTrimAboveFiltersFrames(t *testing.T) {
	frames := []Frame{
		{Function: "runtime.main"},
		{Function: "main.run"},
		{Function: "myapp/handler.Serve"},
		{Function: "myapp/handler.process"},
	}
	result := TrimAbove(frames, "myapp/handler")
	if len(result) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(result))
	}
	if result[0].Function != "myapp/handler.Serve" {
		t.Errorf("expected first frame to be handler.Serve, got %q", result[0].Function)
	}
}

func TestTrimAboveReturnsNilIfNoMatch(t *testing.T) {
	frames := []Frame{
		{Function: "runtime.main"},
		{Function: "main.run"},
	}
	result := TrimAbove(frames, "nonexistent")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTrimAboveEmptySlice(t *testing.T) {
	result := TrimAbove(nil, "pkg")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTrimBelowFiltersFrames(t *testing.T) {
	frames := []Frame{
		{Function: "myapp/handler.Serve"},
		{Function: "myapp/handler.process"},
		{Function: "net/http.serve"},
		{Function: "runtime.goexit"},
	}
	result := TrimBelow(frames, "myapp/handler")
	if len(result) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(result))
	}
	if result[1].Function != "myapp/handler.process" {
		t.Errorf("expected last frame to be handler.process, got %q", result[1].Function)
	}
}

func TestTrimBelowReturnsNilIfNoMatch(t *testing.T) {
	frames := []Frame{
		{Function: "runtime.main"},
	}
	result := TrimBelow(frames, "nonexistent")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTrimBelowEmptySlice(t *testing.T) {
	result := TrimBelow(nil, "pkg")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTrimAboveAndBelowCombined(t *testing.T) {
	frames := []Frame{
		{Function: "runtime.main"},
		{Function: "myapp/server.Start"},
		{Function: "myapp/handler.Serve"},
		{Function: "myapp/handler.validate"},
		{Function: "net/http.ListenAndServe"},
	}
	result := TrimAbove(frames, "myapp/handler")
	result = TrimBelow(result, "myapp/handler")
	if len(result) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(result))
	}
	if result[0].Function != "myapp/handler.Serve" {
		t.Errorf("got %q", result[0].Function)
	}
	if result[1].Function != "myapp/handler.validate" {
		t.Errorf("got %q", result[1].Function)
	}
}
