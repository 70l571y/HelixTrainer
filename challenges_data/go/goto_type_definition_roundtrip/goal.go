//go:build ignore

package main

func main() {
	var formatter Formatter = bannerFormatter{}
	println(formatter.Format("helix"))
}
