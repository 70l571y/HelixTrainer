//go:build ignore

package main

func main() {
	println(helperFunc(), statusLabel)
}

func helperFunc() string {
	return "jump"
}

var statusLabel = "draft"
