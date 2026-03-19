//go:build ignore

package main

func main() {
	println(buildMessage("helix"))
	println(renderCount(3))
	println(makeLabel())
}

func buildMessage(name string) string {
	return name
}

func renderCount(count int) string {
	return "count"
}

func makeLabel() string {
	return "ok"
}
