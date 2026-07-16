package option_test

import (
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/option"
)

func TestOptionZeroValue(t *testing.T) {
	var o option.Option[int]

	if value, ok := o.Get(); value != 0 || ok {
		t.Fatalf("zero Option.Get() = (%d, %t), want (0, false)", value, ok)
	}
	if o.IsSome() {
		t.Fatal("zero Option.IsSome() = true, want false")
	}
	if !o.IsNone() {
		t.Fatal("zero Option.IsNone() = false, want true")
	}
}

func TestSome(t *testing.T) {
	o := option.Some(42)

	if value, ok := o.Get(); value != 42 || !ok {
		t.Fatalf("Some(42).Get() = (%d, %t), want (42, true)", value, ok)
	}
	if !o.IsSome() {
		t.Fatal("Some(42).IsSome() = false, want true")
	}
	if o.IsNone() {
		t.Fatal("Some(42).IsNone() = true, want false")
	}
}

func TestNone(t *testing.T) {
	o := option.None[string]()

	if value, ok := o.Get(); value != "" || ok {
		t.Fatalf("None[string]().Get() = (%q, %t), want (\"\", false)", value, ok)
	}
}

func TestFrom(t *testing.T) {
	if value, ok := option.From("value", true).Get(); value != "value" || !ok {
		t.Fatalf("From(\"value\", true).Get() = (%q, %t), want (\"value\", true)", value, ok)
	}
	if value, ok := option.From("ignored", false).Get(); value != "ignored" || ok {
		t.Fatalf("From(\"ignored\", false).Get() = (%q, %t), want (\"ignored\", false)", value, ok)
	}
}

func TestOrElse(t *testing.T) {
	if got := option.Some(1).OrElse(2); got != 1 {
		t.Fatalf("Some(1).OrElse(2) = %d, want 1", got)
	}
	if got := option.None[int]().OrElse(2); got != 2 {
		t.Fatalf("None[int]().OrElse(2) = %d, want 2", got)
	}
}

func TestOrElseGet(t *testing.T) {
	calls := 0
	fallback := func() int {
		calls++
		return 2
	}

	if got := option.Some(1).OrElseGet(fallback); got != 1 {
		t.Fatalf("Some(1).OrElseGet() = %d, want 1", got)
	}
	if calls != 0 {
		t.Fatalf("fallback called %d times for Some, want 0", calls)
	}

	if got := option.None[int]().OrElseGet(fallback); got != 2 {
		t.Fatalf("None[int]().OrElseGet() = %d, want 2", got)
	}
	if calls != 1 {
		t.Fatalf("fallback called %d times after None, want 1", calls)
	}
}

func TestMap(t *testing.T) {
	mapped := option.Some(42).Map(strconv.Itoa)
	if value, ok := mapped.Get(); value != "42" || !ok {
		t.Fatalf("Some(42).Map().Get() = (%q, %t), want (\"42\", true)", value, ok)
	}

	called := false
	none := option.None[int]().Map(func(value int) string {
		called = true
		return strconv.Itoa(value)
	})
	if value, ok := none.Get(); value != "" || ok {
		t.Fatalf("None[int]().Map().Get() = (%q, %t), want (\"\", false)", value, ok)
	}
	if called {
		t.Fatal("Map transform called for None")
	}
}

func TestFlatMap(t *testing.T) {
	mapped := option.Some("42").FlatMap(func(value string) option.Option[int] {
		parsed, err := strconv.Atoi(value)
		return option.From(parsed, err == nil)
	})
	if value, ok := mapped.Get(); value != 42 || !ok {
		t.Fatalf("Some(\"42\").FlatMap().Get() = (%d, %t), want (42, true)", value, ok)
	}

	missing := option.Some("bad").FlatMap(func(value string) option.Option[int] {
		parsed, err := strconv.Atoi(value)
		return option.From(parsed, err == nil)
	})
	if value, ok := missing.Get(); value != 0 || ok {
		t.Fatalf("Some(\"bad\").FlatMap().Get() = (%d, %t), want (0, false)", value, ok)
	}

	called := false
	option.None[string]().FlatMap(func(string) option.Option[int] {
		called = true
		return option.Some(1)
	})
	if called {
		t.Fatal("FlatMap transform called for None")
	}
}

func TestFilter(t *testing.T) {
	if value, ok := option.Some(2).Filter(func(value int) bool { return value%2 == 0 }).Get(); value != 2 || !ok {
		t.Fatalf("Some(2).Filter().Get() = (%d, %t), want (2, true)", value, ok)
	}
	if value, ok := option.Some(3).Filter(func(value int) bool { return value%2 == 0 }).Get(); value != 0 || ok {
		t.Fatalf("Some(3).Filter().Get() = (%d, %t), want (0, false)", value, ok)
	}

	called := false
	option.None[int]().Filter(func(int) bool {
		called = true
		return true
	})
	if called {
		t.Fatal("Filter predicate called for None")
	}
}
