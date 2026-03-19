//go:build ignore

package main

func buildTask(name string, timeout int) string {
	return name
}

func main() {
	println(buildTask("sync", 30))
	println(buildTask("clean", 15))
}
