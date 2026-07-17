package hashmap

import (
	"hash/maphash"
	"iter"
	"sync"
)

const (
	minSlotCount          = 8
	loadFactorNumerator   = 3
	loadFactorDenominator = 4
)

// Map is a hash table from keys of type K to values of type V. H defines the
// hash function and equality relation used for keys, so K need not satisfy
// comparable.
//
// Use [New] or [NewWithCapacity] to construct a Map. A Map must not be copied
// after first use; use [Map.Clone] to create an independent copy. A Map is not
// safe for concurrent mutation.
type Map[K, V any, H maphash.Hasher[K]] struct {
	hasher      H
	seed        maphash.Seed
	slots       []slot[K, V]
	length      int
	initialized bool
	hashes      sync.Pool
}

type slot[K, V any] struct {
	hash     uint64
	key      K
	value    V
	occupied bool
}

// New returns an empty Map that uses hasher to hash and compare keys.
func New[K, V any, H maphash.Hasher[K]](hasher H) *Map[K, V, H] {
	return NewWithCapacity[K, V](hasher, 0)
}

// NewWithCapacity returns an empty Map with enough initial storage for
// approximately capacity entries. It panics when capacity is negative.
func NewWithCapacity[K, V any, H maphash.Hasher[K]](
	hasher H,
	capacity int,
) *Map[K, V, H] {
	if capacity < 0 {
		panic("hashmap: negative capacity")
	}

	return &Map[K, V, H]{
		hasher:      hasher,
		seed:        maphash.MakeSeed(),
		slots:       make([]slot[K, V], slotCountForCapacity(capacity)),
		initialized: true,
	}
}

// Len returns the number of entries in m. A nil Map has length zero.
func (m *Map[K, V, H]) Len() int {
	if m == nil {
		return 0
	}

	return m.length
}

// Get returns the value associated with key and reports whether key is present.
func (m *Map[K, V, H]) Get(key K) (V, bool) {
	if m == nil || m.length == 0 {
		var zero V
		return zero, false
	}

	hash := m.hashKey(key)
	slotIndex, ok := m.find(key, hash)
	if !ok {
		var zero V
		return zero, false
	}

	return m.slots[slotIndex].value, true
}

// ContainsKey reports whether key is present in m.
func (m *Map[K, V, H]) ContainsKey(key K) bool {
	if m == nil || m.length == 0 {
		return false
	}

	_, ok := m.find(key, m.hashKey(key))
	return ok
}

// Set associates value with key. It returns the previous value and true when
// an equal key was already present. When replacing an entry, Set retains the
// key originally stored in the Map.
func (m *Map[K, V, H]) Set(key K, value V) (V, bool) {
	if m == nil {
		panic("hashmap: Set called on nil Map")
	}
	m.initialize()

	hash := m.hashKey(key)
	slotIndex, ok := m.find(key, hash)
	if ok {
		current := &m.slots[slotIndex]
		previous := current.value
		current.value = value
		return previous, true
	}

	if m.shouldGrow() {
		m.resize(len(m.slots) * 2)
		slotIndex, _ = m.find(key, hash)
	}

	m.slots[slotIndex] = slot[K, V]{
		hash:     hash,
		key:      key,
		value:    value,
		occupied: true,
	}
	m.length++

	var zero V
	return zero, false
}

// Delete removes the entry whose key is equal to key. It returns the removed
// value and reports whether an entry was found.
func (m *Map[K, V, H]) Delete(key K) (V, bool) {
	if m == nil || m.length == 0 {
		var zero V
		return zero, false
	}

	slotIndex, ok := m.find(key, m.hashKey(key))
	if !ok {
		var zero V
		return zero, false
	}

	previous := m.slots[slotIndex].value
	m.deleteAt(slotIndex)
	m.length--
	return previous, true
}

// Clear removes all entries from m while retaining storage for future entries.
// Clear has no effect on a nil Map.
func (m *Map[K, V, H]) Clear() {
	if m == nil {
		return
	}

	clear(m.slots)
	m.length = 0
}

// Clone returns an independent shallow copy of m. Keys and values themselves
// are not cloned. Clone returns nil when m is nil.
func (m *Map[K, V, H]) Clone() *Map[K, V, H] {
	if m == nil {
		return nil
	}

	clone := &Map[K, V, H]{
		hasher:      m.hasher,
		seed:        m.seed,
		slots:       append([]slot[K, V](nil), m.slots...),
		length:      m.length,
		initialized: m.initialized,
	}

	return clone
}

// All returns an iterator over the entries of m in unspecified order. A nil
// Map produces no entries.
//
// The Map must not be mutated while the iterator is running.
func (m *Map[K, V, H]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}

		for i := range m.slots {
			current := &m.slots[i]
			if current.occupied && !yield(current.key, current.value) {
				return
			}
		}
	}
}

// Keys returns an iterator over the keys of m in unspecified order. A nil Map
// produces no keys.
//
// The Map must not be mutated while the iterator is running.
func (m *Map[K, V, H]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range m.All() {
			if !yield(key) {
				return
			}
		}
	}
}

// Values returns an iterator over the values of m in unspecified order. A nil
// Map produces no values.
//
// The Map must not be mutated while the iterator is running.
func (m *Map[K, V, H]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range m.All() {
			if !yield(value) {
				return
			}
		}
	}
}

func (m *Map[K, V, H]) initialize() {
	if m.initialized {
		return
	}

	m.seed = maphash.MakeSeed()
	m.slots = make([]slot[K, V], minSlotCount)
	m.initialized = true
}

func (m *Map[K, V, H]) hashKey(key K) uint64 {
	var hash *maphash.Hash
	if pooled := m.hashes.Get(); pooled != nil {
		hash = pooled.(*maphash.Hash)
	} else {
		hash = new(maphash.Hash)
	}

	hash.SetSeed(m.seed)
	m.hasher.Hash(hash, key)
	sum := hash.Sum64()
	m.hashes.Put(hash)
	return sum
}

func (m *Map[K, V, H]) find(key K, hash uint64) (int, bool) {
	mask := len(m.slots) - 1
	slotIndex := indexFor(hash, len(m.slots))

	for {
		current := &m.slots[slotIndex]
		if !current.occupied {
			return slotIndex, false
		}
		if current.hash == hash && m.hasher.Equal(key, current.key) {
			return slotIndex, true
		}
		slotIndex = (slotIndex + 1) & mask
	}
}

func (m *Map[K, V, H]) shouldGrow() bool {
	threshold := growThreshold(len(m.slots))
	return m.length+1 > threshold
}

func (m *Map[K, V, H]) resize(slotCount int) {
	previous := m.slots
	m.slots = make([]slot[K, V], slotCount)

	for i := range previous {
		if previous[i].occupied {
			m.place(previous[i])
		}
	}
}

func (m *Map[K, V, H]) place(current slot[K, V]) {
	mask := len(m.slots) - 1
	slotIndex := indexFor(current.hash, len(m.slots))
	for m.slots[slotIndex].occupied {
		slotIndex = (slotIndex + 1) & mask
	}
	m.slots[slotIndex] = current
}

func (m *Map[K, V, H]) deleteAt(slotIndex int) {
	mask := len(m.slots) - 1
	hole := slotIndex

	for scan := (slotIndex + 1) & mask; m.slots[scan].occupied; scan = (scan + 1) & mask {
		home := indexFor(m.slots[scan].hash, len(m.slots))
		distanceToHole := (hole - home) & mask
		distanceToEntry := (scan - home) & mask
		if distanceToHole < distanceToEntry {
			m.slots[hole] = m.slots[scan]
			hole = scan
		}
	}

	var zero slot[K, V]
	m.slots[hole] = zero
}

func slotCountForCapacity(capacity int) int {
	slotCount := minSlotCount
	for capacity > growThreshold(slotCount) {
		if slotCount > int(^uint(0)>>1)/2 {
			panic("hashmap: capacity too large")
		}
		slotCount *= 2
	}
	return slotCount
}

func growThreshold(slotCount int) int {
	return slotCount / loadFactorDenominator * loadFactorNumerator
}

func indexFor(hash uint64, slotCount int) int {
	return int(hash & uint64(slotCount-1))
}
