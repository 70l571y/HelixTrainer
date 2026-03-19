//go:build ignore

package main

func main() {
	logStart("api")
	logFinish("api")

	logStart("worker")
	logFinish("worker")

	logStart("admin")
	logDone("admin")
}

func logStart(name string) {
	println("start", name)
}

func logFinish(name string) {
	println("finish", name)
}

func logDone(name string) {
	println("done", name)
}
