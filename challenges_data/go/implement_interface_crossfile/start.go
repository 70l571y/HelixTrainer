//go:build ignore

package main

type Saver interface {
	Save() string
}

func main() {
	var saver Saver = FileStore{}
	println(saver.Save())
}
