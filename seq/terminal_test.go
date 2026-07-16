package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

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
