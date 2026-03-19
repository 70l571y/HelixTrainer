//go:build ignore

package main

func main() {
	stateCode := currentStatus()
	println(stateCode)
	if stateCode > 0 {
		println(statusCodeValue())
	}
}

func currentStatus() int {
	return 3
}

func statusCodeValue() int {
	return 8
}
