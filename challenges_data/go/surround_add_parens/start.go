//go:build ignore

package main

import "fmt"

func calculate() int {
	x := 10
	y := 20
	result := x + y
	return result
}

func main() {
	fmt.Println(calculate())
}
