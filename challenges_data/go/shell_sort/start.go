package main

import "fmt"

func main() {
	items := []string{
		"banana",
		"apple",
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
