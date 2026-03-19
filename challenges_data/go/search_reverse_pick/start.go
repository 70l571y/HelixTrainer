//go:build ignore

package main

func main() {
	statuses := []string{
		"draft",
		"queued",
		"draft",
		"published",
	}

	println(statuses[0], statuses[1], statuses[2], statuses[3])
}
