package hashmap_test

import (
	"hash/maphash"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

func TestMapValues(t *testing.T) {
	original := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	original.Set("one", 1)
	original.Set("two", 2)

	mapped := original.MapValues(func(key string, value int) string {
		return key + "=" + strconv.Itoa(value)
	})

	for key, want := range map[string]string{"one": "one=1", "two": "two=2"} {
		if got, ok := mapped.Get(key); !ok || got != want {
			t.Fatalf("MapValues().Get(%q) = (%q, %t), want (%q, true)", key, got, ok, want)
		}
	}
	if got, _ := original.Get("one"); got != 1 {
		t.Fatalf("MapValues changed original value to %d", got)
	}

	mapped.Set("three", "three=3")
	if original.ContainsKey("three") {
		t.Fatal("mutating MapValues result changed original Map")
	}
}

func TestMapValuesPreservesNonComparableKeys(t *testing.T) {
	original := hashmap.New[[]byte, int](bytesHasher{})
	original.Set([]byte("one"), 1)

	mapped := original.MapValues(func(_ []byte, value int) string {
		return strconv.Itoa(value)
	})

	if got, ok := mapped.Get([]byte("one")); !ok || got != "1" {
		t.Fatalf("MapValues().Get(one) = (%q, %t), want (1, true)", got, ok)
	}
}

func TestMapValuesOfZeroMapCanBeModified(t *testing.T) {
	var original hashmap.Map[string, int, maphash.ComparableHasher[string]]

	mapped := original.MapValues(func(string, int) string {
		t.Fatal("MapValues transform called for empty Map")
		return ""
	})
	mapped.Set("one", "value")

	if got, ok := mapped.Get("one"); !ok || got != "value" {
		t.Fatalf("MapValues().Set/Get = (%q, %t), want (value, true)", got, ok)
	}
}

func TestMapKeys(t *testing.T) {
	original := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	original.Set("one", 1)
	original.Set("two", 2)

	mapped := original.MapKeys(bytesHasher{}, func(key string, _ int) []byte {
		return []byte(key)
	})

	for key, want := range map[string]int{"one": 1, "two": 2} {
		if got, ok := mapped.Get([]byte(key)); !ok || got != want {
			t.Fatalf("MapKeys().Get(%q) = (%d, %t), want (%d, true)", key, got, ok, want)
		}
	}
	if original.Len() != 2 {
		t.Fatalf("MapKeys changed original Len() to %d", original.Len())
	}
}

func TestMapPairs(t *testing.T) {
	original := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	original.Set("one", 1)
	original.Set("two", 2)

	mapped := original.MapPairs(
		bytesHasher{},
		func(key string, value int) ([]byte, string) {
			return []byte(strconv.Itoa(value)), key
		},
	)

	for key, want := range map[string]string{"1": "one", "2": "two"} {
		if got, ok := mapped.Get([]byte(key)); !ok || got != want {
			t.Fatalf("MapPairs().Get(%q) = (%q, %t), want (%q, true)", key, got, ok, want)
		}
	}
}

func TestMapTransformsNilMap(t *testing.T) {
	var original *hashmap.Map[string, int, maphash.ComparableHasher[string]]

	if got := original.MapValues(func(string, int) string {
		t.Fatal("MapValues transform called for nil Map")
		return ""
	}); got != nil {
		t.Fatal("nil Map.MapValues() is non-nil")
	}

	if got := original.MapKeys(bytesHasher{}, func(string, int) []byte {
		t.Fatal("MapKeys transform called for nil Map")
		return nil
	}); got != nil {
		t.Fatal("nil Map.MapKeys() is non-nil")
	}

	if got := original.MapPairs(bytesHasher{}, func(string, int) ([]byte, string) {
		t.Fatal("MapPairs transform called for nil Map")
		return nil, ""
	}); got != nil {
		t.Fatal("nil Map.MapPairs() is non-nil")
	}
}
