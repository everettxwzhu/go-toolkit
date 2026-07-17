// Package hashset provides a generic collection of distinct values whose
// equality and hashing are defined by a [maphash.Hasher].
//
// Unlike package set, hashset does not require its element type to satisfy
// comparable. This permits slices and structs containing slices, among other
// non-comparable values, to be stored in a Set. A Hasher may also define a
// custom equality relation for comparable values, such as case-insensitive
// equality for strings.
//
// A hasher must obey the contract documented by [maphash.Hasher]. Values must
// not be mutated in a way that changes their hash or equality while stored in a
// Set.
package hashset
