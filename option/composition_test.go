package option_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/option"
	"github.com/everettxwzhu/go-toolkit/result"
	"github.com/everettxwzhu/go-toolkit/tuple"
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

func TestZip(t *testing.T) {
	got := option.Zip(option.Some(1), option.Some("one"))
	if value, ok := got.Get(); !ok || value != tuple.New(1, "one") {
		t.Fatalf("Zip(Some, Some).Get() = (%v, %t), want (%v, true)", value, ok, tuple.New(1, "one"))
	}

	if option.Zip(option.None[int](), option.Some("one")).IsSome() {
		t.Fatal("Zip(None, Some).IsSome() = true, want false")
	}
	if option.Zip(option.Some(1), option.None[string]()).IsSome() {
		t.Fatal("Zip(Some, None).IsSome() = true, want false")
	}
}

func TestExists(t *testing.T) {
	if !option.Some(4).Exists(func(value int) bool { return value%2 == 0 }) {
		t.Fatal("Some(4).Exists(even) = false, want true")
	}
	if option.Some(3).Exists(func(value int) bool { return value%2 == 0 }) {
		t.Fatal("Some(3).Exists(even) = true, want false")
	}

	called := false
	if option.None[int]().Exists(func(int) bool {
		called = true
		return true
	}) {
		t.Fatal("None.Exists() = true, want false")
	}
	if called {
		t.Fatal("Exists predicate called for None")
	}
}

func TestContainsBy(t *testing.T) {
	equal := func(left, right string) bool {
		return len(left) == len(right)
	}
	if !option.Some("one").ContainsBy("two", equal) {
		t.Fatal(`Some("one").ContainsBy("two") = false, want true`)
	}
	if option.Some("one").ContainsBy("four", equal) {
		t.Fatal(`Some("one").ContainsBy("four") = true, want false`)
	}

	called := false
	if option.None[string]().ContainsBy("one", func(string, string) bool {
		called = true
		return true
	}) {
		t.Fatal("None.ContainsBy() = true, want false")
	}
	if called {
		t.Fatal("ContainsBy equal called for None")
	}
}

func TestToSlice(t *testing.T) {
	if got, want := option.Some(0).ToSlice(), []int{0}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Some(0).ToSlice() = %v, want %v", got, want)
	}
	if got := option.None[int]().ToSlice(); got == nil || len(got) != 0 {
		t.Fatalf("None[int]().ToSlice() = %#v, want a non-nil empty slice", got)
	}
}

func TestTranspose(t *testing.T) {
	someOK := option.Transpose(option.Some(result.Ok(42)))
	if value, err := someOK.Get(); err != nil {
		t.Fatalf("Transpose(Some(Ok)).Get() error = %v, want nil", err)
	} else if got, ok := value.Get(); !ok || got != 42 {
		t.Fatalf("transposed value.Get() = (%d, %t), want (42, true)", got, ok)
	}

	wantErr := errors.New("failed")
	if _, err := option.Transpose(option.Some(result.Err[int](wantErr))).Get(); err != wantErr {
		t.Fatalf("Transpose(Some(Err)).Get() error = %v, want %v", err, wantErr)
	}

	none := option.Transpose(option.None[result.Result[int]]())
	if value, err := none.Get(); err != nil {
		t.Fatalf("Transpose(None).Get() error = %v, want nil", err)
	} else if value.IsSome() {
		t.Fatal("Transpose(None) contains Some, want None")
	}
}

func TestSequence(t *testing.T) {
	got := option.Sequence([]option.Option[int]{
		option.Some(1),
		option.Some(2),
		option.Some(3),
	})
	if values, ok := got.Get(); !ok || !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("Sequence(all Some).Get() = (%v, %t), want ([1 2 3], true)", values, ok)
	}

	if option.Sequence([]option.Option[int]{
		option.Some(1),
		option.None[int](),
		option.Some(3),
	}).IsSome() {
		t.Fatal("Sequence(with None).IsSome() = true, want false")
	}

	empty, ok := option.Sequence([]option.Option[int]{}).Get()
	if !ok || empty == nil || len(empty) != 0 {
		t.Fatalf("Sequence(empty).Get() = (%#v, %t), want (non-nil empty slice, true)", empty, ok)
	}
}

func TestTraverse(t *testing.T) {
	var visited []int
	got := option.Traverse([]int{1, 2, 3}, func(value int) option.Option[string] {
		visited = append(visited, value)
		return option.Some(fmt.Sprintf("v%d", value))
	})
	if values, ok := got.Get(); !ok || !reflect.DeepEqual(values, []string{"v1", "v2", "v3"}) {
		t.Fatalf("Traverse(all Some).Get() = (%v, %t), want ([v1 v2 v3], true)", values, ok)
	}
	if !reflect.DeepEqual(visited, []int{1, 2, 3}) {
		t.Fatalf("Traverse visit order = %v, want [1 2 3]", visited)
	}

	visited = nil
	failed := option.Traverse([]int{1, 2, 3, 4}, func(value int) option.Option[int] {
		visited = append(visited, value)
		if value == 3 {
			return option.None[int]()
		}
		return option.Some(value * 10)
	})
	if failed.IsSome() {
		t.Fatal("Traverse(with None).IsSome() = true, want false")
	}
	if !reflect.DeepEqual(visited, []int{1, 2, 3}) {
		t.Fatalf("failed Traverse visit order = %v, want [1 2 3]", visited)
	}
}

func TestContains(t *testing.T) {
	if !option.Contains(option.Some("value"), "value") {
		t.Fatal(`Contains(Some("value"), "value") = false, want true`)
	}
	if option.Contains(option.Some("value"), "other") {
		t.Fatal(`Contains(Some("value"), "other") = true, want false`)
	}
	if option.Contains(option.None[string](), "") {
		t.Fatal(`Contains(None, "") = true, want false`)
	}
}
