//go:build ignore

package main

func main() {
	println(helperFn(), statusLabel)
}

func helperFn() string {
	return "jump"
}

var statusLabel = "draf"
