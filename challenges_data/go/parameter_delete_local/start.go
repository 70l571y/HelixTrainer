//go:build ignore

package main

func buildTask(name string, timeout int, retries int) string {
	return name
}

func main() {
	println(buildTask("sync", 30, 2))
	println(buildTask("clean", 15, 1))
}
