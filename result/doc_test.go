package result_test

import (
	"fmt"
	"strconv"

	"github.com/everettxwzhu/go-toolkit/result"
)

func Example() {
	r := result.From(strconv.Atoi("21")).
		Map(func(value int) int { return value * 2 })

	value, err := r.Get()
	fmt.Println(value, err)
	// Output: 42 <nil>
}
