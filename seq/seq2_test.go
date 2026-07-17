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

func TestSeq2MapPairs(t *testing.T) {
	values := seq2FromPairs(
		tuple.New("one", 1),
		tuple.New("two", 2),
	)
	calls := 0
	mapped := values.MapPairs(func(key string, value int) (int, string) {
		calls++
		return len(key), key + ":" + strconv.Itoa(value)
	})
	if calls != 0 {
		t.Fatalf("MapPairs transform called %d times before iteration, want 0", calls)
	}

	got := collectSeq2(mapped)
	want := []tuple.Pair[int, string]{
		tuple.New(3, "one:1"),
		tuple.New(3, "two:2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MapPairs() pairs = %v, want %v", got, want)
	}
	if calls != 2 {
		t.Fatalf("MapPairs transform called %d times, want 2", calls)
	}
}

func TestSeq2MapPairsSupportsEarlyTermination(t *testing.T) {
	calls := 0
	mapped := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	).MapPairs(func(key int, value string) (string, int) {
		calls++
		return value, key
	})

	for range mapped.All() {
		break
	}
	if calls != 1 {
		t.Fatalf("breaking MapPairs iteration called transform %d times, want 1", calls)
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

func TestSeq2Take(t *testing.T) {
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	)
	for _, tt := range []struct {
		name string
		n    int
		want []tuple.Pair[int, string]
	}{
		{
			name: "positive",
			n:    2,
			want: []tuple.Pair[int, string]{
				tuple.New(1, "one"),
				tuple.New(2, "two"),
			},
		},
		{
			name: "larger than sequence",
			n:    5,
			want: []tuple.Pair[int, string]{
				tuple.New(1, "one"),
				tuple.New(2, "two"),
				tuple.New(3, "three"),
			},
		},
		{name: "zero", n: 0},
		{name: "negative", n: -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectSeq2(values.Take(tt.n)); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Take(%d) pairs = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestSeq2TakeIsLazyAndDoesNotOverconsume(t *testing.T) {
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
	taken := values.Take(2)
	if produced != 0 {
		t.Fatalf("Take consumed %d pairs before iteration, want 0", produced)
	}

	if got := collectSeq2(taken); len(got) != 2 {
		t.Fatalf("len(Take(2).Collect()) = %d, want 2", len(got))
	}
	if produced != 2 {
		t.Fatalf("Take(2) consumed %d pairs, want 2", produced)
	}
}

func TestSeq2Drop(t *testing.T) {
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	)
	for _, tt := range []struct {
		name string
		n    int
		want []tuple.Pair[int, string]
	}{
		{
			name: "positive",
			n:    2,
			want: []tuple.Pair[int, string]{tuple.New(3, "three")},
		},
		{name: "larger than sequence", n: 5},
		{
			name: "zero",
			n:    0,
			want: []tuple.Pair[int, string]{
				tuple.New(1, "one"),
				tuple.New(2, "two"),
				tuple.New(3, "three"),
			},
		},
		{
			name: "negative",
			n:    -1,
			want: []tuple.Pair[int, string]{
				tuple.New(1, "one"),
				tuple.New(2, "two"),
				tuple.New(3, "three"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectSeq2(values.Drop(tt.n)); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Drop(%d) pairs = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestSeq2DropSupportsEarlyTermination(t *testing.T) {
	produced := 0
	values := seq.FromSeq2(func(yield func(int, string) bool) {
		for _, pair := range []tuple.Pair[int, string]{
			tuple.New(1, "one"),
			tuple.New(2, "two"),
			tuple.New(3, "three"),
			tuple.New(4, "four"),
		} {
			produced++
			if !yield(pair.First, pair.Second) {
				return
			}
		}
	})

	for range values.Drop(2).All() {
		break
	}
	if produced != 3 {
		t.Fatalf("breaking Drop(2) iteration consumed %d pairs, want 3", produced)
	}
}

func TestSeq2TakeWhile(t *testing.T) {
	calls := 0
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
		tuple.New(1, "one again"),
	)
	taken := values.TakeWhile(func(key int, _ string) bool {
		calls++
		return key < 3
	})
	if calls != 0 {
		t.Fatalf("TakeWhile predicate called %d times before iteration, want 0", calls)
	}

	got := collectSeq2(taken)
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TakeWhile() pairs = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("TakeWhile predicate called %d times, want 3", calls)
	}
}

func TestSeq2DropWhile(t *testing.T) {
	calls := 0
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
		tuple.New(1, "one again"),
	)
	dropped := values.DropWhile(func(key int, _ string) bool {
		calls++
		return key < 3
	})
	if calls != 0 {
		t.Fatalf("DropWhile predicate called %d times before iteration, want 0", calls)
	}

	got := collectSeq2(dropped)
	want := []tuple.Pair[int, string]{
		tuple.New(3, "three"),
		tuple.New(1, "one again"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DropWhile() pairs = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("DropWhile predicate called %d times, want 3", calls)
	}
}

func TestSeq2Concat(t *testing.T) {
	got := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
	).Concat(
		seq.Empty2[int, string](),
		seq2FromPairs(tuple.New(3, "three")),
	).Collect()
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Concat().Collect() = %v, want %v", got, want)
	}

	withoutOthers := seq2FromPairs(tuple.New(1, "one")).Concat().Collect()
	if want := []tuple.Pair[int, string]{tuple.New(1, "one")}; !reflect.DeepEqual(withoutOthers, want) {
		t.Fatalf("Concat without others = %v, want %v", withoutOthers, want)
	}
}

func TestSeq2ConcatIsLazyAndSupportsEarlyTermination(t *testing.T) {
	leftProduced := 0
	left := seq.FromSeq2(func(yield func(int, string) bool) {
		for _, pair := range []tuple.Pair[int, string]{
			tuple.New(1, "one"),
			tuple.New(2, "two"),
		} {
			leftProduced++
			if !yield(pair.First, pair.Second) {
				return
			}
		}
	})
	rightProduced := 0
	right := seq.FromSeq2(func(yield func(int, string) bool) {
		rightProduced++
		yield(3, "three")
	})
	concatenated := left.Concat(right)
	if leftProduced != 0 || rightProduced != 0 {
		t.Fatal("Concat consumed input before iteration")
	}

	for range concatenated.All() {
		break
	}
	if leftProduced != 1 || rightProduced != 0 {
		t.Fatalf(
			"breaking Concat iteration consumed (%d, %d) pairs, want (1, 0)",
			leftProduced,
			rightProduced,
		)
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

func TestSeq2First(t *testing.T) {
	produced := 0
	values := seq.FromSeq2(func(yield func(int, string) bool) {
		for _, pair := range []tuple.Pair[int, string]{
			tuple.New(1, "one"),
			tuple.New(2, "two"),
		} {
			produced++
			if !yield(pair.First, pair.Second) {
				return
			}
		}
	})

	if key, value, ok := values.First(); !ok || key != 1 || value != "one" {
		t.Fatalf("First() = (%d, %q, %t), want (1, \"one\", true)", key, value, ok)
	}
	if produced != 1 {
		t.Fatalf("First consumed %d pairs, want 1", produced)
	}
	if key, value, ok := seq.Empty2[int, string]().First(); ok || key != 0 || value != "" {
		t.Fatalf("empty First() = (%d, %q, %t), want (0, \"\", false)", key, value, ok)
	}
}

func TestSeq2Last(t *testing.T) {
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	)
	if key, value, ok := values.Last(); !ok || key != 3 || value != "three" {
		t.Fatalf("Last() = (%d, %q, %t), want (3, \"three\", true)", key, value, ok)
	}
	if key, value, ok := seq.Empty2[int, string]().Last(); ok || key != 0 || value != "" {
		t.Fatalf("empty Last() = (%d, %q, %t), want (0, \"\", false)", key, value, ok)
	}
}

func TestSeq2Find(t *testing.T) {
	calls := 0
	values := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	)
	key, value, ok := values.Find(func(key int, value string) bool {
		calls++
		return key == 2 && value == "two"
	})
	if !ok || key != 2 || value != "two" {
		t.Fatalf("Find() = (%d, %q, %t), want (2, \"two\", true)", key, value, ok)
	}
	if calls != 2 {
		t.Fatalf("Find predicate called %d times, want 2", calls)
	}

	key, value, ok = values.Find(func(key int, _ string) bool { return key > 3 })
	if ok || key != 0 || value != "" {
		t.Fatalf("non-matching Find() = (%d, %q, %t), want (0, \"\", false)", key, value, ok)
	}
}

func TestSeq2Any(t *testing.T) {
	calls := 0
	got := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	).Any(func(key int, _ string) bool {
		calls++
		return key == 2
	})
	if !got || calls != 2 {
		t.Fatalf("Any() = %t with %d calls, want true with 2 calls", got, calls)
	}

	if seq.Empty2[int, string]().Any(func(int, string) bool {
		t.Fatal("Any predicate called for empty sequence")
		return true
	}) {
		t.Fatal("empty Any() = true, want false")
	}
}

func TestSeq2Every(t *testing.T) {
	calls := 0
	got := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	).Every(func(key int, _ string) bool {
		calls++
		return key < 2
	})
	if got || calls != 2 {
		t.Fatalf("Every() = %t with %d calls, want false with 2 calls", got, calls)
	}

	if !seq.Empty2[int, string]().Every(func(int, string) bool {
		t.Fatal("Every predicate called for empty sequence")
		return false
	}) {
		t.Fatal("empty Every() = false, want true")
	}
}

func TestSeq2Count(t *testing.T) {
	if got := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(3, "three"),
	).Count(); got != 3 {
		t.Fatalf("Count() = %d, want 3", got)
	}
	if got := seq.Empty2[int, string]().Count(); got != 0 {
		t.Fatalf("empty Count() = %d, want 0", got)
	}
}

func TestSeq2Collect(t *testing.T) {
	got := seq2FromPairs(
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(1, "one again"),
	).Collect()
	want := []tuple.Pair[int, string]{
		tuple.New(1, "one"),
		tuple.New(2, "two"),
		tuple.New(1, "one again"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Collect() = %v, want %v", got, want)
	}
	if got := seq.Empty2[int, string]().Collect(); len(got) != 0 {
		t.Fatalf("empty Collect() = %v, want empty", got)
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
