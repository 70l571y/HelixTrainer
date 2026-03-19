//go:build ignore

package main

func main() {
	println(runTask("sync", 30, 2))
	println(runTask("clean", 15, 1))
	println(runNightly())
}
