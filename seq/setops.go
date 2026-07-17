package seq

import (
	"hash/maphash"
	"iter"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

// UnionBy returns the first value for each key encountered in s followed by
// the first new value for each key encountered in other.
func (s Seq[T]) UnionBy[K comparable](other Seq[T], key func(T) K) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		seen := make(map[K]struct{})

		emit := func(values Seq[T]) bool {
			for value := range values.All() {
				valueKey := key(value)
				if _, ok := seen[valueKey]; ok {
					continue
				}

				seen[valueKey] = struct{}{}
				if !yield(value) {
					return false
				}
			}

			return true
		}

		if !emit(s) {
			return
		}
		emit(other)
	})
}

// UnionByHasher returns the first value for each key encountered in s followed
// by the first new value for each key encountered in other. Hasher defines
// hashing and equality for keys, so K need not satisfy comparable.
func (s Seq[T]) UnionByHasher[
	K any,
	H maphash.Hasher[K],
](
	other Seq[T],
	key func(T) K,
	hasher H,
) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		seen := hashmap.New[K, struct{}](hasher)

		emit := func(values Seq[T]) bool {
			for value := range values.All() {
				_, exists := seen.Set(key(value), struct{}{})
				if exists {
					continue
				}

				if !yield(value) {
					return false
				}
			}
			return true
		}

		if !emit(s) {
			return
		}
		emit(other)
	})
}

// IntersectBy returns the first value from s for each key that also occurs in
// other. Values retain their order from s.
func (s Seq[T]) IntersectBy[K comparable](other Seq[T], key func(T) K) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		otherKeys := make(map[K]struct{})
		for value := range other.All() {
			otherKeys[key(value)] = struct{}{}
		}

		emitted := make(map[K]struct{})
		for value := range s.All() {
			valueKey := key(value)
			if _, ok := otherKeys[valueKey]; !ok {
				continue
			}
			if _, ok := emitted[valueKey]; ok {
				continue
			}

			emitted[valueKey] = struct{}{}
			if !yield(value) {
				return
			}
		}
	})
}

// IntersectByHasher returns the first value from s for each key that also
// occurs in other. Hasher defines hashing and equality for keys, so K need not
// satisfy comparable. Values retain their order from s.
func (s Seq[T]) IntersectByHasher[
	K any,
	H maphash.Hasher[K],
](
	other Seq[T],
	key func(T) K,
	hasher H,
) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		otherKeys := hashmap.New[K, struct{}](hasher)
		for value := range other.All() {
			otherKeys.Set(key(value), struct{}{})
		}

		emitted := hashmap.New[K, struct{}](hasher)
		for value := range s.All() {
			valueKey := key(value)
			if !otherKeys.ContainsKey(valueKey) {
				continue
			}
			if _, exists := emitted.Set(valueKey, struct{}{}); exists {
				continue
			}

			if !yield(value) {
				return
			}
		}
	})
}

// ExceptBy returns the first value from s for each key that does not occur in
// other. Values retain their order from s.
func (s Seq[T]) ExceptBy[K comparable](other Seq[T], key func(T) K) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		excluded := make(map[K]struct{})
		for value := range other.All() {
			excluded[key(value)] = struct{}{}
		}

		emitted := make(map[K]struct{})
		for value := range s.All() {
			valueKey := key(value)
			if _, ok := excluded[valueKey]; ok {
				continue
			}
			if _, ok := emitted[valueKey]; ok {
				continue
			}

			emitted[valueKey] = struct{}{}
			if !yield(value) {
				return
			}
		}
	})
}

// ExceptByHasher returns the first value from s for each key that does not
// occur in other. Hasher defines hashing and equality for keys, so K need not
// satisfy comparable. Values retain their order from s.
func (s Seq[T]) ExceptByHasher[
	K any,
	H maphash.Hasher[K],
](
	other Seq[T],
	key func(T) K,
	hasher H,
) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		excluded := hashmap.New[K, struct{}](hasher)
		for value := range other.All() {
			excluded.Set(key(value), struct{}{})
		}

		emitted := hashmap.New[K, struct{}](hasher)
		for value := range s.All() {
			valueKey := key(value)
			if excluded.ContainsKey(valueKey) {
				continue
			}
			if _, exists := emitted.Set(valueKey, struct{}{}); exists {
				continue
			}

			if !yield(value) {
				return
			}
		}
	})
}

// ContainsBy reports whether s contains a value equal to target according to
// equal. It stops at the first match.
func (s Seq[T]) ContainsBy(target T, equal func(T, T) bool) bool {
	for value := range s.All() {
		if equal(value, target) {
			return true
		}
	}

	return false
}

// EqualBy reports whether s and other contain equal values in the same order
// according to equal. It stops at the first mismatch.
func (s Seq[T]) EqualBy(other Seq[T], equal func(T, T) bool) bool {
	nextLeft, stopLeft := iter.Pull(s.All())
	defer stopLeft()
	nextRight, stopRight := iter.Pull(other.All())
	defer stopRight()

	for {
		left, leftOK := nextLeft()
		right, rightOK := nextRight()

		if leftOK != rightOK {
			return false
		}
		if !leftOK {
			return true
		}
		if !equal(left, right) {
			return false
		}
	}
}

// Distinct returns a sequence containing the first occurrence of each value.
//
// Distinct is a top-level function because a method cannot add a comparable
// constraint to the existing type parameter T of Seq[T].
func Distinct[T comparable](s Seq[T]) Seq[T] {
	return s.DistinctBy(func(value T) T {
		return value
	})
}

// Contains reports whether s contains target. It stops at the first match.
//
// Contains is a top-level function because a method cannot add a comparable
// constraint to the existing type parameter T of Seq[T].
func Contains[T comparable](s Seq[T], target T) bool {
	return s.ContainsBy(target, func(left, right T) bool {
		return left == right
	})
}

// Equal reports whether left and right contain equal values in the same order.
//
// Equal is a top-level function because a method cannot add a comparable
// constraint to the existing type parameter T of Seq[T].
func Equal[T comparable](left, right Seq[T]) bool {
	return left.EqualBy(right, func(left, right T) bool {
		return left == right
	})
}
