//go:build ignore

package main

// TODO: Move 'sigmoid' function to matrix_processor.go

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

// dotProduct here

// transpose here
