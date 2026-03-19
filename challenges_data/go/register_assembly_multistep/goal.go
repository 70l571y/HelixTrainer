//go:build ignore

package main

func main() {
	partA := "api"
	partB := "::v2"
	partC := "::ready"
	result := "api::v2::ready"

	println(partA, partB, partC, result)
}
