//go:build ignore

package main

func main() {
	println(runJob("sync", 2, 30))
	println(runNightly())
}
