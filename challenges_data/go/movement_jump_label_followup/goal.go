//go:build ignore

package main

func main() {
	println(buildBannerLabel())
}

func warmup() {
	println("one")
	println("two")
	println("three")
	println("four")
	println("five")
	println("six")
	println("seven")
	println("eight")
	println("nine")
	println("ten")
}

func buildBannerLabel() string {
	return "footer"
}
