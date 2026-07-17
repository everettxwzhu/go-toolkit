package hashset_test

import (
	"bytes"
	"fmt"
	"hash/maphash"
	"strings"

	"github.com/everettxwzhu/go-toolkit/hashset"
)

type bytesHasher struct{}

func (bytesHasher) Hash(hash *maphash.Hash, value []byte) {
	hash.Write(value)
}

func (bytesHasher) Equal(left, right []byte) bool {
	return bytes.Equal(left, right)
}

type caseInsensitiveHasher struct{}

func (caseInsensitiveHasher) Hash(hash *maphash.Hash, value string) {
	hash.WriteString(strings.ToLower(value))
}

func (caseInsensitiveHasher) Equal(left, right string) bool {
	return strings.ToLower(left) == strings.ToLower(right)
}

func ExampleSet_bytes() {
	values := hashset.New[[]byte](
		bytesHasher{},
		[]byte("alpha"),
		[]byte("beta"),
	)

	// A different slice with the same bytes finds the stored value.
	fmt.Println(values.Contains([]byte("alpha")))
	values.Add([]byte("alpha"))
	fmt.Println(values.Len())
	// Output:
	// true
	// 2
}

func ExampleSet_caseInsensitive() {
	names := hashset.New(
		caseInsensitiveHasher{},
		"Go",
		"Rust",
	)

	names.Add("GO")
	fmt.Println(names.Contains("go"))
	fmt.Println(names.Len())
	// Output:
	// true
	// 2
}
