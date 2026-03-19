//go:build ignore

package main

func main() {
	var runner Runner = workerRunner{}
	println(runner.Execute("helix"))
}
