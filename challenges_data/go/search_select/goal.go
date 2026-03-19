//go:build ignore

package main

func processBar() int {
	barList := []int{}
	bar := 10
	if bar > 5 {
		println("bar is big")
	}
	return bar
}
