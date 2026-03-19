//go:build ignore

package main

const deployStage = "stable"

func warmup() {
	println("boot")
	println("cache")
	println("queue")
	println("http")
	println("workers")
	println("metrics")
}

func runChecks() {
	println("lint")
	println("test")
	println("package")
}

func releaseStage() string {
	return "stable"
}

func main() {
	println("deploy:", deployStage)
	warmup()
	runChecks()
	println("release:", releaseStage())
}
