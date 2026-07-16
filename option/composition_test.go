package option_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/everettxwzhu/go-toolkit/option"
	"github.com/everettxwzhu/go-toolkit/result"
)

func TestFold(t *testing.T) {
	someCalls := 0
	noneCalls := 0
	onSome := func(value int) string {
		someCalls++
		return fmt.Sprintf("some:%d", value)
	}
	onNone := func() string {
		noneCalls++
		return "none"
	}

	if got := option.Some(42).Fold(onSome, onNone); got != "some:42" {
		t.Fatalf("Some(42).Fold() = %q, want %q", got, "some:42")
	}
	if someCalls != 1 || noneCalls != 0 {
		t.Fatalf("Some Fold calls = (%d, %d), want (1, 0)", someCalls, noneCalls)
	}

	if got := option.None[int]().Fold(onSome, onNone); got != "none" {
		t.Fatalf("None[int]().Fold() = %q, want %q", got, "none")
	}
	if someCalls != 1 || noneCalls != 1 {
		t.Fatalf("None Fold cumulative calls = (%d, %d), want (1, 1)", someCalls, noneCalls)
	}
}

func TestZipWith(t *testing.T) {
	zipped := option.Some(2).ZipWith(option.Some("x"), func(left int, right string) string {
		return fmt.Sprintf("%d%s", left, right)
	})
	if value, ok := zipped.Get(); value != "2x" || !ok {
		t.Fatalf("Some.ZipWith(Some).Get() = (%q, %t), want (\"2x\", true)", value, ok)
	}

	for _, tt := range []struct {
		name  string
		left  option.Option[int]
		right option.Option[string]
	}{
		{name: "left None", left: option.None[int](), right: option.Some("x")},
		{name: "right None", left: option.Some(2), right: option.None[string]()},
		{name: "both None", left: option.None[int](), right: option.None[string]()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			got := tt.left.ZipWith(tt.right, func(int, string) string {
				called = true
				return "unexpected"
			})
			if _, ok := got.Get(); ok {
				t.Fatal("ZipWith().IsSome() = true, want false")
			}
			if called {
				t.Fatal("ZipWith combine called when an input was None")
			}
		})
	}
}

func TestOr(t *testing.T) {
	if value, ok := option.Some(1).Or(option.Some(2)).Get(); value != 1 || !ok {
		t.Fatalf("Some(1).Or(Some(2)).Get() = (%d, %t), want (1, true)", value, ok)
	}
	if value, ok := option.None[int]().Or(option.Some(2)).Get(); value != 2 || !ok {
		t.Fatalf("None[int]().Or(Some(2)).Get() = (%d, %t), want (2, true)", value, ok)
	}
	if _, ok := option.None[int]().Or(option.None[int]()).Get(); ok {
		t.Fatal("None[int]().Or(None[int]()).IsSome() = true, want false")
	}
}

func TestOrGet(t *testing.T) {
	calls := 0
	fallback := func() option.Option[int] {
		calls++
		return option.Some(2)
	}

	if value, ok := option.Some(1).OrGet(fallback).Get(); value != 1 || !ok {
		t.Fatalf("Some(1).OrGet().Get() = (%d, %t), want (1, true)", value, ok)
	}
	if calls != 0 {
		t.Fatalf("fallback called %d times for Some, want 0", calls)
	}

	if value, ok := option.None[int]().OrGet(fallback).Get(); value != 2 || !ok {
		t.Fatalf("None[int]().OrGet().Get() = (%d, %t), want (2, true)", value, ok)
	}
	if calls != 1 {
		t.Fatalf("fallback called %d times after None, want 1", calls)
	}
}

func TestInspect(t *testing.T) {
	inspected := 0
	o := option.Some(42)
	got := o.Inspect(func(value int) {
		inspected = value
	})

	if value, ok := got.Get(); value != 42 || !ok {
		t.Fatalf("Some(42).Inspect().Get() = (%d, %t), want (42, true)", value, ok)
	}
	if inspected != 42 {
		t.Fatalf("inspected value = %d, want 42", inspected)
	}

	called := false
	option.None[int]().Inspect(func(int) {
		called = true
	})
	if called {
		t.Fatal("Inspect action called for None")
	}
}

func TestFlatten(t *testing.T) {
	if value, ok := option.Flatten(option.Some(option.Some(42))).Get(); value != 42 || !ok {
		t.Fatalf("Flatten(Some(Some(42))).Get() = (%d, %t), want (42, true)", value, ok)
	}
	if _, ok := option.Flatten(option.Some(option.None[int]())).Get(); ok {
		t.Fatal("Flatten(Some(None)).IsSome() = true, want false")
	}
	if _, ok := option.Flatten(option.None[option.Option[int]]()).Get(); ok {
		t.Fatal("Flatten(None).IsSome() = true, want false")
	}
}

func TestToResult(t *testing.T) {
	if value, err := option.Some(42).ToResult(nil).Get(); value != 42 || err != nil {
		t.Fatalf("Some(42).ToResult(nil).Get() = (%d, %v), want (42, nil)", value, err)
	}

	wantErr := errors.New("missing")
	if value, err := option.None[int]().ToResult(wantErr).Get(); value != 0 || err != wantErr {
		t.Fatalf("None[int]().ToResult().Get() = (%d, %v), want (0, %v)", value, err, wantErr)
	}
}

func TestToResultGet(t *testing.T) {
	wantErr := errors.New("missing")
	calls := 0
	errFallback := func() error {
		calls++
		return wantErr
	}

	if value, err := option.Some(42).ToResultGet(errFallback).Get(); value != 42 || err != nil {
		t.Fatalf("Some(42).ToResultGet().Get() = (%d, %v), want (42, nil)", value, err)
	}
	if calls != 0 {
		t.Fatalf("error fallback called %d times for Some, want 0", calls)
	}

	if value, err := option.None[int]().ToResultGet(errFallback).Get(); value != 0 || err != wantErr {
		t.Fatalf("None[int]().ToResultGet().Get() = (%d, %v), want (0, %v)", value, err, wantErr)
	}
	if calls != 1 {
		t.Fatalf("error fallback called %d times after None, want 1", calls)
	}
}

func TestFromResult(t *testing.T) {
	if value, ok := option.FromResult(result.Ok(42)).Get(); value != 42 || !ok {
		t.Fatalf("FromResult(Ok(42)).Get() = (%d, %t), want (42, true)", value, ok)
	}
	if _, ok := option.FromResult(result.Err[int](errors.New("failed"))).Get(); ok {
		t.Fatal("FromResult(Err).IsSome() = true, want false")
	}
}
