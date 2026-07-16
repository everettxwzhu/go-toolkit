package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestFilter(t *testing.T) {
	got := seq.Range(1, 6).Filter(func(value int) bool { return value%2 == 0 }).Collect()
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter().Collect() = %v, want %v", got, want)
	}
}

func TestTake(t *testing.T) {
	for _, tt := range []struct {
		name string
		n    int
		want []int
	}{
		{name: "positive", n: 2, want: []int{1, 2}},
		{name: "larger than sequence", n: 5, want: []int{1, 2, 3}},
		{name: "zero", n: 0, want: nil},
		{name: "negative", n: -1, want: nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := seq.Of(1, 2, 3).Take(tt.n).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Take(%d).Collect() = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestDrop(t *testing.T) {
	for _, tt := range []struct {
		name string
		n    int
		want []int
	}{
		{name: "positive", n: 2, want: []int{3}},
		{name: "larger than sequence", n: 5, want: nil},
		{name: "zero", n: 0, want: []int{1, 2, 3}},
		{name: "negative", n: -1, want: []int{1, 2, 3}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := seq.Of(1, 2, 3).Drop(tt.n).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Drop(%d).Collect() = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	got := seq.Of(1, 2, 3, 1).TakeWhile(func(value int) bool { return value < 3 }).Collect()
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TakeWhile().Collect() = %v, want %v", got, want)
	}
}

func TestDropWhile(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3, 1).DropWhile(func(value int) bool {
		calls++
		return value < 3
	}).Collect()
	want := []int{3, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DropWhile().Collect() = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

func TestDistinctBy(t *testing.T) {
	got := seq.Of("one", "two", "six", "three").DistinctBy(func(value string) int {
		return len(value)
	}).Collect()
	want := []string{"one", "three"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DistinctBy().Collect() = %v, want %v", got, want)
	}
}
