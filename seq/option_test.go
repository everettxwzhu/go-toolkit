package seq_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/option"
	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestFilterOption(t *testing.T) {
	calls := 0
	filtered := seq.Of(1, 2, 3, 4).FilterOption(func(value int) option.Option[string] {
		calls++
		if value%2 != 0 {
			return option.None[string]()
		}
		return option.Some(strconv.Itoa(value))
	})
	if calls != 0 {
		t.Fatalf("FilterOption transform called %d times before iteration, want 0", calls)
	}

	if got, want := filtered.Collect(), []string{"2", "4"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterOption().Collect() = %v, want %v", got, want)
	}
	if calls != 4 {
		t.Fatalf("FilterOption transform called %d times, want 4", calls)
	}
}

func TestFilterOptionSupportsEarlyTermination(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3, 4).FilterOption(func(value int) option.Option[int] {
		calls++
		if value%2 != 0 {
			return option.None[int]()
		}
		return option.Some(value)
	}).Take(1).Collect()

	if want := []int{2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterOption().Take(1).Collect() = %v, want %v", got, want)
	}
	if calls != 2 {
		t.Fatalf("FilterOption transform called %d times, want 2", calls)
	}
}

func TestCollectOption(t *testing.T) {
	got := seq.Of("1", "2", "3").CollectOption(func(value string) option.Option[int] {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return option.None[int]()
		}
		return option.Some(parsed)
	})
	if values, ok := got.Get(); !ok || !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("CollectOption(all Some).Get() = (%v, %t), want ([1 2 3], true)", values, ok)
	}

	empty, ok := seq.Empty[int]().CollectOption(func(value int) option.Option[int] {
		return option.Some(value)
	}).Get()
	if !ok || empty == nil || len(empty) != 0 {
		t.Fatalf("empty CollectOption().Get() = (%#v, %t), want (non-nil empty slice, true)", empty, ok)
	}
}

func TestCollectOptionStopsOnNone(t *testing.T) {
	var visited []int
	got := seq.Of(1, 2, 3, 4).CollectOption(func(value int) option.Option[int] {
		visited = append(visited, value)
		if value == 3 {
			return option.None[int]()
		}
		return option.Some(value * 10)
	})

	if got.IsSome() {
		t.Fatal("CollectOption(with None).IsSome() = true, want false")
	}
	if want := []int{1, 2, 3}; !reflect.DeepEqual(visited, want) {
		t.Fatalf("CollectOption visited = %v, want %v", visited, want)
	}
}

func TestFirstOption(t *testing.T) {
	if got, ok := seq.Of(3, 4).FirstOption().Get(); !ok || got != 3 {
		t.Fatalf("FirstOption().Get() = (%d, %t), want (3, true)", got, ok)
	}
	if seq.Empty[int]().FirstOption().IsSome() {
		t.Fatal("empty FirstOption().IsSome() = true, want false")
	}
}

func TestLastOption(t *testing.T) {
	if got, ok := seq.Of(3, 4).LastOption().Get(); !ok || got != 4 {
		t.Fatalf("LastOption().Get() = (%d, %t), want (4, true)", got, ok)
	}
	if seq.Empty[int]().LastOption().IsSome() {
		t.Fatal("empty LastOption().IsSome() = true, want false")
	}
}

func TestFindOption(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3, 4).FindOption(func(value int) bool {
		calls++
		return value == 3
	})
	if value, ok := got.Get(); !ok || value != 3 {
		t.Fatalf("FindOption().Get() = (%d, %t), want (3, true)", value, ok)
	}
	if calls != 3 {
		t.Fatalf("FindOption predicate called %d times, want 3", calls)
	}

	if seq.Of(1, 2).FindOption(func(value int) bool { return value > 2 }).IsSome() {
		t.Fatal("non-matching FindOption().IsSome() = true, want false")
	}
}
