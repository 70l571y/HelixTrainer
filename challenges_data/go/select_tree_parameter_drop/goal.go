//go:build ignore

package main

func sendRequest(url string, retries int) {
	println(url, retries)
}

func main() {
	sendRequest("api", 2)
	sendRequest("worker", 1)
}
