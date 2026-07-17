package seq_test

import (
	"fmt"
	"hash/maphash"
	"strings"

	"github.com/everettxwzhu/go-toolkit/seq"
)

type exampleCaseInsensitiveHasher struct{}

func (exampleCaseInsensitiveHasher) Hash(
	hash *maphash.Hash,
	value string,
) {
	hash.WriteString(strings.ToLower(value))
}

func (exampleCaseInsensitiveHasher) Equal(left, right string) bool {
	return strings.ToLower(left) == strings.ToLower(right)
}

func ExampleSeq_GroupByHasher() {
	groups := seq.Of("Go", "GO", "Rust").GroupByHasher(
		func(value string) string { return value },
		exampleCaseInsensitiveHasher{},
	)

	goValues, _ := groups.Get("go")
	fmt.Println(goValues)
	// Output:
	// [Go GO]
}
