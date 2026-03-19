//go:build ignore

package main

func fetchData(url string, method string, retries int, cache bool, timeout int) {
	println("Fetching", url, "...")
}

func main() {
	// Call 1: Simple values
	fetchData("api/v1", "GET", 3, true, 30)

	// Call 2: Complex values (gw would be slow here)
	fetchData("api/v2", "POST", 5, false, 6000)

	// Call 3: Mixed types
	fetchData("api/v3", "PUT", 1, nil, 999)
}
