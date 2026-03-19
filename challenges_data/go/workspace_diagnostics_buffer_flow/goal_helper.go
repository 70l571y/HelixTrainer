//go:build ignore

package main

func helperStatus() string {
	return buildStatus()
}

func buildStatus() string {
	return "ok"
}
