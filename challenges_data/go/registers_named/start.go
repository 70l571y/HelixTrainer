//go:build ignore

package main

import "fmt"

func funcA() {
	fmt.Println("Function A")
}

func funcB() {
	fmt.Println("Function B")
}

func funcC() {
	fmt.Println("Function C")
}

func main() {
	funcC()
	funcA()
	funcB()
}
