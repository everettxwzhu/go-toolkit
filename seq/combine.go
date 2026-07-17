package seq

import (
	"iter"

	"github.com/everettxwzhu/go-toolkit/tuple"
)

// Concat returns a sequence that yields all values of s followed by all values
// of each sequence in others, in argument order.
func (s Seq[T]) Concat(others ...Seq[T]) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for value := range s.All() {
			if !yield(value) {
				return
			}
		}

		for _, other := range others {
			for value := range other.All() {
				if !yield(value) {
					return
				}
			}
		}
	})
}

// ZipWith combines corresponding values from s and other. Iteration stops when
// either sequence ends, and combine is called once for each resulting value.
func (s Seq[T]) ZipWith[U, V any](other Seq[U], combine func(T, U) V) Seq[V] {
	return FromSeq(func(yield func(V) bool) {
		next1, stop1 := iter.Pull(s.All())
		defer stop1()
		next2, stop2 := iter.Pull(other.All())
		defer stop2()

		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}

			v2, ok2 := next2()
			if !ok2 {
				return
			}
			if !yield(combine(v1, v2)) {
				return
			}
		}
	})
}

// Zip pairs corresponding values from left and right. Iteration stops when
// either sequence ends.
func Zip[T, U any](left Seq[T], right Seq[U]) Seq[tuple.Pair[T, U]] {
	return FromSeq(func(yield func(tuple.Pair[T, U]) bool) {
		next1, stop1 := iter.Pull(left.All())
		defer stop1()
		next2, stop2 := iter.Pull(right.All())
		defer stop2()

		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}

			v2, ok2 := next2()
			if !ok2 {
				return
			}

			if !yield(tuple.New(v1, v2)) {
				return
			}
		}
	})
}

// Intersperse returns a sequence with separator inserted between consecutive
// values of s. It does not add a separator before the first or after the last
// value.
func (s Seq[T]) Intersperse(separator T) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		first := true

		for v := range s.All() {
			if !first && !yield(separator) {
				return
			}

			first = false
			if !yield(v) {
				return
			}
		}
	})
}

// Prepend returns a sequence that yields values followed by the values of s.
func (s Seq[T]) Prepend(values ...T) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}

		for v := range s.All() {
			if !yield(v) {
				return
			}
		}
	})
}

// Append returns a sequence that yields the values of s followed by values.
func (s Seq[T]) Append(values ...T) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for v := range s.All() {
			if !yield(v) {
				return
			}
		}

		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	})
}
