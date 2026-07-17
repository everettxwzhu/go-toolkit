package seq

import (
	"hash/maphash"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

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

// GroupByHasher groups the values of s by the key returned by key. Hasher
// defines hashing and equality for keys, so K need not satisfy comparable.
// Values within each group retain their original order.
func (s Seq[T]) GroupByHasher[
	K any,
	H maphash.Hasher[K],
](
	key func(T) K,
	hasher H,
) *hashmap.Map[K, []T, H] {
	result := hashmap.New[K, []T](hasher)

	for value := range s.All() {
		groupedKey := key(value)
		group, _ := result.Get(groupedKey)
		result.Set(groupedKey, append(group, value))
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

// ToHashMap transforms each value of s into a key-value pair and returns a
// hash Map of those pairs. Hasher defines hashing and equality for keys, so K
// need not satisfy comparable. When multiple values produce equal keys, the
// last value wins.
func (s Seq[T]) ToHashMap[
	K, V any,
	H maphash.Hasher[K],
](
	transform func(T) (K, V),
	hasher H,
) *hashmap.Map[K, V, H] {
	result := hashmap.New[K, V](hasher)

	for value := range s.All() {
		key, mappedValue := transform(value)
		result.Set(key, mappedValue)
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

// CountByHasher groups the values of s by key and returns the number of values
// in each group. Hasher defines hashing and equality for keys, so K need not
// satisfy comparable.
func (s Seq[T]) CountByHasher[
	K any,
	H maphash.Hasher[K],
](
	key func(T) K,
	hasher H,
) *hashmap.Map[K, int, H] {
	result := hashmap.New[K, int](hasher)

	for value := range s.All() {
		groupedKey := key(value)
		count, _ := result.Get(groupedKey)
		result.Set(groupedKey, count+1)
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

// IndexByHasher returns a hash Map from each key to the last value of s that
// produced that key. Hasher defines hashing and equality for keys, so K need
// not satisfy comparable.
func (s Seq[T]) IndexByHasher[
	K any,
	H maphash.Hasher[K],
](
	key func(T) K,
	hasher H,
) *hashmap.Map[K, T, H] {
	result := hashmap.New[K, T](hasher)

	for value := range s.All() {
		result.Set(key(value), value)
	}

	return result
}
