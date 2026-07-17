package seq_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
	"github.com/everettxwzhu/go-toolkit/tuple"
)

func seq2FromPairs[K, V any](pairs ...tuple.Pair[K, V]) seq.Seq2[K, V] {
	return seq.FromSeq2(func(yield func(K, V) bool) {
		for _, pair := range pairs {
			if !yield(pair.First, pair.Second) {
				return
			}
		}
	})
}

func collectSeq2[K, V any](values seq.Seq2[K, V]) []tuple.Pair[K, V] {
	var pairs []tuple.Pair[K, V]
	for key, value := range values.All() {
		pairs = append(pairs, tuple.New(key, value))
	}
	return pairs
}

func TestSeq2ZeroAndEmpty(t *testing.T) {
	var zero seq.Seq2[int, string]
	if got := collectSeq2(zero); len(got) != 0 {
		t.Fatalf("zero Seq2 pairs = %v, want empty", got)
	}
	if got := collectSeq2(seq.Empty2[int, string]()); len(got) != 0 {
		t.Fatalf("Empty2 pairs = %v, want empty", got)
	}
	if got := collectSeq2(seq.FromSeq2[int, string](nil)); len(got) != 0 {
		t.Fatalf("FromSeq2(nil) pairs = %v, want empty", got)
	}
}

func TestFromMap(t *testing.T) {
	values := map[string]int{"one": 1, "two": 2}
	if got := seq.CollectMap(seq.FromMap(values)); !reflect.DeepEqual(got, values) {
		t.Fatalf("CollectMap(FromMap()) = %v, want %v", got, values)
	}
}

func TestSeq2Map(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("one", 1),
		tuple.New("two", 2),
	)
	got := values.Map(func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	}).Collect()
	want := []string{"one:1", "two:2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Seq2.Map().Collect() = %v, want %v", got, want)
	}
}

func TestSeq2MapKeys(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("one", 1),
		tuple.New("two", 2),
	)
	got := collectSeq2(values.MapKeys(func(key string, value int) int {
		return len(key) + value
	}))
	want := []tuple.Pair[int, int]{
		tuple.New(4, 1),
		tuple.New(5, 2),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MapKeys() pairs = %v, want %v", got, want)
	}
}

func TestSeq2MapValues(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("one", 1),
		tuple.New("two", 2),
	)
	got := collectSeq2(values.MapValues(func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	}))
	want := []tuple.Pair[string, string]{
		tuple.New("one", "one:1"),
		tuple.New("two", "two:2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MapValues() pairs = %v, want %v", got, want)
	}
}

func TestSeq2FilterIsLazyAndSupportsEarlyTermination(t *testing.T) {
	produced := 0
	values := seq.FromSeq2(func(yield func(int, string) bool) {
		for _, pair := range []tuple.Pair[int, string]{
			tuple.New(1, "one"),
			tuple.New(2, "two"),
			tuple.New(3, "three"),
		} {
			produced++
			if !yield(pair.First, pair.Second) {
				return
			}
		}
	})
	predicateCalls := 0
	filtered := values.Filter(func(key int, _ string) bool {
		predicateCalls++
		return key%2 == 0
	})
	if produced != 0 || predicateCalls != 0 {
		t.Fatal("Seq2.Filter consumed input before iteration")
	}

	got := collectSeq2(filtered)
	want := []tuple.Pair[int, string]{tuple.New(2, "two")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter() pairs = %v, want %v", got, want)
	}
	if produced != 3 || predicateCalls != 3 {
		t.Fatalf("Filter counts = produced %d, predicate %d; want 3, 3", produced, predicateCalls)
	}

	produced = 0
	for range values.Filter(func(int, string) bool { return true }).All() {
		break
	}
	if produced != 1 {
		t.Fatalf("breaking Seq2.Filter consumed %d inputs, want 1", produced)
	}
}

func TestSeq2Inspect(t *testing.T) {
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	)
	var inspected []tuple.Pair[int, string]
	result := values.Inspect(func(key int, value string) {
		inspected = append(inspected, tuple.New(key, value))
	})
	if inspected != nil {
		t.Fatalf("Inspect actions before iteration = %v, want nil", inspected)
	}

	var got []tuple.Pair[int, string]
	for key, value := range result.All() {
		got = append(got, tuple.New(key, value))
		if len(got) == 2 {
			break
		}
	}
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Inspect pairs = %v, want %v", got, want)
	}
	if !reflect.DeepEqual(inspected, want) {
		t.Fatalf("Inspect actions = %v, want %v", inspected, want)
	}
}

func TestSeq2KeysValuesAndSwap(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("one", 1),
		tuple.New("two", 2),
	)
	if got, want := values.Keys().Collect(), []string{"one", "two"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys().Collect() = %v, want %v", got, want)
	}
	if got, want := values.Values().Collect(), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Values().Collect() = %v, want %v", got, want)
	}

	got := collectSeq2(values.Swap())
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Swap() pairs = %v, want %v", got, want)
	}
}

func TestSeq2ForEach(t *testing.T) {
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	)
	var got []tuple.Pair[int, string]
	values.ForEach(func(key int, value string) {
		got = append(got, tuple.New(key, value))
	})
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ForEach pairs = %v, want %v", got, want)
	}
}

func TestCollectMapLastValueWins(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("key", 1),
		tuple.New("other", 2),
		tuple.New("key", 3),
	)
	got := seq.CollectMap(values)
	want := map[string]int{"key": 3, "other": 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CollectMap() = %v, want %v", got, want)
	}
}
