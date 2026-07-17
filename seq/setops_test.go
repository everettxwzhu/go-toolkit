package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestUnionBy(t *testing.T) {
	got := seq.Of("one", "two", "three").UnionBy(
		seq.Of("six", "four", "seven"),
		func(value string) int { return len(value) },
	).Collect()
	want := []string{"one", "three", "four"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UnionBy().Collect() = %v, want %v", got, want)
	}
}

func TestUnionByIsLazyAndSupportsEarlyTermination(t *testing.T) {
	leftProduced := 0
	left := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 1, 2, 3} {
			leftProduced++
			if !yield(value) {
				return
			}
		}
	})
	rightProduced := 0
	right := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{4, 5} {
			rightProduced++
			if !yield(value) {
				return
			}
		}
	})

	union := left.UnionBy(right, func(value int) int { return value })
	if leftProduced != 0 || rightProduced != 0 {
		t.Fatal("UnionBy consumed input before iteration")
	}
	if got, want := union.Take(2).Collect(), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("UnionBy().Take(2).Collect() = %v, want %v", got, want)
	}
	if leftProduced != 3 || rightProduced != 0 {
		t.Fatalf(
			"UnionBy().Take(2) produced (%d, %d) inputs, want (3, 0)",
			leftProduced,
			rightProduced,
		)
	}
}

func TestIntersectBy(t *testing.T) {
	got := seq.Of("one", "two", "three", "six", "four").IntersectBy(
		seq.Of("xxx", "zzzz", "also four"),
		func(value string) int { return len(value) },
	).Collect()
	want := []string{"one", "four"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IntersectBy().Collect() = %v, want %v", got, want)
	}
}

func TestExceptBy(t *testing.T) {
	got := seq.Of("one", "two", "three", "six", "four", "five").ExceptBy(
		seq.Of("xxx"),
		func(value string) int { return len(value) },
	).Collect()
	want := []string{"three", "four"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExceptBy().Collect() = %v, want %v", got, want)
	}
}

func TestContainsBy(t *testing.T) {
	calls := 0
	got := seq.Of("one", "three", "four").ContainsBy("seven", func(left, right string) bool {
		calls++
		return len(left) == len(right)
	})
	if !got || calls != 2 {
		t.Fatalf("ContainsBy() = %t with %d calls, want true with 2 calls", got, calls)
	}

	if seq.Empty[string]().ContainsBy("value", func(string, string) bool {
		t.Fatal("ContainsBy equal called for empty sequence")
		return true
	}) {
		t.Fatal("empty ContainsBy() = true, want false")
	}
}

func TestEqualBy(t *testing.T) {
	calls := 0
	if !seq.Of("one", "four", "three").EqualBy(
		seq.Of("two", "five", "seven"),
		func(left, right string) bool {
			calls++
			return len(left) == len(right)
		},
	) {
		t.Fatal("EqualBy(equal lengths) = false, want true")
	}
	if calls != 3 {
		t.Fatalf("EqualBy equal called %d times, want 3", calls)
	}

	calls = 0
	if seq.Of(1, 2, 3).EqualBy(seq.Of(1, 9, 3), func(left, right int) bool {
		calls++
		return left == right
	}) {
		t.Fatal("EqualBy(mismatch) = true, want false")
	}
	if calls != 2 {
		t.Fatalf("mismatching EqualBy equal called %d times, want 2", calls)
	}

	if seq.Of(1, 2).EqualBy(seq.Of(1), func(left, right int) bool { return left == right }) {
		t.Fatal("EqualBy(longer left) = true, want false")
	}
	if seq.Of(1).EqualBy(seq.Of(1, 2), func(left, right int) bool { return left == right }) {
		t.Fatal("EqualBy(longer right) = true, want false")
	}
	if !seq.Empty[int]().EqualBy(seq.Empty[int](), func(left, right int) bool { return left == right }) {
		t.Fatal("EqualBy(empty, empty) = false, want true")
	}
}

func TestDistinct(t *testing.T) {
	got := seq.Distinct(seq.Of(1, 2, 1, 3, 2)).Collect()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Distinct().Collect() = %v, want %v", got, want)
	}
}

func TestContains(t *testing.T) {
	if !seq.Contains(seq.Of("one", "two"), "two") {
		t.Fatal(`Contains(["one", "two"], "two") = false, want true`)
	}
	if seq.Contains(seq.Of("one", "two"), "three") {
		t.Fatal(`Contains(["one", "two"], "three") = true, want false`)
	}
}

func TestEqual(t *testing.T) {
	if !seq.Equal(seq.Of(1, 2, 3), seq.Of(1, 2, 3)) {
		t.Fatal("Equal(equal sequences) = false, want true")
	}
	if seq.Equal(seq.Of(1, 2, 3), seq.Of(1, 3, 2)) {
		t.Fatal("Equal(different order) = true, want false")
	}
	if seq.Equal(seq.Of(1, 2), seq.Of(1, 2, 3)) {
		t.Fatal("Equal(different lengths) = true, want false")
	}
}
