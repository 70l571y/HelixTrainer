//go:build ignore

package main

func runNightly() string {
	return runTask("nightly", 60, 3)
}

func formatTask(name string, timeout int, retries int) string {
	return name
}
