package seq_test

import (
	"bytes"
	"hash/maphash"
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

type bytesHasher struct{}

func (bytesHasher) Hash(hash *maphash.Hash, value []byte) {
	hash.Write(value)
}

func (bytesHasher) Equal(left, right []byte) bool {
	return bytes.Equal(left, right)
}

type keyedValue struct {
	key   []byte
	value string
}

func TestDistinctByHasher(t *testing.T) {
	values := []keyedValue{
		{key: []byte("one"), value: "first one"},
		{key: []byte("two"), value: "first two"},
		{key: []byte("one"), value: "second one"},
	}

	got := seq.FromSlice(values).DistinctByHasher(
		func(value keyedValue) []byte { return value.key },
		bytesHasher{},
	).Map(func(value keyedValue) string {
		return value.value
	}).Collect()

	want := []string{"first one", "first two"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DistinctByHasher().Collect() = %v, want %v", got, want)
	}
}

func TestDistinctByHasherIsLazyAndSupportsEarlyTermination(t *testing.T) {
	produced := 0
	values := seq.FromSeq(func(yield func([]byte) bool) {
		for _, value := range [][]byte{[]byte("one"), []byte("one"), []byte("two"), []byte("three")} {
			produced++
			if !yield(value) {
				return
			}
		}
	})

	distinct := values.DistinctByHasher(
		func(value []byte) []byte { return value },
		bytesHasher{},
	)
	if produced != 0 {
		t.Fatal("DistinctByHasher consumed input before iteration")
	}
	if got := distinct.Take(2).Collect(); !equalBytesList(got, "one", "two") {
		t.Fatalf("DistinctByHasher().Take(2).Collect() = %q, want [one two]", got)
	}
	if produced != 3 {
		t.Fatalf("DistinctByHasher().Take(2) consumed %d values, want 3", produced)
	}
}

func TestGroupByHasher(t *testing.T) {
	values := []keyedValue{
		{key: []byte("one"), value: "first"},
		{key: []byte("two"), value: "second"},
		{key: []byte("one"), value: "third"},
	}
	got := seq.FromSlice(values).GroupByHasher(
		func(value keyedValue) []byte { return value.key },
		bytesHasher{},
	)

	one, ok := got.Get([]byte("one"))
	if !ok || len(one) != 2 || one[0].value != "first" || one[1].value != "third" {
		t.Fatalf("GroupByHasher().Get(one) = (%v, %t), want first and third", one, ok)
	}
	two, ok := got.Get([]byte("two"))
	if !ok || len(two) != 1 || two[0].value != "second" {
		t.Fatalf("GroupByHasher().Get(two) = (%v, %t), want second", two, ok)
	}
}

func TestToHashMapLastValueWins(t *testing.T) {
	values := []keyedValue{
		{key: []byte("one"), value: "first"},
		{key: []byte("one"), value: "second"},
	}
	got := seq.FromSlice(values).ToHashMap(
		func(value keyedValue) ([]byte, string) {
			return value.key, value.value
		},
		bytesHasher{},
	)

	if value, ok := got.Get([]byte("one")); !ok || value != "second" {
		t.Fatalf("ToHashMap().Get(one) = (%q, %t), want (second, true)", value, ok)
	}
}

func TestCountByHasher(t *testing.T) {
	got := seq.Of(
		[]byte("one"),
		[]byte("two"),
		[]byte("one"),
	).CountByHasher(
		func(value []byte) []byte { return value },
		bytesHasher{},
	)

	for key, want := range map[string]int{"one": 2, "two": 1} {
		if count, ok := got.Get([]byte(key)); !ok || count != want {
			t.Fatalf("CountByHasher().Get(%q) = (%d, %t), want (%d, true)", key, count, ok, want)
		}
	}
}

func TestIndexByHasherLastValueWins(t *testing.T) {
	values := []keyedValue{
		{key: []byte("one"), value: "first"},
		{key: []byte("one"), value: "second"},
	}
	got := seq.FromSlice(values).IndexByHasher(
		func(value keyedValue) []byte { return value.key },
		bytesHasher{},
	)

	if value, ok := got.Get([]byte("one")); !ok || value.value != "second" {
		t.Fatalf("IndexByHasher().Get(one) = (%v, %t), want second", value, ok)
	}
}

func TestHasherSetOperations(t *testing.T) {
	left := seq.Of(
		[]byte("one"),
		[]byte("two"),
		[]byte("one"),
		[]byte("three"),
	)
	right := seq.Of(
		[]byte("two"),
		[]byte("four"),
		[]byte("four"),
	)
	key := func(value []byte) []byte { return value }
	hasher := bytesHasher{}

	if got := left.UnionByHasher(right, key, hasher).Collect(); !equalBytesList(got, "one", "two", "three", "four") {
		t.Fatalf("UnionByHasher().Collect() = %q, want [one two three four]", got)
	}
	if got := left.IntersectByHasher(right, key, hasher).Collect(); !equalBytesList(got, "two") {
		t.Fatalf("IntersectByHasher().Collect() = %q, want [two]", got)
	}
	if got := left.ExceptByHasher(right, key, hasher).Collect(); !equalBytesList(got, "one", "three") {
		t.Fatalf("ExceptByHasher().Collect() = %q, want [one three]", got)
	}
}

func equalBytesList(values [][]byte, want ...string) bool {
	if len(values) != len(want) {
		return false
	}
	for i := range values {
		if string(values[i]) != want[i] {
			return false
		}
	}
	return true
}
