// escape_analysis.go
// This file demonstrates Go's escape analysis.
// Run: go run -gcflags="-m" escape_analysis.go

package main

import "fmt"

func escape() *int {
	x := 42
	return &x // x escapes to heap
}

func noEscape() int {
	y := 99
	return y // y does not escape
}

func main() {
	p := escape()
	fmt.Println("Escaped value:", *p)

	v := noEscape()
	fmt.Println("Non-escaped value:", v)
}
