//go:build ignore

package main

// TODO: Move 'dotProduct' and 'transpose' to math_lib.go
// TODO: Move 'sigmoid' from math_lib.go to this file

func processData(data [][]float64) [][]float64 {
	if len(data) == 0 {
		return [][]float64{}
	}
	result := make([][]float64, len(data))
	for i, row := range data {
		result[i] = make([]float64, len(row))
		for j, x := range row {
			result[i][j] = x * 2
		}
	}
	return result
}

// Sigmoid here
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + pow(2.72, -x))
}

func pow(base float64, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		if exp > 0 {
			result *= base
		} else {
			result /= base
		}
	}
	return result
}
