package seq_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestMap(t *testing.T) {
	got := seq.Of(1, 2, 3).Map(strconv.Itoa).Collect()
	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Map().Collect() = %v, want %v", got, want)
	}
}

func TestFilterMap(t *testing.T) {
	got := seq.Of("1", "bad", "3").FilterMap(func(value string) (int, bool) {
		parsed, err := strconv.Atoi(value)
		return parsed, err == nil
	}).Collect()
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterMap().Collect() = %v, want %v", got, want)
	}
}

func TestFlatMap(t *testing.T) {
	got := seq.Of(1, 3).FlatMap(func(value int) seq.Seq[int] {
		return seq.Range(0, value)
	}).Collect()
	want := []int{0, 0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FlatMap().Collect() = %v, want %v", got, want)
	}
}

func TestInspect(t *testing.T) {
	var inspected []int
	values := seq.Of(1, 2, 3).Inspect(func(value int) {
		inspected = append(inspected, value)
	})
	if inspected != nil {
		t.Fatalf("Inspect actions before iteration = %v, want nil", inspected)
	}

	if got, want := values.Take(2).Collect(), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Inspect().Take(2).Collect() = %v, want %v", got, want)
	}
	if want := []int{1, 2}; !reflect.DeepEqual(inspected, want) {
		t.Fatalf("Inspect actions = %v, want %v", inspected, want)
	}
}

func TestScan(t *testing.T) {
	reducerCalls := 0
	scanned := seq.Of(1, 2, 3).Scan(10, func(accumulator, value int) int {
		reducerCalls++
		return accumulator + value
	})
	if reducerCalls != 0 {
		t.Fatalf("Scan reducer called %d times before iteration, want 0", reducerCalls)
	}

	if got, want := scanned.Collect(), []int{11, 13, 16}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Scan().Collect() = %v, want %v", got, want)
	}
	if reducerCalls != 3 {
		t.Fatalf("Scan reducer called %d times, want 3", reducerCalls)
	}
	if got := seq.Empty[int]().Scan(10, func(accumulator, value int) int {
		return accumulator + value
	}).Collect(); len(got) != 0 {
		t.Fatalf("empty Scan().Collect() = %v, want empty", got)
	}
}

func TestEnumerate(t *testing.T) {
	enumerated := seq.Of("a", "b", "c").Enumerate()
	if got, want := enumerated.Keys().Collect(), []int{0, 1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Enumerate().Keys().Collect() = %v, want %v", got, want)
	}
	if got, want := enumerated.Values().Collect(), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Enumerate().Values().Collect() = %v, want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	for _, tt := range []struct {
		name string
		size int
		want [][]int
	}{
		{name: "exact and partial chunks", size: 2, want: [][]int{{1, 2}, {3, 4}, {5}}},
		{name: "larger than input", size: 10, want: [][]int{{1, 2, 3, 4, 5}}},
		{name: "zero", size: 0, want: nil},
		{name: "negative", size: -1, want: nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Chunk(seq.Of(1, 2, 3, 4, 5), tt.size).Collect()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Chunk(size %d).Collect() = %v, want %v", tt.size, got, tt.want)
			}
		})
	}

	chunks := seq.Chunk(seq.Of(1, 2, 3, 4), 2).Collect()
	chunks[0][0] = 99
	if got, want := chunks[1], []int{3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("mutating first chunk changed second to %v, want %v", got, want)
	}
}

func TestChunkIsLazyAndSupportsEarlyTermination(t *testing.T) {
	produced := 0
	values := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3, 4} {
			produced++
			if !yield(value) {
				return
			}
		}
	})
	chunked := seq.Chunk(values, 2)
	if produced != 0 {
		t.Fatalf("Chunk consumed %d values before iteration, want 0", produced)
	}

	if got, want := chunked.Take(1).Collect(), [][]int{{1, 2}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Chunk().Take(1).Collect() = %v, want %v", got, want)
	}
	if produced != 2 {
		t.Fatalf("Chunk().Take(1) consumed %d inputs, want 2", produced)
	}
}

func TestWindow(t *testing.T) {
	for _, tt := range []struct {
		name string
		size int
		want [][]int
	}{
		{name: "sliding windows", size: 3, want: [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}},
		{name: "single values", size: 1, want: [][]int{{1}, {2}, {3}, {4}, {5}}},
		{name: "larger than input", size: 10, want: nil},
		{name: "zero", size: 0, want: nil},
		{name: "negative", size: -1, want: nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Window(seq.Of(1, 2, 3, 4, 5), tt.size).Collect()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Window(size %d).Collect() = %v, want %v", tt.size, got, tt.want)
			}
		})
	}

	windows := seq.Window(seq.Of(1, 2, 3, 4), 3).Collect()
	windows[0][1] = 99
	if got, want := windows[1], []int{2, 3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("mutating first window changed second to %v, want %v", got, want)
	}
}

func TestWindowIsLazyAndSupportsEarlyTermination(t *testing.T) {
	produced := 0
	values := seq.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3, 4} {
			produced++
			if !yield(value) {
				return
			}
		}
	})
	windowed := seq.Window(values, 3)
	if produced != 0 {
		t.Fatalf("Window consumed %d values before iteration, want 0", produced)
	}

	if got, want := windowed.Take(1).Collect(), [][]int{{1, 2, 3}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Window().Take(1).Collect() = %v, want %v", got, want)
	}
	if produced != 3 {
		t.Fatalf("Window().Take(1) consumed %d inputs, want 3", produced)
	}
}
