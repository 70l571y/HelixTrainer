//go:build ignore

package main

import "fmt"

func main() {
	items := []string{
		"apple",
		"banana",
		"cherry",
		"date",
		"elderberry",
		"fig",
		"grape",
	}

	for _, item := range items {
		fmt.Println(item)
	}
}
