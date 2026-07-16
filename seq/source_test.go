package seq_test

import (
	"iter"
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestEmpty(t *testing.T) {
	if got := seq.Empty[int]().Collect(); len(got) != 0 {
		t.Fatalf("Empty().Collect() = %v, want empty slice", got)
	}
}

func TestFromSeq(t *testing.T) {
	source := iter.Seq[int](func(yield func(int) bool) {
		yield(1)
		yield(2)
	})
	if got, want := seq.FromSeq(source).Collect(), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("FromSeq().Collect() = %v, want %v", got, want)
	}

	if got := seq.FromSeq[int](nil).Collect(); len(got) != 0 {
		t.Fatalf("FromSeq(nil).Collect() = %v, want empty slice", got)
	}
}

func TestFromSliceReadsValuesLazily(t *testing.T) {
	values := []int{1, 2}
	s := seq.FromSlice(values)
	values[0] = 9

	if got, want := s.Collect(), []int{9, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("FromSlice().Collect() = %v, want %v", got, want)
	}
}

func TestOf(t *testing.T) {
	if got, want := seq.Of("a", "b").Collect(), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Of().Collect() = %v, want %v", got, want)
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
		want       []int
	}{
		{name: "ascending", start: 2, end: 5, want: []int{2, 3, 4}},
		{name: "equal bounds", start: 2, end: 2, want: nil},
		{name: "descending bounds", start: 5, end: 2, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seq.Range(tt.start, tt.end).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Range(%d, %d).Collect() = %v, want %v", tt.start, tt.end, got, tt.want)
			}
		})
	}
}
