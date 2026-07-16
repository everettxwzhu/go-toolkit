package seq

// Map returns a sequence produced by applying transform to each value of s.
func (s Seq[T]) Map[U any](transform func(T) U) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		for value := range s.All() {
			if !yield(transform(value)) {
				return
			}
		}
	})
}

// FilterMap applies transform to each value of s and returns a sequence of the
// transformed values whose accompanying boolean is true.
func (s Seq[T]) FilterMap[U any](transform func(T) (U, bool)) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		for value := range s.All() {
			transformed, keep := transform(value)
			if !keep {
				continue
			}

			if !yield(transformed) {
				return
			}
		}
	})
}

// FlatMap applies transform to each value of s and returns a sequence formed by
// concatenating the resulting sequences in order.
func (s Seq[T]) FlatMap[U any](transform func(T) Seq[U]) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		for value := range s.All() {
			inner := transform(value)

			for transformed := range inner.All() {
				if !yield(transformed) {
					return
				}
			}
		}
	})
}
