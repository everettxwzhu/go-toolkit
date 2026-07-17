package seq

import (
	"cmp"
	"slices"
)

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

// Last returns the last value of s and true. If s is empty, Last returns the
// zero value of T and false.
func (s Seq[T]) Last() (T, bool) {
	var last T
	found := false

	for value := range s.All() {
		last = value
		found = true
	}

	return last, found
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

// At returns the value at the zero-based index and true. If index is negative
// or outside the sequence, At returns the zero value of T and false.
func (s Seq[T]) At(index int) (T, bool) {
	if index < 0 {
		var zero T
		return zero, false
	}

	current := 0
	for value := range s.All() {
		if current == index {
			return value, true
		}
		current++
	}

	var zero T
	return zero, false
}

// Single returns the only value of s and true. It returns the zero value of T
// and false when s is empty or contains more than one value.
func (s Seq[T]) Single() (T, bool) {
	var single T
	found := false

	for value := range s.All() {
		if found {
			var zero T
			return zero, false
		}

		single = value
		found = true
	}

	return single, found
}

// ReduceFirst combines the values of s from left to right using the first value
// as the initial accumulator. It returns false when s is empty.
func (s Seq[T]) ReduceFirst(reducer func(T, T) T) (T, bool) {
	var result T
	found := false

	for value := range s.All() {
		if !found {
			result = value
			found = true
			continue
		}

		result = reducer(result, value)
	}

	return result, found
}

// Partition consumes s and separates its values according to predicate.
// Matched and unmatched values retain their original relative order.
func (s Seq[T]) Partition(predicate func(T) bool) (matched, unmatched []T) {
	for value := range s.All() {
		if predicate(value) {
			matched = append(matched, value)
		} else {
			unmatched = append(unmatched, value)
		}
	}

	return matched, unmatched
}

// MinBy returns the value whose key is smallest and true. It returns the zero
// value of T and false when s is empty. When multiple values have the smallest
// key, MinBy returns the first.
func (s Seq[T]) MinBy[K cmp.Ordered](key func(T) K) (T, bool) {
	var minimum T
	var minimumKey K
	found := false

	for value := range s.All() {
		valueKey := key(value)
		if !found || valueKey < minimumKey {
			minimum = value
			minimumKey = valueKey
			found = true
		}
	}

	return minimum, found
}

// MaxBy returns the value whose key is largest and true. It returns the zero
// value of T and false when s is empty. When multiple values have the largest
// key, MaxBy returns the first.
func (s Seq[T]) MaxBy[K cmp.Ordered](key func(T) K) (T, bool) {
	var maximum T
	var maximumKey K
	found := false

	for value := range s.All() {
		valueKey := key(value)
		if !found || valueKey > maximumKey {
			maximum = value
			maximumKey = valueKey
			found = true
		}
	}

	return maximum, found
}

// MinFunc returns the smallest value according to compare and true. It returns
// the zero value of T and false when s is empty.
//
// Compare must return a negative value when its first argument is smaller, zero
// when the arguments are equal, and a positive value otherwise.
func (s Seq[T]) MinFunc(compare func(T, T) int) (T, bool) {
	var minimum T
	found := false

	for value := range s.All() {
		if !found || compare(value, minimum) < 0 {
			minimum = value
			found = true
		}
	}

	return minimum, found
}

// MaxFunc returns the largest value according to compare and true. It returns
// the zero value of T and false when s is empty.
//
// Compare follows the same contract as MinFunc.
func (s Seq[T]) MaxFunc(compare func(T, T) int) (T, bool) {
	var maximum T
	found := false

	for value := range s.All() {
		if !found || compare(value, maximum) > 0 {
			maximum = value
			found = true
		}
	}

	return maximum, found
}
