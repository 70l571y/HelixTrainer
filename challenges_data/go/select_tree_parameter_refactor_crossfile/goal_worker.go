//go:build ignore

package main

func runNightly() string {
	return runTask("nightly", 3)
}

func formatTask(name string, retries int) string {
	return name
}
