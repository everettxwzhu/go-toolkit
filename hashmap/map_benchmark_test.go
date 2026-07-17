package hashmap_test

import (
	"bytes"
	"fmt"
	"hash/maphash"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/hashmap"
	"github.com/rogpeppe/generic/anyhash"
)

var (
	benchmarkValueSink int
	benchmarkBoolSink  bool
)

func BenchmarkGetString(b *testing.B) {
	for _, size := range []int{16, 1_024, 65_536} {
		keys := benchmarkStringKeys(size)

		b.Run(fmt.Sprintf("Builtin/%d", size), func(b *testing.B) {
			values := make(map[string]int, size)
			for i, key := range keys {
				values[key] = i
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				value, ok = values[keys[i%len(keys)]]
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})

		b.Run(fmt.Sprintf("HashMap/%d", size), func(b *testing.B) {
			values := hashmap.NewWithCapacity[string, int](
				maphash.ComparableHasher[string]{},
				size,
			)
			for i, key := range keys {
				values.Set(key, i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				value, ok = values.Get(keys[i%len(keys)])
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})

		b.Run(fmt.Sprintf("AnyHash/%d", size), func(b *testing.B) {
			values := anyhash.NewMap[
				string,
				int,
				maphash.ComparableHasher[string],
			](maphash.ComparableHasher[string]{})
			for i, key := range keys {
				values.Set(key, i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				_, value, ok = values.Get(keys[i%len(keys)])
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})
	}
}

func BenchmarkInsertString(b *testing.B) {
	for _, size := range []int{16, 1_024, 65_536} {
		keys := benchmarkStringKeys(size)

		b.Run(fmt.Sprintf("Builtin/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := make(map[string]int)
				for i, key := range keys {
					values[key] = i
				}
				benchmarkValueSink = len(values)
			}
		})

		b.Run(fmt.Sprintf("BuiltinPreallocated/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := make(map[string]int, size)
				for i, key := range keys {
					values[key] = i
				}
				benchmarkValueSink = len(values)
			}
		})

		b.Run(fmt.Sprintf("HashMap/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := hashmap.New[string, int](
					maphash.ComparableHasher[string]{},
				)
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})

		b.Run(fmt.Sprintf("HashMapPreallocated/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := hashmap.NewWithCapacity[string, int](
					maphash.ComparableHasher[string]{},
					size,
				)
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})

		b.Run(fmt.Sprintf("AnyHash/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := anyhash.NewMap[
					string,
					int,
					maphash.ComparableHasher[string],
				](maphash.ComparableHasher[string]{})
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})
	}
}

func BenchmarkGetBytes(b *testing.B) {
	for _, size := range []int{16, 1_024} {
		keys := benchmarkByteKeys(size)

		b.Run(fmt.Sprintf("StringMap/%d", size), func(b *testing.B) {
			values := make(map[string]int, size)
			for i, key := range keys {
				values[string(key)] = i
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				value, ok = values[string(keys[i%len(keys)])]
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})

		b.Run(fmt.Sprintf("HashMap/%d", size), func(b *testing.B) {
			values := hashmap.NewWithCapacity[[]byte, int](bytesHasher{}, size)
			for i, key := range keys {
				values.Set(key, i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				value, ok = values.Get(keys[i%len(keys)])
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})

		b.Run(fmt.Sprintf("AnyHash/%d", size), func(b *testing.B) {
			values := anyhash.NewMap[[]byte, int, bytesHasher](bytesHasher{})
			for i, key := range keys {
				values.Set(key, i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				_, value, ok = values.Get(keys[i%len(keys)])
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})

		b.Run(fmt.Sprintf("Linear/%d", size), func(b *testing.B) {
			var values linearBytesMap
			for i, key := range keys {
				values.Set(key, i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			value := 0
			ok := false
			for i := 0; i < b.N; i++ {
				value, ok = values.Get(keys[i%len(keys)])
			}
			benchmarkValueSink = value
			benchmarkBoolSink = ok
		})
	}
}

func BenchmarkInsertBytes(b *testing.B) {
	for _, size := range []int{16, 1_024} {
		keys := benchmarkByteKeys(size)

		b.Run(fmt.Sprintf("StringMap/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := make(map[string]int)
				for i, key := range keys {
					values[string(key)] = i
				}
				benchmarkValueSink = len(values)
			}
		})

		b.Run(fmt.Sprintf("StringMapPreallocated/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := make(map[string]int, size)
				for i, key := range keys {
					values[string(key)] = i
				}
				benchmarkValueSink = len(values)
			}
		})

		b.Run(fmt.Sprintf("HashMap/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := hashmap.New[[]byte, int](bytesHasher{})
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})

		b.Run(fmt.Sprintf("HashMapPreallocated/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := hashmap.NewWithCapacity[[]byte, int](
					bytesHasher{},
					size,
				)
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})

		b.Run(fmt.Sprintf("AnyHash/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				values := anyhash.NewMap[[]byte, int, bytesHasher](bytesHasher{})
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})

		b.Run(fmt.Sprintf("Linear/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ReportMetric(float64(size), "entries/op")
			for range b.N {
				var values linearBytesMap
				for i, key := range keys {
					values.Set(key, i)
				}
				benchmarkValueSink = values.Len()
			}
		})
	}
}

type linearBytesMap struct {
	entries []linearBytesEntry
}

type linearBytesEntry struct {
	key   []byte
	value int
}

func (m *linearBytesMap) Get(key []byte) (int, bool) {
	for i := range m.entries {
		if bytes.Equal(m.entries[i].key, key) {
			return m.entries[i].value, true
		}
	}
	return 0, false
}

func (m *linearBytesMap) Set(key []byte, value int) {
	for i := range m.entries {
		if bytes.Equal(m.entries[i].key, key) {
			m.entries[i].value = value
			return
		}
	}
	m.entries = append(m.entries, linearBytesEntry{key: key, value: value})
}

func (m *linearBytesMap) Len() int {
	return len(m.entries)
}

func benchmarkStringKeys(size int) []string {
	keys := make([]string, size)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}

func benchmarkByteKeys(size int) [][]byte {
	stringKeys := benchmarkStringKeys(size)
	keys := make([][]byte, size)
	for i := range stringKeys {
		keys[i] = []byte(stringKeys[i])
	}
	return keys
}
