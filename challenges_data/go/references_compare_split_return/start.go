//go:build ignore

package main

func main() {
	println(normalizeName("helix"))
	println(serviceName())
	println(handlerName())
}

func normalizeName(input string) string {
	return input
}
