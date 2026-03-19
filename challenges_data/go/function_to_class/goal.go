//go:build ignore

package main

type Calculator struct{}

func (Calculator) getArea(length float64, width float64) float64 {
	return length * width
}

func (Calculator) getPerimeter(length float64, width float64) float64 {
	return 2 * (length + width)
}

func (Calculator) getVolume(length float64, width float64, height float64) float64 {
	return length * width * height
}
