//go:build ignore

package main

func processData(d []int) int {
	x := 0
	for _, i := range d {
		x += i
	}
	return x
}

func average(d []int) float64 {
	x := processData(d)
	return float64(x) / float64(len(d))
}
