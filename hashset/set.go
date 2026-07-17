package hashset

import (
	"hash/maphash"
	"iter"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

type empty = struct{}

// Set is a mutable collection of distinct values. H defines the hash function
// and equality relation used for values, so K need not satisfy comparable.
//
// Copying a Set shares its underlying storage; use [Set.Clone] when an
// independent copy is required. A Set is not safe for concurrent mutation.
type Set[K any, H maphash.Hasher[K]] struct {
	hasher H
	values *hashmap.Map[K, empty, H]
}

func makeSet[K any, H maphash.Hasher[K]](
	hasher H,
	capacity int,
) Set[K, H] {
	return Set[K, H]{
		hasher: hasher,
		values: hashmap.NewWithCapacity[K, empty](hasher, capacity),
	}
}

// New returns an initialized Set containing values.
// Duplicate values, as defined by hasher, are stored only once.
func New[K any, H maphash.Hasher[K]](hasher H, values ...K) Set[K, H] {
	s := makeSet(hasher, len(values))
	s.Add(values...)
	return s
}

// NewWithCapacity returns an empty Set with enough initial storage for
// approximately capacity values. It panics when capacity is negative.
func NewWithCapacity[K any, H maphash.Hasher[K]](
	hasher H,
	capacity int,
) Set[K, H] {
	return makeSet(hasher, capacity)
}

// Add inserts values into s. Adding an existing value has no effect.
//
// The zero value of Set can be modified when the zero value of H is a valid
// Hasher.
func (s *Set[K, H]) Add(values ...K) {
	s.initialize()
	for _, value := range values {
		s.values.Set(value, empty{})
	}
}

// Remove deletes values from s. Removing a missing value has no effect.
func (s *Set[K, H]) Remove(values ...K) {
	if s.values == nil {
		return
	}

	for _, value := range values {
		s.values.Delete(value)
	}
}

// Contains reports whether s contains value.
func (s Set[K, H]) Contains(value K) bool {
	return s.values.ContainsKey(value)
}

// Len returns the number of values in s.
func (s Set[K, H]) Len() int {
	return s.values.Len()
}

// Clone returns an independent shallow copy of s. Values themselves are not
// cloned.
func (s Set[K, H]) Clone() Set[K, H] {
	return Set[K, H]{
		hasher: s.hasher,
		values: s.values.Clone(),
	}
}

// All returns an iterator over the values in s in unspecified order.
//
// The Set must not be mutated while the iterator is running.
func (s Set[K, H]) All() iter.Seq[K] {
	return s.values.Keys()
}

// Filter returns a new Set containing the values for which predicate returns
// true. It does not modify s.
func (s Set[K, H]) Filter(predicate func(K) bool) Set[K, H] {
	result := makeSet(s.hasher, s.Len())
	for value := range s.All() {
		if predicate(value) {
			result.Add(value)
		}
	}
	return result
}

// Map returns a new Set containing the transformed values from s. Hasher
// defines hashing and equality for the transformed values. Multiple input
// values may produce equal output values. Map does not modify s.
func (s Set[K, H]) Map[
	L any,
	G maphash.Hasher[L],
](
	hasher G,
	transform func(K) L,
) Set[L, G] {
	result := makeSet(hasher, s.Len())
	for value := range s.All() {
		result.Add(transform(value))
	}
	return result
}

// Union returns a new Set containing values present in s or other.
// It does not modify either input Set.
func (s Set[K, H]) Union(other Set[K, H]) Set[K, H] {
	result := makeSet(s.hasher, s.Len()+other.Len())
	for value := range s.All() {
		result.Add(value)
	}
	for value := range other.All() {
		result.Add(value)
	}
	return result
}

// Intersection returns a new Set containing values present in both s and
// other. It does not modify either input Set.
func (s Set[K, H]) Intersection(other Set[K, H]) Set[K, H] {
	small, large := s, other
	if small.Len() > large.Len() {
		small, large = large, small
	}

	result := makeSet(s.hasher, small.Len())
	for value := range small.All() {
		if large.Contains(value) {
			result.Add(value)
		}
	}
	return result
}

// Difference returns a new Set containing values present in s but not other.
// It does not modify either input Set.
func (s Set[K, H]) Difference(other Set[K, H]) Set[K, H] {
	result := makeSet(s.hasher, s.Len())
	for value := range s.All() {
		if !other.Contains(value) {
			result.Add(value)
		}
	}
	return result
}

// SymmetricDifference returns a new Set containing values present in exactly
// one of s and other. It does not modify either input Set.
func (s Set[K, H]) SymmetricDifference(other Set[K, H]) Set[K, H] {
	result := makeSet(s.hasher, s.Len()+other.Len())
	for value := range s.All() {
		if !other.Contains(value) {
			result.Add(value)
		}
	}
	for value := range other.All() {
		if !s.Contains(value) {
			result.Add(value)
		}
	}
	return result
}

// IsSubset reports whether every value in s is also present in other.
func (s Set[K, H]) IsSubset(other Set[K, H]) bool {
	if s.Len() > other.Len() {
		return false
	}

	for value := range s.All() {
		if !other.Contains(value) {
			return false
		}
	}
	return true
}

// Equal reports whether s and other contain the same values.
func (s Set[K, H]) Equal(other Set[K, H]) bool {
	return s.Len() == other.Len() && s.IsSubset(other)
}

func (s *Set[K, H]) initialize() {
	if s.values == nil {
		s.values = hashmap.New[K, empty](s.hasher)
	}
}
