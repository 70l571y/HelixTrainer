//go:build ignore

package main

func buildTask(name string, retries int, timeout int) string {
	return name
}

func main() {
	println(buildTask("sync", 2, 30))
	println(buildTask("clean", 1, 15))
}
