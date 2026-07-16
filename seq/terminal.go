package seq

import "slices"

// Collect consumes s and returns its values in iteration order as a slice.
func (s Seq[T]) Collect() []T {
	return slices.Collect(s.All())
}

// Reduce consumes s by combining initial and each value with reducer from left
// to right, then returns the accumulated result.
func (s Seq[T]) Reduce[U any](initial U, reducer func(U, T) U) U {
	result := initial

	for value := range s.All() {
		result = reducer(result, value)
	}

	return result
}

// First returns the first value of s and true. If s is empty, First returns the
// zero value of T and false.
func (s Seq[T]) First() (T, bool) {
	for value := range s.All() {
		return value, true
	}

	var zero T
	return zero, false
}

// Find returns the first value of s for which predicate returns true. If no
// value matches, Find returns the zero value of T and false.
func (s Seq[T]) Find(predicate func(T) bool) (T, bool) {
	for value := range s.All() {
		if predicate(value) {
			return value, true
		}
	}

	var zero T
	return zero, false
}

// Count consumes s and returns the number of values it produces.
func (s Seq[T]) Count() int {
	count := 0

	for range s.All() {
		count++
	}

	return count
}

// Any reports whether predicate returns true for at least one value of s.
// It stops iteration at the first match.
func (s Seq[T]) Any(predicate func(T) bool) bool {
	for value := range s.All() {
		if predicate(value) {
			return true
		}
	}

	return false
}

// Every reports whether predicate returns true for every value of s.
// It stops iteration at the first non-matching value and returns true for an
// empty sequence.
func (s Seq[T]) Every(predicate func(T) bool) bool {
	for value := range s.All() {
		if !predicate(value) {
			return false
		}
	}

	return true
}

// ForEach calls action once for each value of s in iteration order.
func (s Seq[T]) ForEach(action func(T)) {
	for value := range s.All() {
		action(value)
	}
}
