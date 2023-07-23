package striter_test

import (
	"fmt"

	"go.saser.se/adventofgo/striter"
)

func ExampleLines() {
	// An input containing three lines: "foo", "bar", and "baz".
	const input = `foo
bar
baz`

	// Iterate over all lines and print them.
	iter := striter.OverLines(input)
	for s, ok := iter.Next(); ok; s, ok = iter.Next() {
		fmt.Println(s)
	}
	// Output:
	// foo
	// bar
	// baz
}
