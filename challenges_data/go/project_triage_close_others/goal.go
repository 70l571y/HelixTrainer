//go:build ignore

package main

func main() {
	println(loadReport())
	println(buildSummary())
}

func buildSummary() string {
	return "summary"
}
