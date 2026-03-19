//go:build ignore

package main

// TODO: Move 'dotProduct' and 'transpose' to math_lib.go
// TODO: Move 'sigmoid' from math_lib.go to this file

func dotProduct(v1 []float64, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic("Vectors must be the same length")
	}
	sum := 0.0
	for i := range v1 {
		sum += v1[i] * v2[i]
	}
	return sum
}

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

func transpose(matrix [][]float64) [][]float64 {
	if len(matrix) == 0 {
		return [][]float64{}
	}
	rows := len(matrix)
	cols := len(matrix[0])
	result := make([][]float64, cols)
	for i := range result {
		result[i] = make([]float64, rows)
		for j := range result[i] {
			result[i][j] = matrix[j][i]
		}
	}
	return result
}

// Sigmoid here
