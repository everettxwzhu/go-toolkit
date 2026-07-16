package seq_test

import (
	"fmt"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func Example() {
	values := seq.Range(1, 7).
		Filter(func(value int) bool { return value%2 == 0 }).
		Map(func(value int) int { return value * value }).
		Collect()

	fmt.Println(values)
	// Output: [4 16 36]
}
