package seq

import "iter"

// Empty returns a sequence that contains no values.
func Empty[T any]() Seq[T] {
	return Seq[T]{source: func(yield func(T) bool) {}}
}

// FromSeq wraps source as a Seq.
// A nil source is treated as an empty sequence when iterated.
func FromSeq[T any](source iter.Seq[T]) Seq[T] {
	return Seq[T]{source: source}
}

// FromSlice returns a sequence that yields the values in values in order.
// The slice is read during iteration and is not copied.
func FromSlice[T any](values []T) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	})
}

// Of returns a sequence that yields values in argument order.
func Of[T any](values ...T) Seq[T] {
	return FromSlice(values)
}

// Range returns a sequence of consecutive integers from start, inclusive, to
// end, exclusive. It is empty when start is greater than or equal to end.
func Range(start, end int) Seq[int] {
	return FromSeq(func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	})
}
