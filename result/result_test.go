package result_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/result"
)

func TestResultZeroValue(t *testing.T) {
	var r result.Result[int]

	if value, err := r.Get(); value != 0 || err != nil {
		t.Fatalf("zero Result.Get() = (%d, %v), want (0, nil)", value, err)
	}
	if !r.IsOK() {
		t.Fatal("zero Result.IsOK() = false, want true")
	}
	if r.IsErr() {
		t.Fatal("zero Result.IsErr() = true, want false")
	}
}

func TestOk(t *testing.T) {
	r := result.Ok(42)

	if value, err := r.Get(); value != 42 || err != nil {
		t.Fatalf("Ok(42).Get() = (%d, %v), want (42, nil)", value, err)
	}
	if !r.IsOK() {
		t.Fatal("Ok(42).IsOK() = false, want true")
	}
	if r.IsErr() {
		t.Fatal("Ok(42).IsErr() = true, want false")
	}
}

func TestErr(t *testing.T) {
	wantErr := errors.New("failed")
	r := result.Err[int](wantErr)

	if value, err := r.Get(); value != 0 || err != wantErr {
		t.Fatalf("Err[int]().Get() = (%d, %v), want (0, %v)", value, err, wantErr)
	}
	if r.IsOK() {
		t.Fatal("Err[int]().IsOK() = true, want false")
	}
	if !r.IsErr() {
		t.Fatal("Err[int]().IsErr() = false, want true")
	}
}

func TestErrPanicsWithNil(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("Err[int](nil) did not panic")
		}
	}()

	result.Err[int](nil)
}

func TestFrom(t *testing.T) {
	if value, err := result.From("value", nil).Get(); value != "value" || err != nil {
		t.Fatalf("From(\"value\", nil).Get() = (%q, %v), want (\"value\", nil)", value, err)
	}

	wantErr := errors.New("failed")
	if value, err := result.From("value", wantErr).Get(); value != "value" || err != wantErr {
		t.Fatalf("From(\"value\", wantErr).Get() = (%q, %v), want (\"value\", %v)", value, err, wantErr)
	}
}

func TestOrElse(t *testing.T) {
	if got := result.Ok(1).OrElse(2); got != 1 {
		t.Fatalf("Ok(1).OrElse(2) = %d, want 1", got)
	}
	if got := result.Err[int](errors.New("failed")).OrElse(2); got != 2 {
		t.Fatalf("Err[int]().OrElse(2) = %d, want 2", got)
	}
}

func TestOrElseGet(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0
	fallback := func(err error) int {
		calls++
		if err != wantErr {
			t.Fatalf("fallback error = %v, want %v", err, wantErr)
		}
		return 2
	}

	if got := result.Ok(1).OrElseGet(fallback); got != 1 {
		t.Fatalf("Ok(1).OrElseGet() = %d, want 1", got)
	}
	if calls != 0 {
		t.Fatalf("fallback called %d times for Ok, want 0", calls)
	}

	if got := result.Err[int](wantErr).OrElseGet(fallback); got != 2 {
		t.Fatalf("Err[int]().OrElseGet() = %d, want 2", got)
	}
	if calls != 1 {
		t.Fatalf("fallback called %d times after Err, want 1", calls)
	}
}

func TestMap(t *testing.T) {
	mapped := result.Ok(42).Map(strconv.Itoa)
	if value, err := mapped.Get(); value != "42" || err != nil {
		t.Fatalf("Ok(42).Map().Get() = (%q, %v), want (\"42\", nil)", value, err)
	}

	wantErr := errors.New("failed")
	called := false
	failed := result.Err[int](wantErr).Map(func(value int) string {
		called = true
		return strconv.Itoa(value)
	})
	if value, err := failed.Get(); value != "" || err != wantErr {
		t.Fatalf("Err[int]().Map().Get() = (%q, %v), want (\"\", %v)", value, err, wantErr)
	}
	if called {
		t.Fatal("Map transform called for Err")
	}
}

func TestFlatMap(t *testing.T) {
	mapped := result.Ok("42").FlatMap(func(value string) result.Result[int] {
		return result.From(strconv.Atoi(value))
	})
	if value, err := mapped.Get(); value != 42 || err != nil {
		t.Fatalf("Ok(\"42\").FlatMap().Get() = (%d, %v), want (42, nil)", value, err)
	}

	mappedErr := result.Ok("bad").FlatMap(func(value string) result.Result[int] {
		return result.From(strconv.Atoi(value))
	})
	if _, err := mappedErr.Get(); err == nil {
		t.Fatal("Ok(\"bad\").FlatMap().Get() error = nil, want non-nil")
	}

	wantErr := errors.New("failed")
	called := false
	failed := result.Err[string](wantErr).FlatMap(func(string) result.Result[int] {
		called = true
		return result.Ok(1)
	})
	if value, err := failed.Get(); value != 0 || err != wantErr {
		t.Fatalf("Err[string]().FlatMap().Get() = (%d, %v), want (0, %v)", value, err, wantErr)
	}
	if called {
		t.Fatal("FlatMap transform called for Err")
	}
}
