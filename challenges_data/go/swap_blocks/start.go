//go:build ignore

package main

import "fmt"

func main() {
	data := []int{1, 2, 3}
	fmt.Println(processData(data))
}

func processData(d []int) int {
	return sum(d) * 2
}

func sum(d []int) int {
	total := 0
	for _, v := range d {
		total += v
	}
	return total
}
