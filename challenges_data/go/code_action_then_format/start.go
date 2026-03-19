//go:build ignore

package main

import "strings"

func main() {
	data := map[string]string{
		"name": " helix ",
		"mode": " normal ",
	}

	println(strings.TrimSpace(data["name"]))
}
