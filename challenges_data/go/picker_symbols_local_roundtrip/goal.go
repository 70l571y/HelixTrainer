//go:build ignore

package main

func buildHeader() string {
	return "header"
}

func buildBannerFooter() string {
	return "footer"
}

func buildSummary() string {
	return buildHeader() + ":" + buildBannerFooter()
}

func main() {
	println(buildSummary())
}
