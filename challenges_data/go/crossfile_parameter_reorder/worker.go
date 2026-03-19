//go:build ignore

package main

func runNightly() string {
	return runJob("nightly", 60, 3)
}
