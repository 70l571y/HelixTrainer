//go:build ignore

package main

import "strings"

func main() {
	println(strings.TrimSpace("  ready  "))
	println(buildCode())
}
