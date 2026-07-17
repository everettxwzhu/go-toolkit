package result_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/result"
	"github.com/everettxwzhu/go-toolkit/tuple"
)

type codedError struct {
	code int
}

func (e *codedError) Error() string {
	return fmt.Sprintf("code %d", e.code)
}

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

func TestZip(t *testing.T) {
	got := result.Zip(result.Ok(1), result.Ok("one"))
	if value, err := got.Get(); err != nil || value != tuple.New(1, "one") {
		t.Fatalf("Zip(Ok, Ok).Get() = (%v, %v), want (%v, nil)", value, err, tuple.New(1, "one"))
	}

	leftErr := errors.New("left")
	rightErr := errors.New("right")
	for _, tt := range []struct {
		name    string
		left    result.Result[int]
		right   result.Result[string]
		wantErr error
	}{
		{name: "left Err", left: result.Err[int](leftErr), right: result.Ok("one"), wantErr: leftErr},
		{name: "right Err", left: result.Ok(1), right: result.Err[string](rightErr), wantErr: rightErr},
		{name: "both Err", left: result.Err[int](leftErr), right: result.Err[string](rightErr), wantErr: leftErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := result.Zip(tt.left, tt.right).Get(); err != tt.wantErr {
				t.Fatalf("Zip().Get() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsure(t *testing.T) {
	ok := result.Ok(4)
	errCalls := 0
	got := ok.Ensure(
		func(value int) bool { return value%2 == 0 },
		func(int) error {
			errCalls++
			return errors.New("unexpected")
		},
	)
	if value, err := got.Get(); value != 4 || err != nil {
		t.Fatalf("Ok(4).Ensure(true).Get() = (%d, %v), want (4, nil)", value, err)
	}
	if errCalls != 0 {
		t.Fatalf("Ensure error factory called %d times for matching value, want 0", errCalls)
	}

	wantErr := errors.New("not positive")
	got = result.Ok(-1).Ensure(
		func(value int) bool { return value > 0 },
		func(value int) error {
			if value != -1 {
				t.Fatalf("Ensure error factory value = %d, want -1", value)
			}
			return wantErr
		},
	)
	if _, err := got.Get(); err != wantErr {
		t.Fatalf("Ok(-1).Ensure(false).Get() error = %v, want %v", err, wantErr)
	}

	originalErr := errors.New("original")
	predicateCalled := false
	errorFactoryCalled := false
	failed := result.From(7, originalErr).Ensure(
		func(int) bool {
			predicateCalled = true
			return true
		},
		func(int) error {
			errorFactoryCalled = true
			return errors.New("unexpected")
		},
	)
	if value, err := failed.Get(); value != 7 || err != originalErr {
		t.Fatalf("failed Ensure().Get() = (%d, %v), want (7, %v)", value, err, originalErr)
	}
	if predicateCalled || errorFactoryCalled {
		t.Fatalf("failed Ensure callbacks called = (%t, %t), want (false, false)", predicateCalled, errorFactoryCalled)
	}
}

func TestContainsBy(t *testing.T) {
	equal := func(left, right string) bool {
		return len(left) == len(right)
	}
	if !result.Ok("one").ContainsBy("two", equal) {
		t.Fatal(`Ok("one").ContainsBy("two") = false, want true`)
	}
	if result.Ok("one").ContainsBy("four", equal) {
		t.Fatal(`Ok("one").ContainsBy("four") = true, want false`)
	}

	called := false
	if result.Err[string](errors.New("failed")).ContainsBy("one", func(string, string) bool {
		called = true
		return true
	}) {
		t.Fatal("Err.ContainsBy() = true, want false")
	}
	if called {
		t.Fatal("ContainsBy equal called for Err")
	}
}

func TestSequence(t *testing.T) {
	got := result.Sequence([]result.Result[int]{
		result.Ok(1),
		result.Ok(2),
		result.Ok(3),
	})
	if values, err := got.Get(); err != nil || !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("Sequence(all Ok).Get() = (%v, %v), want ([1 2 3], nil)", values, err)
	}

	wantErr := errors.New("failed")
	if _, err := result.Sequence([]result.Result[int]{
		result.Ok(1),
		result.Err[int](wantErr),
		result.Err[int](errors.New("later")),
	}).Get(); err != wantErr {
		t.Fatalf("Sequence(with Err).Get() error = %v, want %v", err, wantErr)
	}

	empty, err := result.Sequence([]result.Result[int]{}).Get()
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("Sequence(empty).Get() = (%#v, %v), want (non-nil empty slice, nil)", empty, err)
	}
}

func TestTraverse(t *testing.T) {
	var visited []int
	got := result.Traverse([]int{1, 2, 3}, func(value int) result.Result[string] {
		visited = append(visited, value)
		return result.Ok(fmt.Sprintf("v%d", value))
	})
	if values, err := got.Get(); err != nil || !reflect.DeepEqual(values, []string{"v1", "v2", "v3"}) {
		t.Fatalf("Traverse(all Ok).Get() = (%v, %v), want ([v1 v2 v3], nil)", values, err)
	}
	if !reflect.DeepEqual(visited, []int{1, 2, 3}) {
		t.Fatalf("Traverse visit order = %v, want [1 2 3]", visited)
	}

	wantErr := errors.New("three")
	visited = nil
	failed := result.Traverse([]int{1, 2, 3, 4}, func(value int) result.Result[int] {
		visited = append(visited, value)
		if value == 3 {
			return result.Err[int](wantErr)
		}
		return result.Ok(value * 10)
	})
	if _, err := failed.Get(); err != wantErr {
		t.Fatalf("Traverse(with Err).Get() error = %v, want %v", err, wantErr)
	}
	if !reflect.DeepEqual(visited, []int{1, 2, 3}) {
		t.Fatalf("failed Traverse visit order = %v, want [1 2 3]", visited)
	}
}

func TestWrapErr(t *testing.T) {
	ok := result.Ok(42).WrapErr("context")
	if value, err := ok.Get(); value != 42 || err != nil {
		t.Fatalf("Ok.WrapErr().Get() = (%d, %v), want (42, nil)", value, err)
	}

	original := &codedError{code: 7}
	wrapped := result.Err[int](original).WrapErr("read value")
	_, err := wrapped.Get()
	if got, want := err.Error(), "read value: code 7"; got != want {
		t.Fatalf("WrapErr error = %q, want %q", got, want)
	}
	if !errors.Is(err, original) {
		t.Fatalf("errors.Is(%v, original) = false, want true", err)
	}
	var coded *codedError
	if !errors.As(err, &coded) || coded != original {
		t.Fatalf("errors.As(%v) = %v, want original coded error", err, coded)
	}
}
