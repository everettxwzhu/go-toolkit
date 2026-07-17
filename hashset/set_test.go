package hashset_test

import (
	"hash/maphash"
	"slices"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/hashset"
)

func assertBytesSet(
	t *testing.T,
	got hashset.Set[[]byte, bytesHasher],
	want ...string,
) {
	t.Helper()

	if got.Len() != len(want) {
		t.Fatalf("Set.Len() = %d, want %d", got.Len(), len(want))
	}
	for _, value := range want {
		if !got.Contains([]byte(value)) {
			t.Fatalf("Set does not contain %q", value)
		}
	}
}

func TestNewAddRemoveAndContains(t *testing.T) {
	s := hashset.New[[]byte](
		bytesHasher{},
		[]byte("one"),
		[]byte("two"),
		[]byte("one"),
	)
	assertBytesSet(t, s, "one", "two")

	s.Add([]byte("three"), []byte("two"))
	assertBytesSet(t, s, "one", "two", "three")

	s.Remove([]byte("two"), []byte("missing"))
	assertBytesSet(t, s, "one", "three")
}

func TestCustomEqualityRetainsFirstValue(t *testing.T) {
	s := hashset.New(caseInsensitiveHasher{}, "Go", "GO", "Rust")

	if s.Len() != 2 || !s.Contains("go") {
		t.Fatalf("case-insensitive Set = Len %d, Contains(go) %t; want 2, true", s.Len(), s.Contains("go"))
	}

	values := slices.Collect(s.All())
	slices.Sort(values)
	if !slices.Equal(values, []string{"Go", "Rust"}) {
		t.Fatalf("stored values = %v, want [Go Rust]", values)
	}
}

func TestZeroValue(t *testing.T) {
	var s hashset.Set[[]byte, bytesHasher]

	s.Add([]byte("value"))
	assertBytesSet(t, s, "value")
	s.Remove([]byte("value"))
	assertBytesSet(t, s)
}

func TestCopiesShareStorage(t *testing.T) {
	original := hashset.New[[]byte](bytesHasher{}, []byte("one"))
	copied := original

	copied.Add([]byte("two"))
	assertBytesSet(t, original, "one", "two")
}

func TestClone(t *testing.T) {
	original := hashset.New[[]byte](bytesHasher{}, []byte("one"), []byte("two"))
	cloned := original.Clone()

	cloned.Add([]byte("three"))
	cloned.Remove([]byte("one"))

	assertBytesSet(t, original, "one", "two")
	assertBytesSet(t, cloned, "two", "three")
}

func TestAllSupportsEarlyTermination(t *testing.T) {
	s := hashset.New[[]byte](bytesHasher{}, []byte("one"), []byte("two"))
	calls := 0

	s.All()(func([]byte) bool {
		calls++
		return false
	})
	if calls != 1 {
		t.Fatalf("All yield called %d times, want 1", calls)
	}
}

func TestFilterAndMap(t *testing.T) {
	original := hashset.New[[]byte](
		bytesHasher{},
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	)

	filtered := original.Filter(func(value []byte) bool {
		return len(value) == 3
	})
	assertBytesSet(t, filtered, "one", "two")

	mapped := original.Map(
		maphash.ComparableHasher[int]{},
		func(value []byte) int { return len(value) },
	)
	if mapped.Len() != 2 || !mapped.Contains(3) || !mapped.Contains(5) {
		t.Fatalf("mapped lengths = %v, want {3, 5}", slices.Collect(mapped.All()))
	}
	assertBytesSet(t, original, "one", "two", "three")
}

func TestSetOperations(t *testing.T) {
	left := hashset.New[[]byte](
		bytesHasher{},
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	)
	right := hashset.New[[]byte](
		bytesHasher{},
		[]byte("two"),
		[]byte("three"),
		[]byte("four"),
	)

	assertBytesSet(t, left.Union(right), "one", "two", "three", "four")
	assertBytesSet(t, left.Intersection(right), "two", "three")
	assertBytesSet(t, left.Difference(right), "one")
	assertBytesSet(t, left.SymmetricDifference(right), "one", "four")

	subset := hashset.New[[]byte](bytesHasher{}, []byte("one"), []byte("two"))
	if !subset.IsSubset(left) {
		t.Fatal("{one, two}.IsSubset(left) = false, want true")
	}
	if left.IsSubset(subset) {
		t.Fatal("left.IsSubset({one, two}) = true, want false")
	}
	if !left.Equal(left.Clone()) {
		t.Fatal("Set is not equal to its clone")
	}
	if left.Equal(right) {
		t.Fatal("different Sets compare equal")
	}
}

func TestNewWithCapacityPanicsOnNegativeCapacity(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewWithCapacity(-1) did not panic")
		}
	}()

	hashset.NewWithCapacity[[]byte](bytesHasher{}, -1)
}

func TestMapChangesElementType(t *testing.T) {
	s := hashset.New(
		maphash.ComparableHasher[int]{},
		1,
		2,
		3,
	)
	got := s.Map(
		maphash.ComparableHasher[string]{},
		strconv.Itoa,
	)

	if got.Len() != 3 {
		t.Fatalf("Map Len() = %d, want 3", got.Len())
	}
	for _, value := range []string{"1", "2", "3"} {
		if !got.Contains(value) {
			t.Fatalf("Map does not contain %q", value)
		}
	}
}
