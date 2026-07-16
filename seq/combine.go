package seq

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
