//go:build ignore

package main

func buildHeader() string {
	return "header"
}

func buildFooter() string {
	return "footer"
}

func buildSummary() string {
	return buildHeader() + ":" + buildFooter()
}

func main() {
	println(buildSummary())
}
