//go:build ignore

package main

func main() {
	println(sanitizeCode("main"))
	println(runService())
	println(runWorker())
}
