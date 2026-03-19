//go:build ignore

package main

func main() {
	statusCode := currentStatus()
	println(statusCode)
	if statusCode > 0 {
		println(statusCodeValue())
	}
}

func currentStatus() int {
	return 3
}

func statusCodeValue() int {
	return 8
}
