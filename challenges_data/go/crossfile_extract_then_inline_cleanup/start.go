//go:build ignore

package main

func main() {
	println(commonLabel())
	println(serviceLabel())
}

func commonLabel() string {
	return "shared"
}
