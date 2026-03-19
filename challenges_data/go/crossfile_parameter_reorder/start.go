//go:build ignore

package main

func main() {
	println(runJob("sync", 30, 2))
	println(runNightly())
}
