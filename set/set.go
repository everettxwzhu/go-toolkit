// Package set provides a generic collection of distinct comparable values.
package set

import (
	"iter"
	"maps"
)

type empty = struct{}

// Set is a mutable collection of distinct comparable values.
//
// Use [New] to construct a Set before modifying it. Copying a Set shares its
// underlying storage; use [Set.Clone] when an independent copy is required.
// A Set is not safe for concurrent use.
type Set[T comparable] struct {
	values map[T]empty
}

func makeSet[T comparable](capacity int) Set[T] {
	return Set[T]{
		values: make(map[T]empty, capacity),
	}
}

// New returns an initialized Set containing values.
// Duplicate values are stored only once.
func New[T comparable](values ...T) Set[T] {
	s := Set[T]{
		values: make(map[T]empty, len(values)),
	}
	s.Add(values...)
	return s
}

// Add inserts values into s. Adding an existing value has no effect.
// The Set must have been constructed with [New].
func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		s.values[v] = empty{}
	}
}

// Remove deletes values from s. Removing a missing value has no effect.
// The Set must have been constructed with [New].
func (s *Set[T]) Remove(values ...T) {
	for _, v := range values {
		delete(s.values, v)
	}
}

// Contains reports whether s contains value.
func (s Set[T]) Contains(value T) bool {
	_, ok := s.values[value]
	return ok
}

// Len returns the number of values in s.
func (s Set[T]) Len() int {
	return len(s.values)
}

// Clone returns an independent copy of s. Mutating the returned Set does not
// affect s.
func (s Set[T]) Clone() Set[T] {
	return Set[T]{
		values: maps.Clone(s.values),
	}
}

// All returns an iterator over the values in s. The iteration order is
// unspecified. The iterator reads the underlying Set during iteration and must
// not be used concurrently with mutations.
func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range s.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Filter returns a new Set containing the values for which predicate returns
// true. It does not modify s.
func (s Set[T]) Filter(predicate func(T) bool) Set[T] {
	result := makeSet[T](s.Len())

	for v := range s.values {
		if predicate(v) {
			result.Add(v)
		}
	}

	return result
}

// Map returns a new Set containing the transformed values from s. Multiple
// input values may produce the same output value. It does not modify s.
func (s Set[T]) Map[U comparable](transform func(T) U) Set[U] {
	result := makeSet[U](s.Len())

	for v := range s.values {
		result.Add(transform(v))
	}

	return result
}

// Union returns a new Set containing values present in s or other.
// It does not modify either input Set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := makeSet[T](s.Len() + other.Len())

	for v := range s.values {
		result.Add(v)
	}
	for v := range other.values {
		result.Add(v)
	}

	return result
}

// Intersection returns a new Set containing values present in both s and
// other. It does not modify either input Set.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := makeSet[T](min(s.Len(), other.Len()))

	small, large := s, other
	if small.Len() > other.Len() {
		small, large = large, small
	}

	for v := range small.values {
		if large.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// Difference returns a new Set containing values present in s but not other.
// It does not modify either input Set.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := makeSet[T](s.Len())

	for v := range s.values {
		if !other.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// SymmetricDifference returns a new Set containing values present in exactly
// one of s and other. It does not modify either input Set.
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	result := makeSet[T](s.Len() + other.Len())

	for value := range s.values {
		if !other.Contains(value) {
			result.Add(value)
		}
	}
	for value := range other.values {
		if !s.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

// IsSubset reports whether every value in s is also present in other.
func (s Set[T]) IsSubset(other Set[T]) bool {
	if s.Len() > other.Len() {
		return false
	}

	for v := range s.values {
		if !other.Contains(v) {
			return false
		}
	}

	return true
}

// Equal reports whether s and other contain the same values.
func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for v := range s.values {
		if !other.Contains(v) {
			return false
		}
	}

	return true
}
