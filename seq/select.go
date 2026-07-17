package seq

// Filter returns a sequence containing only values for which predicate returns
// true. Values retain their original order.
func (s Seq[T]) Filter(predicate func(T) bool) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for value := range s.All() {
			if !predicate(value) {
				continue
			}

			if !yield(value) {
				return
			}
		}
	})
}

// Take returns a sequence containing at most the first n values of s.
// It returns an empty sequence when n is not positive.
func (s Seq[T]) Take(n int) Seq[T] {
	if n <= 0 {
		return Empty[T]()
	}

	return FromSeq(func(yield func(T) bool) {
		count := 0

		for value := range s.All() {
			if !yield(value) {
				return
			}

			count++
			if n == count {
				return
			}
		}
	})
}

// Drop returns a sequence that skips the first n values of s.
// It returns s unchanged when n is not positive.
func (s Seq[T]) Drop(n int) Seq[T] {
	if n <= 0 {
		return s
	}

	return FromSeq(func(yield func(T) bool) {
		dropped := 0

		for value := range s.All() {
			if dropped < n {
				dropped++
				continue
			}

			if !yield(value) {
				return
			}
		}
	})
}

// TakeWhile returns the longest prefix of s for which predicate returns true.
func (s Seq[T]) TakeWhile(predicate func(T) bool) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		for value := range s.All() {
			if !predicate(value) {
				return
			}

			if !yield(value) {
				return
			}
		}
	})
}

// DropWhile returns s without its longest prefix for which predicate returns
// true. After the first false result, predicate is not called again.
func (s Seq[T]) DropWhile(predicate func(T) bool) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		dropping := true

		for value := range s.All() {
			if dropping && predicate(value) {
				continue
			}

			dropping = false

			if !yield(value) {
				return
			}
		}
	})
}

// DistinctBy returns a sequence containing the first value for each distinct
// key produced by key. Values retain their original order.
func (s Seq[T]) DistinctBy[K comparable](key func(T) K) Seq[T] {
	return FromSeq(func(yield func(T) bool) {
		seen := make(map[K]struct{})

		for value := range s.All() {
			k := key(value)

			if _, exists := seen[k]; exists {
				continue
			}

			seen[k] = struct{}{}

			if !yield(value) {
				return
			}
		}
	})
}

// StepBy returns a sequence containing the first value of s and every step-th
// value after it. It returns an empty sequence when step is not positive.
func (s Seq[T]) StepBy(step int) Seq[T] {
	if step <= 0 {
		return Empty[T]()
	}

	return FromSeq(func(yield func(T) bool) {
		i := 0
		for v := range s.All() {
			if i%step == 0 && !yield(v) {
				return
			}
			i++
		}
	})
}
