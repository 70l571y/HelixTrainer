//go:build ignore

package main

func main() {
	printLn(loadReport())
	println(buildSummary())
}

func buildSummary() string {
	return "summary"
}
