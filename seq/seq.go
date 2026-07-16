package seq

import "iter"

// Seq is a lazy sequence of values produced by an [iter.Seq].
// Its zero value is an empty sequence.
type Seq[T any] struct {
	source iter.Seq[T]
}

// All returns the iterator that produces the values in s.
// For the zero value of Seq, All returns an empty iterator.
func (s Seq[T]) All() iter.Seq[T] {
	if s.source == nil {
		return func(yield func(T) bool) {}
	}

	return s.source
}
