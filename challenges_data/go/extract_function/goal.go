//go:build ignore

package main

func areaCircle(r float64) float64 {
	return 3.14 * r * r
}

func volumeCylinder(r float64, h float64) float64 {
	baseArea := areaCircle(r)
	return baseArea * h
}
