//go:build ignore

package main

import "strings"

func main() {
	name := " helix "
	trimmed := name
	ready := strings.TrimSpace((name))

	println(trimmed, ready)
}
