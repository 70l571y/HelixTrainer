//go:build ignore

package main

func main() {
	var formatter Formatter = upperFormatter{}
	println(formatter.Format("helix"))
}
