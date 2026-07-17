package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

type rankedValue struct {
	name string
	rank int
}

func TestCollect(t *testing.T) {
	if got, want := seq.Of(1, 2, 3).Collect(), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Collect() = %v, want %v", got, want)
	}
}

func TestReduce(t *testing.T) {
	got := seq.Of(1, 2, 3).Reduce("0", func(acc string, value int) string {
		return acc + string(rune('0'+value))
	})
	if want := "0123"; got != want {
		t.Fatalf("Reduce() = %q, want %q", got, want)
	}
}

func TestFirst(t *testing.T) {
	if got, ok := seq.Of(3, 4).First(); !ok || got != 3 {
		t.Fatalf("First() = (%d, %t), want (3, true)", got, ok)
	}
	if got, ok := seq.Empty[int]().First(); ok || got != 0 {
		t.Fatalf("empty First() = (%d, %t), want (0, false)", got, ok)
	}
}

func TestLast(t *testing.T) {
	if got, ok := seq.Of(3, 4, 5).Last(); !ok || got != 5 {
		t.Fatalf("Last() = (%d, %t), want (5, true)", got, ok)
	}
	if got, ok := seq.Empty[int]().Last(); ok || got != 0 {
		t.Fatalf("empty Last() = (%d, %t), want (0, false)", got, ok)
	}
}

func TestFind(t *testing.T) {
	if got, ok := seq.Of(1, 2, 3).Find(func(value int) bool { return value%2 == 0 }); !ok || got != 2 {
		t.Fatalf("Find() = (%d, %t), want (2, true)", got, ok)
	}
	if got, ok := seq.Of(1, 3).Find(func(value int) bool { return value%2 == 0 }); ok || got != 0 {
		t.Fatalf("non-matching Find() = (%d, %t), want (0, false)", got, ok)
	}
}

func TestCount(t *testing.T) {
	if got := seq.Range(5, 9).Count(); got != 4 {
		t.Fatalf("Count() = %d, want 4", got)
	}
}

func TestAnyShortCircuits(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3).Any(func(value int) bool {
		calls++
		return value == 2
	})
	if !got || calls != 2 {
		t.Fatalf("Any() = %t with %d calls, want true with 2 calls", got, calls)
	}
}

func TestEvery(t *testing.T) {
	calls := 0
	got := seq.Of(2, 4, 5, 6).Every(func(value int) bool {
		calls++
		return value%2 == 0
	})
	if got || calls != 3 {
		t.Fatalf("Every() = %t with %d calls, want false with 3 calls", got, calls)
	}
	if !seq.Empty[int]().Every(func(int) bool { return false }) {
		t.Fatal("Every() on an empty sequence = false, want true")
	}
}

func TestForEach(t *testing.T) {
	var got []int
	seq.Of(1, 2, 3).ForEach(func(value int) { got = append(got, value) })
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ForEach values = %v, want %v", got, want)
	}
}

func TestAt(t *testing.T) {
	if got, ok := seq.Of(10, 20, 30).At(1); !ok || got != 20 {
		t.Fatalf("At(1) = (%d, %t), want (20, true)", got, ok)
	}
	if got, ok := seq.Of(10, 20, 30).At(3); ok || got != 0 {
		t.Fatalf("At(3) = (%d, %t), want (0, false)", got, ok)
	}

	iterated := false
	values := seq.FromSeq(func(yield func(int) bool) {
		iterated = true
		yield(1)
	})
	if got, ok := values.At(-1); ok || got != 0 {
		t.Fatalf("At(-1) = (%d, %t), want (0, false)", got, ok)
	}
	if iterated {
		t.Fatal("At(-1) iterated its source")
	}

	produced := 0
	values = seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{10, 20, 30} {
			produced++
			if !yield(value) {
				return
			}
		}
	})
	values.At(1)
	if produced != 2 {
		t.Fatalf("At(1) consumed %d values, want 2", produced)
	}
}

func TestSingle(t *testing.T) {
	if got, ok := seq.Of(42).Single(); !ok || got != 42 {
		t.Fatalf("single-value Single() = (%d, %t), want (42, true)", got, ok)
	}
	if got, ok := seq.Empty[int]().Single(); ok || got != 0 {
		t.Fatalf("empty Single() = (%d, %t), want (0, false)", got, ok)
	}

	produced := 0
	values := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			produced++
			if !yield(value) {
				return
			}
		}
	})
	if got, ok := values.Single(); ok || got != 0 {
		t.Fatalf("multi-value Single() = (%d, %t), want (0, false)", got, ok)
	}
	if produced != 2 {
		t.Fatalf("multi-value Single() consumed %d values, want 2", produced)
	}
}

func TestReduceFirst(t *testing.T) {
	got, ok := seq.Of(1, 2, 3, 4).ReduceFirst(func(left, right int) int {
		return left - right
	})
	if !ok || got != -8 {
		t.Fatalf("ReduceFirst(subtract) = (%d, %t), want (-8, true)", got, ok)
	}

	called := false
	if got, ok := seq.Empty[int]().ReduceFirst(func(left, right int) int {
		called = true
		return left + right
	}); ok || got != 0 {
		t.Fatalf("empty ReduceFirst() = (%d, %t), want (0, false)", got, ok)
	}
	if called {
		t.Fatal("ReduceFirst reducer called for an empty sequence")
	}
}

func TestPartition(t *testing.T) {
	matched, unmatched := seq.Of(1, 2, 3, 4, 5).Partition(func(value int) bool {
		return value%2 == 0
	})
	if want := []int{2, 4}; !reflect.DeepEqual(matched, want) {
		t.Fatalf("Partition matched = %v, want %v", matched, want)
	}
	if want := []int{1, 3, 5}; !reflect.DeepEqual(unmatched, want) {
		t.Fatalf("Partition unmatched = %v, want %v", unmatched, want)
	}
}

func TestMinBy(t *testing.T) {
	values := []rankedValue{
		{name: "middle", rank: 2},
		{name: "first minimum", rank: 1},
		{name: "second minimum", rank: 1},
	}
	keyCalls := 0
	got, ok := seq.Of(values...).MinBy(func(value rankedValue) int {
		keyCalls++
		return value.rank
	})
	if !ok || got != values[1] {
		t.Fatalf("MinBy() = (%v, %t), want (%v, true)", got, ok, values[1])
	}
	if keyCalls != len(values) {
		t.Fatalf("MinBy key called %d times, want %d", keyCalls, len(values))
	}

	if got, ok := seq.Empty[rankedValue]().MinBy(func(value rankedValue) int {
		return value.rank
	}); ok || got != (rankedValue{}) {
		t.Fatalf("empty MinBy() = (%v, %t), want (zero, false)", got, ok)
	}
}

func TestMaxBy(t *testing.T) {
	values := []rankedValue{
		{name: "middle", rank: 2},
		{name: "first maximum", rank: 3},
		{name: "second maximum", rank: 3},
	}
	keyCalls := 0
	got, ok := seq.Of(values...).MaxBy(func(value rankedValue) int {
		keyCalls++
		return value.rank
	})
	if !ok || got != values[1] {
		t.Fatalf("MaxBy() = (%v, %t), want (%v, true)", got, ok, values[1])
	}
	if keyCalls != len(values) {
		t.Fatalf("MaxBy key called %d times, want %d", keyCalls, len(values))
	}

	if got, ok := seq.Empty[rankedValue]().MaxBy(func(value rankedValue) int {
		return value.rank
	}); ok || got != (rankedValue{}) {
		t.Fatalf("empty MaxBy() = (%v, %t), want (zero, false)", got, ok)
	}
}

func TestMinFunc(t *testing.T) {
	compare := func(left, right rankedValue) int {
		return left.rank - right.rank
	}
	values := []rankedValue{
		{name: "high", rank: 3},
		{name: "first low", rank: 1},
		{name: "second low", rank: 1},
	}

	if got, ok := seq.Of(values...).MinFunc(compare); !ok || got != values[1] {
		t.Fatalf("MinFunc() = (%v, %t), want (%v, true)", got, ok, values[1])
	}
	if got, ok := seq.Empty[rankedValue]().MinFunc(compare); ok || got != (rankedValue{}) {
		t.Fatalf("empty MinFunc() = (%v, %t), want (zero, false)", got, ok)
	}
}

func TestMaxFunc(t *testing.T) {
	compare := func(left, right rankedValue) int {
		return left.rank - right.rank
	}
	values := []rankedValue{
		{name: "low", rank: 1},
		{name: "first high", rank: 3},
		{name: "second high", rank: 3},
	}

	if got, ok := seq.Of(values...).MaxFunc(compare); !ok || got != values[1] {
		t.Fatalf("MaxFunc() = (%v, %t), want (%v, true)", got, ok, values[1])
	}
	if got, ok := seq.Empty[rankedValue]().MaxFunc(compare); ok || got != (rankedValue{}) {
		t.Fatalf("empty MaxFunc() = (%v, %t), want (zero, false)", got, ok)
	}
}
