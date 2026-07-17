package hashmap_test

import (
	"bytes"
	"hash/maphash"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

type bytesHasher struct{}

func (bytesHasher) Hash(hash *maphash.Hash, key []byte) {
	hash.Write(key)
}

func (bytesHasher) Equal(left, right []byte) bool {
	return bytes.Equal(left, right)
}

type constantBytesHasher struct{}

func (constantBytesHasher) Hash(*maphash.Hash, []byte) {}

func (constantBytesHasher) Equal(left, right []byte) bool {
	return bytes.Equal(left, right)
}

type caseInsensitiveHasher struct{}

func (caseInsensitiveHasher) Hash(hash *maphash.Hash, key string) {
	hash.WriteString(strings.ToLower(key))
}

func (caseInsensitiveHasher) Equal(left, right string) bool {
	return strings.ToLower(left) == strings.ToLower(right)
}

func TestSetGetAndContainsKey(t *testing.T) {
	m := hashmap.New[string, int](maphash.ComparableHasher[string]{})

	if m.Len() != 0 {
		t.Fatalf("new Map Len() = %d, want 0", m.Len())
	}
	if value, ok := m.Get("missing"); ok || value != 0 {
		t.Fatalf("Get(missing) = (%d, %t), want (0, false)", value, ok)
	}

	if previous, replaced := m.Set("one", 1); replaced || previous != 0 {
		t.Fatalf("first Set() = (%d, %t), want (0, false)", previous, replaced)
	}
	if value, ok := m.Get("one"); !ok || value != 1 {
		t.Fatalf("Get(one) = (%d, %t), want (1, true)", value, ok)
	}
	if !m.ContainsKey("one") {
		t.Fatal("ContainsKey(one) = false, want true")
	}

	if previous, replaced := m.Set("one", 11); !replaced || previous != 1 {
		t.Fatalf("replacing Set() = (%d, %t), want (1, true)", previous, replaced)
	}
	if value, ok := m.Get("one"); !ok || value != 11 {
		t.Fatalf("Get(one) after replacement = (%d, %t), want (11, true)", value, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("Len() after replacement = %d, want 1", m.Len())
	}
}

func TestNonComparableKeys(t *testing.T) {
	m := hashmap.New[[]byte, string](bytesHasher{})

	firstKey := []byte("key")
	m.Set(firstKey, "first")

	equivalentKey := append([]byte(nil), firstKey...)
	if value, ok := m.Get(equivalentKey); !ok || value != "first" {
		t.Fatalf("Get(equivalent []byte) = (%q, %t), want (\"first\", true)", value, ok)
	}

	if previous, replaced := m.Set(equivalentKey, "second"); !replaced || previous != "first" {
		t.Fatalf("Set(equivalent []byte) = (%q, %t), want (\"first\", true)", previous, replaced)
	}
	if m.Len() != 1 {
		t.Fatalf("Len() after equivalent key = %d, want 1", m.Len())
	}
}

func TestCustomEqualityRetainsStoredKey(t *testing.T) {
	m := hashmap.New[string, int](caseInsensitiveHasher{})
	m.Set("Name", 1)
	m.Set("NAME", 2)

	if value, ok := m.Get("name"); !ok || value != 2 {
		t.Fatalf("Get(name) = (%d, %t), want (2, true)", value, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", m.Len())
	}

	var storedKey string
	for key := range m.All() {
		storedKey = key
	}
	if storedKey != "Name" {
		t.Fatalf("stored key = %q, want original key %q", storedKey, "Name")
	}
}

func TestHashCollisions(t *testing.T) {
	m := hashmap.New[[]byte, int](constantBytesHasher{})

	for i := range 100 {
		m.Set([]byte{byte(i)}, i)
	}
	if m.Len() != 100 {
		t.Fatalf("Len() = %d, want 100", m.Len())
	}

	for i := range 100 {
		if value, ok := m.Get([]byte{byte(i)}); !ok || value != i {
			t.Fatalf("Get(%d) = (%d, %t), want (%d, true)", i, value, ok, i)
		}
	}
}

func TestResizePreservesEntries(t *testing.T) {
	m := hashmap.NewWithCapacity[int, int](maphash.ComparableHasher[int]{}, 1)

	for i := range 10_000 {
		m.Set(i, i*i)
	}
	if m.Len() != 10_000 {
		t.Fatalf("Len() = %d, want 10000", m.Len())
	}

	for i := range 10_000 {
		if value, ok := m.Get(i); !ok || value != i*i {
			t.Fatalf("Get(%d) = (%d, %t), want (%d, true)", i, value, ok, i*i)
		}
	}
}

func TestDelete(t *testing.T) {
	m := hashmap.New[[]byte, int](constantBytesHasher{})
	for i, key := range []string{"one", "two", "three"} {
		m.Set([]byte(key), i+1)
	}

	if value, deleted := m.Delete([]byte("two")); !deleted || value != 2 {
		t.Fatalf("Delete(two) = (%d, %t), want (2, true)", value, deleted)
	}
	if m.ContainsKey([]byte("two")) {
		t.Fatal("ContainsKey(two) after Delete = true, want false")
	}
	if m.Len() != 2 {
		t.Fatalf("Len() after Delete = %d, want 2", m.Len())
	}
	if value, deleted := m.Delete([]byte("missing")); deleted || value != 0 {
		t.Fatalf("Delete(missing) = (%d, %t), want (0, false)", value, deleted)
	}

	for key, want := range map[string]int{"one": 1, "three": 3} {
		if value, ok := m.Get([]byte(key)); !ok || value != want {
			t.Fatalf("Get(%s) = (%d, %t), want (%d, true)", key, value, ok, want)
		}
	}
}

func TestClear(t *testing.T) {
	m := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	m.Set("one", 1)
	m.Set("two", 2)
	m.Clear()

	if m.Len() != 0 {
		t.Fatalf("Len() after Clear = %d, want 0", m.Len())
	}
	if _, ok := m.Get("one"); ok {
		t.Fatal("Get(one) after Clear found an entry")
	}

	m.Set("three", 3)
	if value, ok := m.Get("three"); !ok || value != 3 {
		t.Fatalf("Get(three) after reuse = (%d, %t), want (3, true)", value, ok)
	}
}

func TestClone(t *testing.T) {
	original := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	original.Set("one", 1)
	original.Set("two", 2)

	cloned := original.Clone()
	cloned.Set("three", 3)
	cloned.Delete("one")

	if original.Len() != 2 || !original.ContainsKey("one") || original.ContainsKey("three") {
		t.Fatal("mutating Clone changed original Map")
	}
	if cloned.Len() != 2 || cloned.ContainsKey("one") || !cloned.ContainsKey("three") {
		t.Fatal("Clone does not contain expected independent entries")
	}
}

func TestAllKeysAndValues(t *testing.T) {
	m := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	wantEntries := map[string]int{"one": 1, "two": 2, "three": 3}
	for key, value := range wantEntries {
		m.Set(key, value)
	}

	gotEntries := make(map[string]int)
	for key, value := range m.All() {
		gotEntries[key] = value
	}
	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Fatalf("All() entries = %v, want %v", gotEntries, wantEntries)
	}

	gotKeys := make([]string, 0, m.Len())
	for key := range m.Keys() {
		gotKeys = append(gotKeys, key)
	}
	slices.Sort(gotKeys)
	if want := []string{"one", "three", "two"}; !reflect.DeepEqual(gotKeys, want) {
		t.Fatalf("Keys() = %v, want %v", gotKeys, want)
	}

	gotValues := make([]int, 0, m.Len())
	for value := range m.Values() {
		gotValues = append(gotValues, value)
	}
	slices.Sort(gotValues)
	if want := []int{1, 2, 3}; !reflect.DeepEqual(gotValues, want) {
		t.Fatalf("Values() = %v, want %v", gotValues, want)
	}
}

func TestAllSupportsEarlyTermination(t *testing.T) {
	m := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	calls := 0
	m.All()(func(string, int) bool {
		calls++
		return false
	})
	if calls != 1 {
		t.Fatalf("All yield called %d times, want 1", calls)
	}
}

func TestNilMapReads(t *testing.T) {
	var m *hashmap.Map[string, int, maphash.ComparableHasher[string]]

	if m.Len() != 0 {
		t.Fatalf("nil Map Len() = %d, want 0", m.Len())
	}
	if value, ok := m.Get("key"); ok || value != 0 {
		t.Fatalf("nil Map Get() = (%d, %t), want (0, false)", value, ok)
	}
	if m.ContainsKey("key") {
		t.Fatal("nil Map ContainsKey() = true, want false")
	}
	if value, ok := m.Delete("key"); ok || value != 0 {
		t.Fatalf("nil Map Delete() = (%d, %t), want (0, false)", value, ok)
	}
	if m.Clone() != nil {
		t.Fatal("nil Map Clone() is non-nil")
	}
	m.Clear()
	for range m.All() {
		t.Fatal("nil Map All() produced an entry")
	}
}

func TestSetPanicsOnNilMap(t *testing.T) {
	var m *hashmap.Map[string, int, maphash.ComparableHasher[string]]

	defer func() {
		if recover() == nil {
			t.Fatal("Set on nil Map did not panic")
		}
	}()
	m.Set("key", 1)
}

func TestZeroValueWithZeroValueHasher(t *testing.T) {
	var m hashmap.Map[string, int, maphash.ComparableHasher[string]]
	m.Set("key", 1)

	if value, ok := m.Get("key"); !ok || value != 1 {
		t.Fatalf("zero Map Get(key) = (%d, %t), want (1, true)", value, ok)
	}
}

func TestConcurrentReads(t *testing.T) {
	m := hashmap.NewWithCapacity[int, int](maphash.ComparableHasher[int]{}, 1_000)
	for i := range 1_000 {
		m.Set(i, i*i)
	}

	var wait sync.WaitGroup
	for range 8 {
		wait.Go(func() {
			for i := range 10_000 {
				key := i % 1_000
				if value, ok := m.Get(key); !ok || value != key*key {
					t.Errorf("Get(%d) = (%d, %t), want (%d, true)", key, value, ok, key*key)
					return
				}
			}
		})
	}
	wait.Wait()
}

func TestNewWithCapacityPanicsOnNegativeCapacity(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewWithCapacity(-1) did not panic")
		}
	}()

	hashmap.NewWithCapacity[string, int](maphash.ComparableHasher[string]{}, -1)
}

func FuzzMapAgainstBuiltin(f *testing.F) {
	f.Add([]byte("initial corpus"))
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	f.Fuzz(func(t *testing.T, operations []byte) {
		got := hashmap.New[byte, int](maphash.ComparableHasher[byte]{})
		want := make(map[byte]int)

		for index, operation := range operations {
			key := operation & 31
			switch operation % 3 {
			case 0:
				previous, replaced := got.Set(key, index)
				wantPrevious, wantReplaced := want[key]
				if previous != wantPrevious || replaced != wantReplaced {
					t.Fatalf(
						"Set(%d) = (%d, %t), want (%d, %t)",
						key,
						previous,
						replaced,
						wantPrevious,
						wantReplaced,
					)
				}
				want[key] = index
			case 1:
				value, deleted := got.Delete(key)
				wantValue, wantDeleted := want[key]
				if value != wantValue || deleted != wantDeleted {
					t.Fatalf(
						"Delete(%d) = (%d, %t), want (%d, %t)",
						key,
						value,
						deleted,
						wantValue,
						wantDeleted,
					)
				}
				delete(want, key)
			case 2:
				value, ok := got.Get(key)
				wantValue, wantOK := want[key]
				if value != wantValue || ok != wantOK {
					t.Fatalf(
						"Get(%d) = (%d, %t), want (%d, %t)",
						key,
						value,
						ok,
						wantValue,
						wantOK,
					)
				}
			}

			if got.Len() != len(want) {
				t.Fatalf("Len() = %d, want %d", got.Len(), len(want))
			}
		}
	})
}
