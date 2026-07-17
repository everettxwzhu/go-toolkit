package seq

// GroupBy groups the values of s by the key returned by key.
// Values within each group retain their original order.
func (s Seq[T]) GroupBy[K comparable](key func(T) K) map[K][]T {
	result := make(map[K][]T)

	for value := range s.All() {
		groupedKey := key(value)
		result[groupedKey] = append(result[groupedKey], value)
	}

	return result
}

// ToMap transforms each value of s into a key-value pair and returns a map of
// those pairs. When multiple values produce the same key, the last value wins.
func (s Seq[T]) ToMap[K comparable, V any](transform func(T) (K, V)) map[K]V {
	result := make(map[K]V)

	for value := range s.All() {
		k, v := transform(value)
		result[k] = v
	}

	return result
}

// CountBy groups the values of s by key and returns the number of values in
// each group.
func (s Seq[T]) CountBy[K comparable](key func(T) K) map[K]int {
	result := make(map[K]int)

	for value := range s.All() {
		k := key(value)
		result[k]++
	}

	return result
}

// IndexBy returns a map from each key to the last value of s that produced that
// key.
func (s Seq[T]) IndexBy[K comparable](key func(T) K) map[K]T {
	result := make(map[K]T)

	for value := range s.All() {
		k := key(value)
		result[k] = value
	}

	return result
}
