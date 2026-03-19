//go:build ignore

package main

func sendRequest(url string, timeout int, retries int) {
	println(url, timeout, retries)
}

func main() {
	sendRequest("api", 30, 2)
	sendRequest("worker", 10, 1)
}
