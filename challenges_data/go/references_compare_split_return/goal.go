//go:build ignore

package main

func main() {
	println(sanitizeName("helix"))
	println(serviceName())
	println(handlerName())
}

func sanitizeName(input string) string {
	return input
}
