package seq

import (
	"iter"
	"maps"

	"github.com/everettxwzhu/go-toolkit/tuple"
)

// Seq2 is a lazy sequence of key-value pairs produced by an
// iter.Seq2. Its zero value is an empty sequence.
type Seq2[K, V any] struct {
	source iter.Seq2[K, V]
}

// Empty2 returns a key-value sequence that contains no pairs.
func Empty2[K, V any]() Seq2[K, V] {
	return Seq2[K, V]{source: func(yield func(K, V) bool) {}}
}

// FromSeq2 wraps source as a Seq2. A nil source is treated as an empty sequence
// when iterated.
func FromSeq2[K, V any](source iter.Seq2[K, V]) Seq2[K, V] {
	if source == nil {
		return Empty2[K, V]()
	}

	return Seq2[K, V]{source: source}
}

// FromMap returns a key-value sequence containing the entries of values. The
// iteration order is unspecified.
func FromMap[K comparable, V any](values map[K]V) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		for k, v := range values {
			if !yield(k, v) {
				return
			}
		}
	})
}

// All returns the iterator that produces the key-value pairs in s.
func (s Seq2[K, V]) All() iter.Seq2[K, V] {
	if s.source == nil {
		return func(yield func(K, V) bool) {}
	}

	return s.source
}

// Map applies transform to each key-value pair and returns the resulting value
// sequence.
func (s Seq2[K, V]) Map[U any](transform func(K, V) U) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		for k, v := range s.All() {
			if !yield(transform(k, v)) {
				return
			}
		}
	})
}

// MapPairs applies transform to each key-value pair and returns a sequence
// containing the transformed pairs.
func (s Seq2[K, V]) MapPairs[L, U any](transform func(K, V) (L, U)) Seq2[L, U] {
	return FromSeq2(func(yield func(L, U) bool) {
		for k, v := range s.All() {
			transformedKey, transformedValue := transform(k, v)
			if !yield(transformedKey, transformedValue) {
				return
			}
		}
	})
}

// MapKeys applies transform to each key-value pair and returns a sequence
// containing the transformed keys and original values.
func (s Seq2[K, V]) MapKeys[L any](transform func(K, V) L) Seq2[L, V] {
	return FromSeq2(func(yield func(L, V) bool) {
		for k, v := range s.All() {
			if !yield(transform(k, v), v) {
				return
			}
		}
	})
}

// MapValues applies transform to each key-value pair and returns a sequence
// containing the original keys and transformed values.
func (s Seq2[K, V]) MapValues[U any](transform func(K, V) U) Seq2[K, U] {
	return FromSeq2(func(yield func(K, U) bool) {
		for k, v := range s.All() {
			if !yield(k, transform(k, v)) {
				return
			}
		}
	})
}

// Filter returns a key-value sequence containing only pairs for which predicate
// returns true.
func (s Seq2[K, V]) Filter(predicate func(K, V) bool) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		for k, v := range s.All() {
			if !predicate(k, v) {
				continue
			}

			if !yield(k, v) {
				return
			}
		}
	})
}

// Take returns a key-value sequence containing at most the first n pairs of s.
// It returns an empty sequence when n is not positive.
func (s Seq2[K, V]) Take(n int) Seq2[K, V] {
	if n <= 0 {
		return Empty2[K, V]()
	}

	return FromSeq2(func(yield func(K, V) bool) {
		count := 0

		for k, v := range s.All() {
			if !yield(k, v) {
				return
			}

			count++
			if count == n {
				return
			}
		}
	})
}

// Drop returns a key-value sequence that skips the first n pairs of s. It
// returns s unchanged when n is not positive.
func (s Seq2[K, V]) Drop(n int) Seq2[K, V] {
	if n <= 0 {
		return s
	}

	return FromSeq2(func(yield func(K, V) bool) {
		dropped := 0

		for k, v := range s.All() {
			if dropped < n {
				dropped++
				continue
			}

			if !yield(k, v) {
				return
			}
		}
	})
}

// TakeWhile returns the longest prefix of s for which predicate returns true.
func (s Seq2[K, V]) TakeWhile(predicate func(K, V) bool) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		for k, v := range s.All() {
			if !predicate(k, v) {
				return
			}

			if !yield(k, v) {
				return
			}
		}
	})
}

// DropWhile returns s without its longest prefix for which predicate returns
// true. After the first false result, predicate is not called again.
func (s Seq2[K, V]) DropWhile(predicate func(K, V) bool) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		dropping := true

		for k, v := range s.All() {
			if dropping && predicate(k, v) {
				continue
			}

			dropping = false
			if !yield(k, v) {
				return
			}
		}
	})
}

// Inspect calls action for each key-value pair immediately before yielding it
// and returns the pairs unchanged.
func (s Seq2[K, V]) Inspect(action func(K, V)) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		for k, v := range s.All() {
			action(k, v)

			if !yield(k, v) {
				return
			}
		}
	})
}

// Keys returns a sequence containing the keys produced by s.
func (s Seq2[K, V]) Keys() Seq[K] {
	return FromSeq(func(yield func(K) bool) {
		for k := range s.All() {
			if !yield(k) {
				return
			}
		}
	})
}

// Values returns a sequence containing the values produced by s.
func (s Seq2[K, V]) Values() Seq[V] {
	return FromSeq(func(yield func(V) bool) {
		for _, v := range s.All() {
			if !yield(v) {
				return
			}
		}
	})
}

// Swap returns a key-value sequence with each key and value exchanged.
func (s Seq2[K, V]) Swap() Seq2[V, K] {
	return FromSeq2(func(yield func(V, K) bool) {
		for k, v := range s.All() {
			if !yield(v, k) {
				return
			}
		}
	})
}

// Concat returns a key-value sequence that yields all pairs of s followed by
// all pairs of each sequence in others, in argument order.
func (s Seq2[K, V]) Concat(others ...Seq2[K, V]) Seq2[K, V] {
	return FromSeq2(func(yield func(K, V) bool) {
		for k, v := range s.All() {
			if !yield(k, v) {
				return
			}
		}

		for _, other := range others {
			for k, v := range other.All() {
				if !yield(k, v) {
					return
				}
			}
		}
	})
}

// First returns the first pair of s and true. If s is empty, First returns the
// zero values of K and V and false.
func (s Seq2[K, V]) First() (K, V, bool) {
	for k, v := range s.All() {
		return k, v, true
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Last returns the last pair of s and true. If s is empty, Last returns the
// zero values of K and V and false.
func (s Seq2[K, V]) Last() (K, V, bool) {
	var lastK K
	var lastV V
	found := false

	for k, v := range s.All() {
		lastK = k
		lastV = v
		found = true
	}

	return lastK, lastV, found
}

// Find returns the first pair of s for which predicate returns true. If no pair
// matches, Find returns the zero values of K and V and false.
func (s Seq2[K, V]) Find(predicate func(K, V) bool) (K, V, bool) {
	for k, v := range s.All() {
		if predicate(k, v) {
			return k, v, true
		}
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Any reports whether predicate returns true for at least one pair of s. It
// stops iteration at the first match.
func (s Seq2[K, V]) Any(predicate func(K, V) bool) bool {
	for k, v := range s.All() {
		if predicate(k, v) {
			return true
		}
	}

	return false
}

// Every reports whether predicate returns true for every pair of s. It stops
// iteration at the first non-matching pair and returns true for an empty
// sequence.
func (s Seq2[K, V]) Every(predicate func(K, V) bool) bool {
	for k, v := range s.All() {
		if !predicate(k, v) {
			return false
		}
	}

	return true
}

// Count consumes s and returns the number of pairs it produces.
func (s Seq2[K, V]) Count() int {
	count := 0

	for range s.All() {
		count++
	}

	return count
}

// ForEach calls action once for each key-value pair in iteration order.
func (s Seq2[K, V]) ForEach(action func(K, V)) {
	for k, v := range s.All() {
		action(k, v)
	}
}

// Collect consumes s and returns its pairs in iteration order as a slice.
func (s Seq2[K, V]) Collect() []tuple.Pair[K, V] {
	var result []tuple.Pair[K, V]

	for k, v := range s.All() {
		result = append(result, tuple.New(k, v))
	}

	return result
}

// CollectMap consumes s and returns a map containing its pairs. When multiple
// pairs contain the same key, the last value wins.
//
// CollectMap is a top-level function because Seq2 permits non-comparable keys.
func CollectMap[K comparable, V any](s Seq2[K, V]) map[K]V {
	return maps.Collect(s.All())
}
