package option_test

import (
	"fmt"
	"strings"

	"github.com/everettxwzhu/go-toolkit/option"
)

func Example() {
	value := option.Some("gopher").
		Filter(func(value string) bool { return len(value) > 3 }).
		Map(strings.ToUpper).
		OrElse("UNKNOWN")

	fmt.Println(value)
	// Output: GOPHER
}
