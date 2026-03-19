//go:build ignore

package main

func runTask(name string, timeout int, retries int) string {
	return formatTask(name, timeout, retries)
}
