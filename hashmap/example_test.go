package hashmap_test

import (
	"fmt"
	"go/types"
	"hash/maphash"
	"strconv"

	"github.com/everettxwzhu/go-toolkit/hashmap"
)

func ExampleMap_goTypesHasher() {
	metadata := hashmap.New[types.Type, string](types.Hasher{})

	// These two slice types are separately allocated but semantically
	// identical according to go/types.Identical.
	storedType := types.NewSlice(types.Typ[types.String])
	equivalentType := types.NewSlice(types.Typ[types.String])
	metadata.Set(storedType, "list of names")

	value, ok := metadata.Get(equivalentType)
	fmt.Println(types.Identical(storedType, equivalentType))
	fmt.Println(value, ok)
	// Output:
	// true
	// list of names true
}

func ExampleMap_MapValues() {
	lengths := hashmap.New[string, int](maphash.ComparableHasher[string]{})
	lengths.Set("one", 3)
	lengths.Set("three", 5)

	labels := lengths.MapValues(func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	})

	value, _ := labels.Get("three")
	fmt.Println(value)
	// Output:
	// three:5
}
