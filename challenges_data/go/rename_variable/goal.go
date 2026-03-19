//go:build ignore

package main

func processData(data []int) int {
	total := 0
	for _, i := range data {
		total += i
	}
	return total
}

func average(data []int) float64 {
	total := processData(data)
	return float64(total) / float64(len(data))
}
