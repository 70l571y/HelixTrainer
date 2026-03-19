//go:build ignore

package main

func main() {
	println(apiError())
	println(serviceError())
	println(workerError())
}

func apiError() string {
	return "error: api"
}
