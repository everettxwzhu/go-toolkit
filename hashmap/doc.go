// Package hashmap provides a generic hash map whose keys may be of any type.
//
// Unlike Go's built-in map, Map does not require its key type to satisfy
// comparable. Instead, each Map is parameterized by a [maphash.Hasher] that
// defines both hashing and equality for its keys. This permits slices and
// structs containing slices, among other non-comparable values, to be used as
// keys.
//
// A hasher must obey the contract documented by [maphash.Hasher]: keys that
// are equal must produce identical hashes. Keys must not be mutated in a way
// that changes their hash or equality while they are stored in a Map.
package hashmap
