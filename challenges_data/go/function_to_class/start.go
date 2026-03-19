//go:build ignore

package main

func calculateArea(length float64, width float64) float64 {
	result := length * width
	return result
}

func calculatePerimeter(length float64, width float64) float64 {
	result := 2 * (length + width)
	return result
}

func calculateVolume(length float64, width float64, height float64) float64 {
	result := length * width * height
	return result
}
