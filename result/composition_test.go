package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/everettxwzhu/go-toolkit/result"
)

func TestFold(t *testing.T) {
	wantErr := errors.New("failed")
	okCalls := 0
	errCalls := 0
	onOK := func(value int) string {
		okCalls++
		return fmt.Sprintf("ok:%d", value)
	}
	onErr := func(err error) string {
		errCalls++
		if err != wantErr {
			t.Fatalf("onErr error = %v, want %v", err, wantErr)
		}
		return "error"
	}

	if got := result.Ok(42).Fold(onOK, onErr); got != "ok:42" {
		t.Fatalf("Ok(42).Fold() = %q, want %q", got, "ok:42")
	}
	if okCalls != 1 || errCalls != 0 {
		t.Fatalf("Ok Fold calls = (%d, %d), want (1, 0)", okCalls, errCalls)
	}

	if got := result.Err[int](wantErr).Fold(onOK, onErr); got != "error" {
		t.Fatalf("Err[int]().Fold() = %q, want %q", got, "error")
	}
	if okCalls != 1 || errCalls != 1 {
		t.Fatalf("Err Fold cumulative calls = (%d, %d), want (1, 1)", okCalls, errCalls)
	}
}

func TestMapErr(t *testing.T) {
	called := false
	okResult := result.Ok(42).MapErr(func(error) error {
		called = true
		return errors.New("unexpected")
	})
	if value, err := okResult.Get(); value != 42 || err != nil {
		t.Fatalf("Ok(42).MapErr().Get() = (%d, %v), want (42, nil)", value, err)
	}
	if called {
		t.Fatal("MapErr transform called for Ok")
	}

	wantErr := errors.New("failed")
	mapped := result.Err[int](wantErr).MapErr(func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})
	if value, err := mapped.Get(); value != 0 || !errors.Is(err, wantErr) {
		t.Fatalf("Err[int]().MapErr().Get() = (%d, %v), want (0, wrapping %v)", value, err, wantErr)
	}
}

func TestZipWith(t *testing.T) {
	zipped := result.Ok(2).ZipWith(result.Ok("x"), func(left int, right string) string {
		return fmt.Sprintf("%d%s", left, right)
	})
	if value, err := zipped.Get(); value != "2x" || err != nil {
		t.Fatalf("Ok.ZipWith(Ok).Get() = (%q, %v), want (\"2x\", nil)", value, err)
	}

	leftErr := errors.New("left")
	rightErr := errors.New("right")
	for _, tt := range []struct {
		name    string
		left    result.Result[int]
		right   result.Result[string]
		wantErr error
	}{
		{name: "left Err", left: result.Err[int](leftErr), right: result.Ok("x"), wantErr: leftErr},
		{name: "right Err", left: result.Ok(2), right: result.Err[string](rightErr), wantErr: rightErr},
		{name: "both Err", left: result.Err[int](leftErr), right: result.Err[string](rightErr), wantErr: leftErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			got := tt.left.ZipWith(tt.right, func(int, string) string {
				called = true
				return "unexpected"
			})
			if value, err := got.Get(); value != "" || err != tt.wantErr {
				t.Fatalf("ZipWith().Get() = (%q, %v), want (\"\", %v)", value, err, tt.wantErr)
			}
			if called {
				t.Fatal("ZipWith combine called when an input was Err")
			}
		})
	}
}

func TestOr(t *testing.T) {
	originalErr := errors.New("original")
	fallbackErr := errors.New("fallback")

	if value, err := result.Ok(1).Or(result.Ok(2)).Get(); value != 1 || err != nil {
		t.Fatalf("Ok(1).Or(Ok(2)).Get() = (%d, %v), want (1, nil)", value, err)
	}
	if value, err := result.Err[int](originalErr).Or(result.Ok(2)).Get(); value != 2 || err != nil {
		t.Fatalf("Err.Or(Ok(2)).Get() = (%d, %v), want (2, nil)", value, err)
	}
	if value, err := result.Err[int](originalErr).Or(result.Err[int](fallbackErr)).Get(); value != 0 || err != fallbackErr {
		t.Fatalf("Err.Or(Err).Get() = (%d, %v), want (0, %v)", value, err, fallbackErr)
	}
}

func TestOrGet(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0
	fallback := func(err error) result.Result[int] {
		calls++
		if err != wantErr {
			t.Fatalf("fallback error = %v, want %v", err, wantErr)
		}
		return result.Ok(2)
	}

	if value, err := result.Ok(1).OrGet(fallback).Get(); value != 1 || err != nil {
		t.Fatalf("Ok(1).OrGet().Get() = (%d, %v), want (1, nil)", value, err)
	}
	if calls != 0 {
		t.Fatalf("fallback called %d times for Ok, want 0", calls)
	}

	if value, err := result.Err[int](wantErr).OrGet(fallback).Get(); value != 2 || err != nil {
		t.Fatalf("Err[int]().OrGet().Get() = (%d, %v), want (2, nil)", value, err)
	}
	if calls != 1 {
		t.Fatalf("fallback called %d times after Err, want 1", calls)
	}
}

func TestInspect(t *testing.T) {
	inspected := 0
	got := result.Ok(42).Inspect(func(value int) {
		inspected = value
	})
	if value, err := got.Get(); value != 42 || err != nil {
		t.Fatalf("Ok(42).Inspect().Get() = (%d, %v), want (42, nil)", value, err)
	}
	if inspected != 42 {
		t.Fatalf("inspected value = %d, want 42", inspected)
	}

	called := false
	result.Err[int](errors.New("failed")).Inspect(func(int) {
		called = true
	})
	if called {
		t.Fatal("Inspect action called for Err")
	}
}

func TestInspectErr(t *testing.T) {
	wantErr := errors.New("failed")
	var inspected error
	got := result.Err[int](wantErr).InspectErr(func(err error) {
		inspected = err
	})
	if value, err := got.Get(); value != 0 || err != wantErr {
		t.Fatalf("Err[int]().InspectErr().Get() = (%d, %v), want (0, %v)", value, err, wantErr)
	}
	if inspected != wantErr {
		t.Fatalf("inspected error = %v, want %v", inspected, wantErr)
	}

	called := false
	result.Ok(42).InspectErr(func(error) {
		called = true
	})
	if called {
		t.Fatal("InspectErr action called for Ok")
	}
}

func TestFlatten(t *testing.T) {
	if value, err := result.Flatten(result.Ok(result.Ok(42))).Get(); value != 42 || err != nil {
		t.Fatalf("Flatten(Ok(Ok(42))).Get() = (%d, %v), want (42, nil)", value, err)
	}

	innerErr := errors.New("inner")
	if value, err := result.Flatten(result.Ok(result.Err[int](innerErr))).Get(); value != 0 || err != innerErr {
		t.Fatalf("Flatten(Ok(Err)).Get() = (%d, %v), want (0, %v)", value, err, innerErr)
	}

	outerErr := errors.New("outer")
	nested := result.From(result.Err[int](innerErr), outerErr)
	if value, err := result.Flatten(nested).Get(); value != 0 || err != outerErr {
		t.Fatalf("Flatten(Err).Get() = (%d, %v), want (0, %v)", value, err, outerErr)
	}
}
