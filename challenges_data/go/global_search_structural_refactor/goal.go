//go:build ignore

package main

func main() {
	println(modernWrap("main"))
	println(runService())
	println(runWorker())
}

func legacyWrap(name string) string {
	return name
}

func modernWrap(name string) string {
	return name
}
