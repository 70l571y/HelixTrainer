//go:build ignore

package main

func main() {
	println(runTask("sync", 2))
	println(runTask("clean", 1))
	println(runNightly())
}
