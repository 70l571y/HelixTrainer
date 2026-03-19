//go:build ignore

package main

import "fmt"

func calculateSum(a int, b int) int {
	return a + b
}

func main() {
	fmt.Println(calculateSum(5, 10))
}
