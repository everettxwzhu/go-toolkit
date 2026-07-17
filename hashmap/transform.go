package hashmap

import "hash/maphash"

// MapValues returns a Map with the same keys as m and values produced by
// applying transform to each entry. It does not modify m.
//
// Keys are shallow-copied and must not be mutated while either Map contains
// them. The original Map must not be mutated while transform is running.
// MapValues returns nil when m is nil.
func (m *Map[K, V, H]) MapValues[U any](
	transform func(K, V) U,
) *Map[K, U, H] {
	if m == nil {
		return nil
	}

	result := &Map[K, U, H]{
		hasher:      m.hasher,
		seed:        m.seed,
		slots:       make([]slot[K, U], len(m.slots)),
		length:      m.length,
		initialized: m.initialized,
	}

	for i := range m.slots {
		current := &m.slots[i]
		if !current.occupied {
			continue
		}

		result.slots[i] = slot[K, U]{
			hash:     current.hash,
			key:      current.key,
			value:    transform(current.key, current.value),
			occupied: true,
		}
	}

	return result
}

// MapKeys returns a Map whose keys are produced by applying transform to each
// entry of m and whose values are unchanged. Hasher defines hashing and
// equality for the transformed keys. It does not modify m.
//
// When multiple entries produce equal keys, one of their values is retained.
// Which value is retained is unspecified because Map iteration order is
// unspecified. The original Map must not be mutated while transform is
// running. MapKeys returns nil when m is nil.
func (m *Map[K, V, H]) MapKeys[
	L any,
	G maphash.Hasher[L],
](
	hasher G,
	transform func(K, V) L,
) *Map[L, V, G] {
	if m == nil {
		return nil
	}

	result := NewWithCapacity[L, V](hasher, m.Len())
	for key, value := range m.All() {
		result.Set(transform(key, value), value)
	}
	return result
}

// MapPairs returns a Map whose keys and values are produced by applying
// transform to each entry of m. Hasher defines hashing and equality for the
// transformed keys. It does not modify m.
//
// When multiple entries produce equal keys, one of their values is retained.
// Which value is retained is unspecified because Map iteration order is
// unspecified. The original Map must not be mutated while transform is
// running. MapPairs returns nil when m is nil.
func (m *Map[K, V, H]) MapPairs[
	L, U any,
	G maphash.Hasher[L],
](
	hasher G,
	transform func(K, V) (L, U),
) *Map[L, U, G] {
	if m == nil {
		return nil
	}

	result := NewWithCapacity[L, U](hasher, m.Len())
	for key, value := range m.All() {
		transformedKey, transformedValue := transform(key, value)
		result.Set(transformedKey, transformedValue)
	}
	return result
}
