//go:build ignore

package main

func main() {
	auditStart("api")
	auditFinish("api")

	auditStart("worker")
	auditFinish("worker")

	auditStart("admin")
	auditDone("admin")
}

func auditStart(name string) {
	println("start", name)
}

func auditFinish(name string) {
	println("finish", name)
}

func auditDone(name string) {
	println("done", name)
}
