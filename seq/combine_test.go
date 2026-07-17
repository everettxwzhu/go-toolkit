package seq_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
	"github.com/everettxwzhu/go-toolkit/tuple"
)

func TestConcat(t *testing.T) {
	got := seq.Of(1, 2).Concat(seq.Empty[int](), seq.Of(3, 4)).Collect()
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Concat().Collect() = %v, want %v", got, want)
	}
}

func TestConcatWithoutOthers(t *testing.T) {
	got := seq.Of(1, 2).Concat().Collect()
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Concat().Collect() = %v, want %v", got, want)
	}
}

func TestZipWith(t *testing.T) {
	leftProduced := 0
	left := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1} {
			leftProduced++
			if !yield(value) {
				return
			}
		}
	})
	rightProduced := 0
	right := seq.FromSeq(func(yield func(string) bool) {
		for _, value := range []string{"a", "b"} {
			rightProduced++
			if !yield(value) {
				return
			}
		}
	})

	combineCalls := 0
	zipped := left.ZipWith(right, func(left int, right string) string {
		combineCalls++
		return strconv.Itoa(left) + right
	})
	if leftProduced != 0 || rightProduced != 0 || combineCalls != 0 {
		t.Fatal("ZipWith consumed input before iteration")
	}

	if got, want := zipped.Collect(), []string{"1a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ZipWith().Collect() = %v, want %v", got, want)
	}
	if leftProduced != 1 || rightProduced != 1 || combineCalls != 1 {
		t.Fatalf(
			"ZipWith counts = left %d, right %d, combine %d; want 1, 1, 1",
			leftProduced,
			rightProduced,
			combineCalls,
		)
	}
}

func TestZip(t *testing.T) {
	got := seq.Zip(seq.Of(1, 2, 3), seq.Of("one", "two")).Collect()
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Zip().Collect() = %v, want %v", got, want)
	}

	leftProduced := 0
	left := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			leftProduced++
			if !yield(value) {
				return
			}
		}
	})
	rightProduced := 0
	right := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{10, 20, 30} {
			rightProduced++
			if !yield(value) {
				return
			}
		}
	})

	seq.Zip(left, right).Take(1).Collect()
	if leftProduced != 1 || rightProduced != 1 {
		t.Fatalf("Zip().Take(1) produced (%d, %d) inputs, want (1, 1)", leftProduced, rightProduced)
	}
}

func TestIntersperse(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values seq.Seq[int]
		want   []int
	}{
		{name: "empty", values: seq.Empty[int](), want: nil},
		{name: "single", values: seq.Of(1), want: []int{1}},
		{name: "multiple", values: seq.Of(1, 2, 3), want: []int{1, 0, 2, 0, 3}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.values.Intersperse(0).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Intersperse(0).Collect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrepend(t *testing.T) {
	got := seq.Of(3, 4).Prepend(1, 2).Collect()
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Prepend().Collect() = %v, want %v", got, want)
	}
}

func TestAppend(t *testing.T) {
	got := seq.Of(1, 2).Append(3, 4).Collect()
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Append().Collect() = %v, want %v", got, want)
	}
}
