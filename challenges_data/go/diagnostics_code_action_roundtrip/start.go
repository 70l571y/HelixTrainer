//go:build ignore

package main

import "strings"

func main() {
	value := "  helix  "
	println(strings.TrimSpace(value))
}
