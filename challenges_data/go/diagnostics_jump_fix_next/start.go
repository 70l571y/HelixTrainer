//go:build ignore

package main

func main() {
	println(buildMesage("helix"))
	println(renderCount(3))
	println(makeLable())
}

func buildMessage(name string) string {
	return "editor:" + name
}

func renderCount(n int) string {
	return "count"
}

func makeLabel() string {
	return "ok"
}
