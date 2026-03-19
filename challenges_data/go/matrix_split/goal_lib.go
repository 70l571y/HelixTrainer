//go:build ignore

package main

// TODO: Move 'sigmoid' function to matrix_processor.go

// dotProduct here
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

// transpose here
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
