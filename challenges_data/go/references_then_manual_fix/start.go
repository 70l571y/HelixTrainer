//go:build ignore

package main

func main() {
	println(buildCode())
	println(runService())
	println(runWorker())
}

func buildCode() string {
	return "local"
}
