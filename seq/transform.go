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

// Inspect returns a sequence that calls action for each value immediately
// before yielding it. Inspect does not modify the yielded values.
func (s Seq[T]) Inspect(action func(T)) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for value := range s.All() {
			action(value)
			if !yield(value) {
				return
			}
		}
	})
}

// Scan returns a sequence containing each intermediate accumulated value
// produced by combining initial and the values of s from left to right.
//
// The initial value itself is not yielded.
func (s Seq[T]) Scan[U any](initial U, reducer func(U, T) U) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		accumulator := initial

		for value := range s.All() {
			accumulator = reducer(accumulator, value)
			if !yield(accumulator) {
				return
			}
		}
	})
}

// Enumerate returns a key-value sequence that pairs each value of s with its
// zero-based position.
func (s Seq[T]) Enumerate() Seq2[int, T] {
	return FromSeq2(func(yield func(int, T) bool) {
		index := 0

		for value := range s.All() {
			if !yield(index, value) {
				return
			}
			index++
		}
	})
}

// Chunk returns a sequence of non-overlapping slices containing at most size
// values from s. It returns an empty sequence when size is not positive.
//
// Each yielded slice must own its backing storage and remain valid after the
// sequence advances.
//
// Chunk is a top-level function because defining it as a method would create an
// instantiation cycle from Seq[T] to Seq[[]T].
func Chunk[T any](s Seq[T], size int) Seq[[]T] {
	if size <= 0 {
		return Empty[[]T]()
	}

	return FromSeq(func(yield func([]T) bool) {
		chunk := make([]T, 0, size)

		for value := range s.All() {
			chunk = append(chunk, value)
			if len(chunk) < size {
				continue
			}

			if !yield(chunk) {
				return
			}
			chunk = make([]T, 0, size)
		}

		if len(chunk) > 0 {
			yield(chunk)
		}
	})
}

// Window returns a sequence containing each consecutive sliding window of the
// requested size. It yields no values when size is not positive or s contains
// fewer than size values.
//
// Each yielded slice must own its backing storage and remain valid after the
// sequence advances.
//
// Window is a top-level function because defining it as a method would create
// an instantiation cycle from Seq[T] to Seq[[]T].
func Window[T any](s Seq[T], size int) Seq[[]T] {
	if size <= 0 {
		return Empty[[]T]()
	}

	return FromSeq(func(yield func([]T) bool) {
		window := make([]T, 0, size)

		for value := range s.All() {
			if len(window) < size {
				window = append(window, value)
			} else {
				copy(window, window[1:])
				window[size-1] = value
			}

			if len(window) == size {
				current := append([]T(nil), window...)
				if !yield(current) {
					return
				}
			}
		}
	})
}
